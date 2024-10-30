package fuzzers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cfsdes/nucke/internal/parsers"
	"github.com/cfsdes/nucke/pkg/globals"
	"github.com/cfsdes/nucke/pkg/plugins/detections"
	"github.com/cfsdes/nucke/pkg/plugins/utils"
	"github.com/cfsdes/nucke/pkg/requests"
)

func FuzzFormData(r *http.Request, client *http.Client, pluginDir string, payloads []string, matcher detections.Matcher) (bool, string, string, string, string, string, []detections.Result) {
	req := requests.CloneReq(r)

	// Result channel
	resultChan := make(chan detections.Result)

	// Array com os resultados de cada teste executado falho
	var logScans []detections.Result

	// Counter of channels opened
	var channelsOpened int

	// Check if method is POST and content type is application/x-www-form-urlencoded
	if !(req.Method == http.MethodPost && strings.Contains(req.Header.Get("Content-Type"), "application/x-www-form-urlencoded")) {
		return false, "", "", "", "", "", nil
	}

	// Get form data parameters from request body
	if err := req.ParseForm(); err != nil {
		if globals.Debug {
			fmt.Println("fuzzFormData:", err)
		}
		return false, "", "", "", "", "", nil
	}

	// Get request body
	body := req.PostForm.Encode()

	// For each parameter, send a new request with the parameter replaced by a payload
	for key, values := range req.PostForm {
		for _, payload := range payloads {

			// Delay between requests
			time.Sleep(time.Duration(globals.Delay) * time.Millisecond)

			// Update payloads {{.params}}
			payload = parsers.ParsePayload(payload)

			// Create a new request body with the parameter replaced by a payload
			var newBody string

			payload = strings.Replace(payload, "{{.original}}", values[0], -1)
			newBody = strings.Replace(body, fmt.Sprintf("%s=%s", key, url.QueryEscape(values[0])), fmt.Sprintf("%s=%s", key, url.QueryEscape(payload)), -1)

			// Set request body
			reqBody := strings.NewReader(newBody)

			// Create a new request with the updated form data
			newReq, err := http.NewRequest(req.Method, req.URL.String(), reqBody)
			if err != nil {
				if globals.Debug {
					fmt.Println("fuzzFormData:", err)
				}
				return false, "", "", "", "", "", nil
			}

			// Copy headers from original request to new request
			newReq.Header = req.Header

			// Get raw request
			rawReq := requests.RequestToRaw(newReq)

			// Send request
			start := time.Now()
			responses := requests.Do(newReq, client)

			// Get response time
			elapsed := int(time.Since(start).Seconds())

			// Extract OOB ID
			oobID := utils.ExtractOobID(payload)

			// Check if match vulnerability
			for _, resp := range responses {
				channelsOpened++
				go detections.MatchCheck(pluginDir, matcher, resp, elapsed, oobID, rawReq, payload, key, resultChan)
			}
		}
	}

	// Wait for any goroutine to send a result to the channel
	for i := 0; i < channelsOpened; i++ {
		res := <-resultChan
		log := detections.Result{
			Found:   res.Found,
			URL:     res.URL,
			Payload: res.Payload,
			Param:   res.Param,
			RawReq:  res.RawReq,
			RawResp: res.RawResp,
			ResBody: res.ResBody,
		}
		logScans = append(logScans, log)
	}

	for _, res := range logScans {
		if res.Found {
			return true, res.URL, res.Payload, res.Param, res.RawReq, res.RawResp, logScans
		}
	}

	return false, "", "", "", "", "", logScans
}
