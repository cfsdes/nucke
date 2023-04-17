package fuzzers

import (
	"net/http"
    "net/url"
    "io/ioutil"
    "bytes"
    "time"
    "fmt"

    "github.com/cfsdes/nucke/pkg/plugins/utils"
    internalUtils "github.com/cfsdes/nucke/internal/utils"
)

func FuzzQuery(r *http.Request, w http.ResponseWriter, client *http.Client, payloads []string, matcher utils.Matcher) (bool, string) {
    req := utils.CloneRequest(r)
    
    // Update payloads {{.oob}} to interact url
    payloads = internalUtils.ReplaceOob(payloads)
    
    // Extract parameters from URL
    params := req.URL.Query()

    // Get request body, if method is POST
    var body []byte
    if req.Method == http.MethodPost {
        var err error
        body, err = ioutil.ReadAll(req.Body)
        if err != nil {
            // handle error
            fmt.Println(err)
            return false, ""
        }
    }

    // For each parameter, send a new request with the parameter replaced by a payload
    for key, _ := range params {
        // Create a new query string with the parameter replaced by a payload
        for _, payload := range payloads {
            newParams := make(url.Values)
            for k, v := range params {
                if k == key {
                    newParams.Set(k, payload)
                } else {
                    newParams.Set(k, v[0])
                }
            }
            req.URL.RawQuery = newParams.Encode()

            // Add request body, if method is POST
            if req.Method == http.MethodPost {
                req.Body = ioutil.NopCloser(bytes.NewReader(body))
            }

            // Get raw request
            rawReq := utils.RequestToRaw(req)

            // Send request
            start := time.Now()
            resp, err := client.Do(req)
            if err != nil {
                // handle error
                fmt.Println(err)
                return false, ""
            }
            defer resp.Body.Close()

            // Get response time
            elapsed := int(time.Since(start).Seconds())

            // Extract OOB ID
            oobID := internalUtils.ExtractOobID(payload)

            // Check if match vulnerability
            found := utils.MatchChek(matcher, resp, elapsed, oobID)
            if found {
                fmt.Println(err)
                return true, rawReq
            }
        }
    }

    return false, ""
}
