package fuzzers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/cfsdes/nucke/internal/parsers"
	"github.com/cfsdes/nucke/pkg/globals"
	"github.com/cfsdes/nucke/pkg/plugins/detections"
	"github.com/cfsdes/nucke/pkg/plugins/utils"
	"github.com/cfsdes/nucke/pkg/requests"
)

func FuzzPath(r *http.Request, client *http.Client, pluginDir string, payloads []string, matcher detections.Matcher, location string) (bool, string, string, string, string, string, []detections.Result) {
	req := requests.CloneReq(r)

	// Extract segments from URL path
	segments := strings.Split(req.URL.Path, "/")

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
			responses := requests.Do(reqCopy, client)

			// Get response time
			elapsed := int(time.Since(start).Seconds())

			// Extract OOB ID
			oobID := utils.ExtractOobID(payload)

			// Check if match vulnerability
			for _, resp := range responses {
				res := detections.MatchCheck(pluginDir, matcher, resp, elapsed, oobID, rawReq, payload, fmt.Sprintf("segment %d", index))
				logScans = append(logScans, res)
			}
		}
	}

	for _, res := range logScans {
		if res.Found {
			return true, res.URL, res.Payload, res.Param, res.RawReq, res.RawResp, logScans
		}
	}

	return false, "", "", "", "", "", logScans
}
