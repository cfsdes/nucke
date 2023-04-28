package fuzzers

import (
	"net/http"
	"fmt"
    "net/url"
    "strings"
    "time"
    
    "github.com/cfsdes/nucke/pkg/plugins/detections"
    "github.com/cfsdes/nucke/pkg/requests"
    internalUtils "github.com/cfsdes/nucke/internal/utils"
)

func FuzzFormData(r *http.Request, client *http.Client, payloads []string, matcher detections.Matcher, keepOriginalKey bool) (bool, string, string, string, string) {
    req := requests.CloneReq(r)

    // Update payloads {{.oob}} to interact url
    payloads = internalUtils.ReplaceOob(payloads) 
    
    // Result channel
    resultChan := make(chan detections.Result)

    // Check if method is POST and content type is application/x-www-form-urlencoded
    if req.Method != http.MethodPost || req.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
        return false, "", "", "", ""
    }

    // Get form data parameters from request body
    if err := req.ParseForm(); err != nil {
        fmt.Println(err)
        return false, "", "", "", "" 
    }

    // Get request body
    body := req.PostForm.Encode()

    // For each parameter, send a new request with the parameter replaced by a payload
    for key, values := range req.PostForm {
        for _, payload := range payloads {
            
            // Create a new request body with the parameter replaced by a payload
            var newBody string
            if keepOriginalKey {
                newBody = strings.Replace(body, fmt.Sprintf("%s=%s", key, url.QueryEscape(values[0])), fmt.Sprintf("%s=%s", key, url.QueryEscape(values[0]+payload)), -1)
            } else {
                newBody = strings.Replace(body, fmt.Sprintf("%s=%s", key, url.QueryEscape(values[0])), fmt.Sprintf("%s=%s", key, url.QueryEscape(payload)), -1)
            }

            

            // Set request body
            reqBody := strings.NewReader(newBody)

            // Create a new request with the updated form data
            newReq, err := http.NewRequest(req.Method, req.URL.String(), reqBody)
            if err != nil {
                fmt.Println(err)
                return false, "", "", "", ""
            }

            // Copy headers from original request to new request
            newReq.Header = req.Header

            // Get raw request
            rawReq := requests.RequestToRaw(newReq)

            // Send request
            start := time.Now()
            resp, err := client.Do(newReq)
            if err != nil {
                fmt.Println(err)
                return false, "", "", "", ""
            }

            // Get response time
            elapsed := int(time.Since(start).Seconds())

            // Extract OOB ID
            oobID := internalUtils.ExtractOobID(payload)

            // Check if match vulnerability
            go detections.MatchCheck(matcher, resp, elapsed, oobID, rawReq, payload, key, resultChan)
        }
    }

    // Wait for any goroutine to send a result to the channel
    for i := 0; i < len(req.PostForm)*len(payloads); i++ {
        res := <-resultChan
        if res.Found {
            return true, res.RawReq, res.URL, res.Payload, res.Param
        }
    }

    return false, "", "", "", ""
}

