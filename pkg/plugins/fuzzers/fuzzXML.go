package fuzzers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/cfsdes/nucke/internal/parsers"
	"github.com/cfsdes/nucke/pkg/globals"
	"github.com/cfsdes/nucke/pkg/plugins/detections"
	"github.com/cfsdes/nucke/pkg/plugins/utils"
	"github.com/cfsdes/nucke/pkg/requests"
)

func FuzzXML(r *http.Request, client *http.Client, pluginDir string, payloads []string, matcher detections.Matcher) (bool, string, string, string, string, string, []detections.Result) {
	req := requests.CloneReq(r)

	// Result channel
	resultChan := make(chan detections.Result)

	// Counter of channels opened
	var channelsOpened int

	// Array com os resultados de cada teste executado falho
	var logScans []detections.Result

	// Check if content type is XML
	if req.Header.Get("Content-Type") != "application/xml" && req.Header.Get("Content-Type") != "text/xml" {
		return false, "", "", "", "", "", nil
	}

	// Get request body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		if globals.Debug {
			fmt.Println("fuzzXML:", err)
		}
		return false, "", "", "", "", "", nil
	}

	// Restore request body
	req.Body = ioutil.NopCloser(bytes.NewReader(body))

	// Find XML tags in request body and replace them with payloads
	re := regexp.MustCompile(`<([^/][^>]+)>([^<]+)</([^>]+)>`)
	matches := re.FindAllStringSubmatch(string(body), -1)
	for _, match := range matches {
		for _, payload := range payloads {

			// Delay between requests
			time.Sleep(time.Duration(globals.Delay) * time.Millisecond)

			// Update payloads {{.params}}
			payload = parsers.ParsePayload(payload)

			// Copy Request
			reqCopy := requests.CloneReq(req)

			// Create a new request body with the tag replaced by a payload
			var newBody string

			payload = strings.Replace(payload, "{{.original}}", match[2], -1)
			newBody = strings.Replace(string(body), match[0], fmt.Sprintf("<%s>%s</%s>", match[1], payload, match[3]), -1)

			// Set request body
			reqCopy.Body = ioutil.NopCloser(strings.NewReader(newBody))

			// Get raw request
			rawReq := requests.RequestToRaw(reqCopy)

			// Send request
			start := time.Now()
			responses := requests.Do(reqCopy, client)

			// Get response time
			elapsed := int(time.Since(start).Seconds())

			// Extract OOB ID
			oobID := utils.ExtractOobID(payload)

			// Check if match vulnerability
			for _, resp := range responses {
				channelsOpened++
				go detections.MatchCheck(pluginDir, matcher, resp, elapsed, oobID, rawReq, payload, match[0], resultChan)
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
