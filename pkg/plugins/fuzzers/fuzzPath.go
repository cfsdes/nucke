package fuzzers

import (
	"net/http"
	"io/ioutil"
	"bytes"
	"time"
	"fmt"
	"strings"

	"github.com/cfsdes/nucke/pkg/plugins/detections"
	"github.com/cfsdes/nucke/pkg/requests"
	"github.com/cfsdes/nucke/pkg/globals"
	"github.com/cfsdes/nucke/internal/parsers"
	"github.com/cfsdes/nucke/pkg/plugins/utils"
)

func FuzzPath(r *http.Request, client *http.Client, payloads []string, matcher detections.Matcher, location string) (bool, string, string, string, string, string, []detections.Result) {
	req := requests.CloneReq(r)

	// Extract segments from URL path
	segments := strings.Split(req.URL.Path, "/")

	// Result channel
	resultChan := make(chan detections.Result)

	// Array com os resultados de cada teste executado falho
    var logScans []detections.Result

	// Get request body, if method is POST
	var body []byte
	var err error
	body, err = ioutil.ReadAll(req.Body)
	if err != nil {
		// handle error
		if globals.Debug {
			fmt.Println("FuzzPath:", err)
		}
		return false, "", "", "", "", "", nil
	}

	// Determine the location to inject payloads
	injectIndexes := []int{}
	if location == "*" {
		injectIndexes = make([]int, len(segments))
		for i := range injectIndexes {
			injectIndexes[i] = i
		}
	} else if location == "last" {
		lastIndex := len(segments) - 1
		injectIndexes = []int{lastIndex}
	}

	// For each segment of the path, inject payloads according to the location
	for _, index := range injectIndexes {
		segment := segments[index]

		// Create a new payload with the original segment replaced
		for _, payload := range payloads {
			// Delay between requests
            time.Sleep(time.Duration(globals.Delay) * time.Millisecond)
			
			// Replace "{{.original}}" with the current segment in the payload
			payload = strings.Replace(payload, "{{.original}}", segment, -1)

			// Update payloads {{.params}}
			payload = parsers.ParsePayload(payload)

			// Create a new path with the payload segment
			newSegments := make([]string, len(segments))
			copy(newSegments, segments)
			newSegments[index] = payload

			// Join the modified segments to form the new path
			newPath := strings.Join(newSegments, "/")
			reqCopy := requests.CloneReq(req)
			reqCopy.URL.Path = newPath

			// Add request body
			reqCopy.Body = ioutil.NopCloser(bytes.NewReader(body))

			// Get raw request
			rawReq := requests.RequestToRaw(reqCopy)

			// Send request
			start := time.Now()
			resp, err := client.Do(reqCopy)
			if err != nil {
				// handle error
				if globals.Debug {
					fmt.Println("FuzzPath:", err)
				}
				return false, "", "", "", "", "", nil
			}

			// Get response time
			elapsed := int(time.Since(start).Seconds())

			// Extract OOB ID
			oobID := utils.ExtractOobID(payload)

			// Check if match vulnerability
			go detections.MatchCheck(matcher, resp, elapsed, oobID, rawReq, payload, fmt.Sprintf("segment %d", index), resultChan)
		}
	}

	// Wait for any goroutine to send a result to the channel
	for i := 0; i < len(injectIndexes)*len(payloads); i++ {
		res := <-resultChan
		if res.Found {
			return true, res.RawReq, res.URL, res.Payload, res.Param, res.RawResp, nil
		} else {
            log := detections.Result{
                Found: false,
                RawReq: res.RawReq,
                URL: res.URL,
                Payload: res.Payload,
                Param: res.Param,
                RawResp: res.RawResp,
				ResBody: res.ResBody,
            }
            logScans = append(logScans, log)
        }
	}

	return false, "", "", "", "", "", logScans
}
