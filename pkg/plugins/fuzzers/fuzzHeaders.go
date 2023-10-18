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

func FuzzHeaders(r *http.Request, client *http.Client, payloads []string, headers []string, matcher detections.Matcher, behavior string) (bool, string, string, string, string, string, []detections.Result) {
    req := requests.CloneReq(r)

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
            fmt.Println("fuzzHeaders:", err)
        }
        return false, "", "", "", "", "", nil
    }

    totalResults := 0

    if behavior == "all" {
        totalResults = len(payloads)
        // Create a new request for each payload
        for _, payload := range payloads {

            // Delay between requests
            time.Sleep(time.Duration(globals.Delay) * time.Millisecond)

            // Update payloads {{.params}}
            payload = parsers.ParsePayload(payload)

            req2 := requests.CloneReq(req)

            // Inject payload into all headers
            for _, header := range headers {
                currentValue := req.Header.Get(header)
                payload = strings.Replace(payload, "{{.original}}", currentValue, -1)
                req2.Header.Set(header, payload)
            }

            // Add request body
            req2.Body = ioutil.NopCloser(bytes.NewReader(body))

            // Get raw request
            rawReq := requests.RequestToRaw(req2)

            // Send request
            start := time.Now()
            responses := requests.Do(req2, client)

            // Get response time
            elapsed := int(time.Since(start).Seconds())

            // Extract OOB ID
            oobID := utils.ExtractOobID(payload)

            // Check if match vulnerability
            for _, resp := range responses {
                go detections.MatchCheck(matcher, resp, elapsed, oobID, rawReq, payload, strings.Join(headers, ","), resultChan)
            }
        }
    } else {
        totalResults = len(headers) * len(payloads)
        // Inject one payload at a time into each header
        for _, header := range headers {
            // Create a new request for each header and payload
            for _, payload := range payloads {

                // Delay between requests
                time.Sleep(time.Duration(globals.Delay) * time.Millisecond)

                // Update payloads {{.params}}
                payload = parsers.ParsePayload(payload)

                req2 := requests.CloneReq(req)

                currentValue := req.Header.Get(header)
                payload = strings.Replace(payload, "{{.original}}", currentValue, -1)
                req2.Header.Set(header, payload)

                // Add request body, if method is POST
                if req2.Method == http.MethodPost {
                    req2.Body = ioutil.NopCloser(bytes.NewReader(body))
                }

                // Get raw request
                rawReq := requests.RequestToRaw(req2)

                // Send request
                start := time.Now()
                responses := requests.Do(req2, client)

                // Get response time
                elapsed := int(time.Since(start).Seconds())

                // Extract OOB ID
                oobID := utils.ExtractOobID(payload)

                // Check if match vulnerability
                for _, resp := range responses {
                    go detections.MatchCheck(matcher, resp, elapsed, oobID, rawReq, payload, header, resultChan)
                }
            }
        }
    }

    // Wait for the expected number of results from goroutines
    for i := 0; i < totalResults; i++ {
        res := <-resultChan
        log := detections.Result{
            Found: res.Found,
            RawReq: res.RawReq,
            URL: res.URL,
            Payload: res.Payload,
            Param: res.Param,
            RawResp: res.RawResp,
            ResBody: res.ResBody,
        }
        logScans = append(logScans, log)
    }

    for _, res := range logScans {
		if res.Found {
			return true, res.RawReq, res.URL, res.Payload, res.Param, res.RawResp, logScans
		}
	}

    return false, "", "", "", "", "", logScans
}
