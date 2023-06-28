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
    "github.com/cfsdes/nucke/internal/globals"
    "github.com/cfsdes/nucke/internal/parsers"
    "github.com/cfsdes/nucke/pkg/plugins/utils"
)

func FuzzHeaders(r *http.Request, client *http.Client, payloads []string, headers []string, matcher detections.Matcher) (bool, string, string, string, string, string) {
    req := requests.CloneReq(r)

    // Result channel
    resultChan := make(chan detections.Result)

    // Get request body, if method is POST
    var body []byte
    if req.Method == http.MethodPost {
        var err error
        body, err = ioutil.ReadAll(req.Body)
        if err != nil {
            // handle error
            if globals.Debug {
                fmt.Println("fuzzHeaders:",err)
            }
            return false, "", "", "", "", ""
        }
    }

    // For each header, send a new request with the header replaced by a payload
    for _, header := range headers {
        // Create a new request with the header replaced by a payload
        for _, payload := range payloads {
            
            // Update payloads {{.params}}
            payload = parsers.ParsePayload(payload)

            req2 := requests.CloneReq(req)

            currentValue := req.Header.Get(header)
            payload  = strings.Replace(payload, "{{.original}}", currentValue, -1)
            req2.Header.Set(header, payload)
            
            // Add request body, if method is POST
            if req2.Method == http.MethodPost {
                req2.Body = ioutil.NopCloser(bytes.NewReader(body))
            }

            // Get raw request
            rawReq := requests.RequestToRaw(req2)

            // Send request
            start := time.Now()
            resp, err := client.Do(req2)
            if err != nil {
                // handle error
                if globals.Debug {
                    fmt.Println("fuzzHeaders:",err)
                }
                return false, "", "", "", "", ""
            }

            // Get response time
            elapsed := int(time.Since(start).Seconds())

            // Extract OOB ID
            oobID := utils.ExtractOobID(payload)

            // Check if match vulnerability
            go detections.MatchCheck(matcher, resp, elapsed, oobID, rawReq, payload, header, resultChan)
        }
    }

    // Wait for any goroutine to send a result to the channel
    for i := 0; i < len(headers)*len(payloads); i++ {
        res := <-resultChan
        if res.Found {
            return true, res.RawReq, res.URL, res.Payload, res.Param, res.RawResp
        }
    }

    return false, "", "", "", "", ""
}



