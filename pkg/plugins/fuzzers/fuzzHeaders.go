package fuzzers

import (
	"net/http"
    "io/ioutil"
    "bytes"
    "time"

    "github.com/cfsdes/nucke/pkg/plugins/utils"
    internalUtils "github.com/cfsdes/nucke/internal/utils"
)

func FuzzHeaders(r *http.Request, w http.ResponseWriter, client *http.Client, payloads []string, headers []string, matcher utils.Matcher) (bool, string, string, error) {
    req := utils.CloneRequest(r, w)

    // Update payloads {{.oob}} to interact url
    payloads = internalUtils.ReplaceOob(payloads)
    
    // Get request body, if method is POST
    var body []byte
    if req.Method == http.MethodPost {
        var err error
        body, err = ioutil.ReadAll(req.Body)
        if err != nil {
            // handle error
            return false, "", "", err
        }
    }

    // For each header, send a new request with the header replaced by a payload
    for _, header := range headers {
        // Create a new request with the header replaced by a payload
        for _, payload := range payloads {
            req2 := utils.CloneRequest(req, w)
            req2.Header.Set(header, payload)

            // Add request body, if method is POST
            if req2.Method == http.MethodPost {
                req2.Body = ioutil.NopCloser(bytes.NewReader(body))
            }

            // Send request
            start := time.Now()
            resp, err := client.Do(req2)
            if err != nil {
                // handle error
                return false, "", "", err
            }
            defer resp.Body.Close()

            // Get response time
            elapsed := int(time.Since(start).Seconds())

            // Extract OOB ID
            oobID := internalUtils.ExtractOobID(payload)

            // Check if match vulnerability
            found := utils.MatchChek(matcher, resp, elapsed, oobID)
            if found {
                return true, header, payload, nil
            }
        }
    }

    return false, "", "", nil
}



