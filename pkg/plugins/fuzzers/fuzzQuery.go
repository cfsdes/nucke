package fuzzers

import (
	"net/http"
    "net/url"
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

func FuzzQuery(r *http.Request, client *http.Client, payloads []string, matcher detections.Matcher) (bool, string, string, string, string) {
    req := requests.CloneReq(r)
    
    // Extract parameters from URL
    params := req.URL.Query()

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
                fmt.Println("fuzzQuery:",err)
            }
            return false, "", "", "", ""
        }
    }

    // For each parameter, send a new request with the parameter replaced by a payload
    for key, _ := range params {
        // Create a new query string with the parameter replaced by a payload
        for _, payload := range payloads {

            // Update payloads {{.params}}
            payload = parsers.ParsePayload(payload)

            newParams := make(url.Values)
            for k, v := range params {
                if k == key {
                    payload  = strings.Replace(payload, "{{.original}}", v[0], -1)
                    newParams.Set(k, payload)
                } else {
                    newParams.Set(k, v[0])
                }
            }

            // Copy Request
            reqCopy := requests.CloneReq(req)
            reqCopy.URL.RawQuery = newParams.Encode()

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
                    fmt.Println("fuzzQuery:",err)
                }
                return false, "", "", "", ""
            }
            
            // Get response time
            elapsed := int(time.Since(start).Seconds())

            // Extract OOB ID
            oobID := initializers.ExtractOobID(payload)

            // Check if match vulnerability
            go detections.MatchCheck(matcher, resp, elapsed, oobID, rawReq, payload, key, resultChan)
        }
    }

    // Wait for any goroutine to send a result to the channel
    for i := 0; i < len(params)*len(payloads); i++ {
        res := <-resultChan
        if res.Found {
            return true, res.RawReq, res.URL, res.Payload, res.Param, res.RawResp
        }
    }

    return false, "", "", "", ""
}
