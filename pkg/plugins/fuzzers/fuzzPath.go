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
	"github.com/cfsdes/nucke/internal/initializers"
	"github.com/cfsdes/nucke/internal/parsers"
)

func FuzzPath(r *http.Request, client *http.Client, payloads []string, matcher detections.Matcher) (bool, string, string, string, string, string) {
	req := requests.CloneReq(r)

	// Extract parameters from URL
	segments := strings.Split(req.URL.Path, "/")

	// Result channel
	resultChan := make(chan detections.Result)

	// Get request body, if method is POST
	var body []byte
	if req.Method == http.MethodPost {
		var err error
		body, err = ioutil.ReadAll(req.Body)
		if err != nil {
			// handle error
			if initializers.Debug {
				fmt.Println("FuzzPath:", err)
			}
			return false, "", "", "", "", ""
		}
	}

	// For each segment of the path, send a new request with the segment replaced by a payload
	for i, _ := range segments {
		// Create a new path with the segment replaced by a payload
		for _, payload := range payloads {

			// Update payloads {{.params}}
			payload = parsers.ParsePayload(payload)

			newSegments := make([]string, len(segments))
			copy(newSegments, segments)
			newSegments[i] = payload

			// Join the modified segments to form the new path
			newPath := strings.Join(newSegments, "/")
			reqCopy := requests.CloneReq(req)
			reqCopy.URL.Path = newPath

			// Add request body, if method is POST
			if reqCopy.Method == http.MethodPost {
				reqCopy.Body = ioutil.NopCloser(bytes.NewReader(body))
			}

			// Get raw request
			rawReq := requests.RequestToRaw(reqCopy)

			// Send request
			start := time.Now()
			resp, err := client.Do(reqCopy)
			if err != nil {
				// handle error
				if initializers.Debug {
					fmt.Println("FuzzPath:", err)
				}
				return false, "", "", "", "", ""
			}

			// Get response time
			elapsed := int(time.Since(start).Seconds())

			// Extract OOB ID
			oobID := initializers.ExtractOobID(payload)

			// Check if match vulnerability
			go detections.MatchCheck(matcher, resp, elapsed, oobID, rawReq, payload, fmt.Sprintf("segment %d", i), resultChan)
		}
	}

	// Wait for any goroutine to send a result to the channel
	for i := 0; i < len(segments)*len(payloads); i++ {
		res := <-resultChan
		if res.Found {
			return true, res.RawReq, res.URL, res.Payload, res.Param, res.RawResp
		}
	}

	return false, "", "", "", "", ""
}
