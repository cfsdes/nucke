package fuzzers

import (
	"net/http"
    "io/ioutil"
    "bytes"
    "time"
    "fmt"

    "github.com/cfsdes/nucke/pkg/plugins/detections"
    "github.com/cfsdes/nucke/pkg/requests"
    "github.com/cfsdes/nucke/internal/initializers"
)

func FuzzHeaders(r *http.Request, client *http.Client, payloads []string, headers []string, matcher detections.Matcher, keepOriginalKey bool) (bool, string, string, string, string) {
    req := requests.CloneReq(r)

    // Update payloads {{.oob}} to interact url
    payloads = initializers.ReplaceOob(payloads)
    
    // Result channel
    resultChan := make(chan detections.Result)

    // Get request body, if method is POST
    var body []byte
    if req.Method == http.MethodPost {
        var err error
        body, err = ioutil.ReadAll(req.Body)
        if err != nil {
            // handle error
            fmt.Println(err)
            return false, "", "", "", ""
        }
    }

    // For each header, send a new request with the header replaced by a payload
    for _, header := range headers {
        // Create a new request with the header replaced by a payload
        for _, payload := range payloads {
            req2 := requests.CloneReq(req)

            if keepOriginalKey {
                currentValue := req.Header.Get(header)
                req2.Header.Set(header, currentValue+payload)
            } else {
                req2.Header.Set(header, payload)
            }
            

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
                fmt.Println(err)
                return false, "", "", "", ""
            }

            // Get response time
            elapsed := int(time.Since(start).Seconds())

            // Extract OOB ID
            oobID := initializers.ExtractOobID(payload)

            // Check if match vulnerability
            go detections.MatchCheck(matcher, resp, elapsed, oobID, rawReq, payload, header, resultChan)
        }
    }

    // Wait for any goroutine to send a result to the channel
    for i := 0; i < len(headers)*len(payloads); i++ {
        res := <-resultChan
        if res.Found {
            return true, res.RawReq, res.URL, res.Payload, res.Param
        }
    }

    return false, "", "", "", ""
}



