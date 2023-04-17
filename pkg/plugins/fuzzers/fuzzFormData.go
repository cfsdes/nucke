package fuzzers

import (
	"net/http"
	"fmt"
    "net/url"
    "strings"
    "time"
    
    "github.com/cfsdes/nucke/pkg/plugins/utils"
    internalUtils "github.com/cfsdes/nucke/internal/utils"
)

func FuzzFormData(r *http.Request, w http.ResponseWriter, client *http.Client, payloads []string, matcher utils.Matcher) (bool, string, string) {
    req := utils.CloneRequest(r)

    // Update payloads {{.oob}} to interact url
    payloads = internalUtils.ReplaceOob(payloads) 
    
    // Check if method is POST and content type is application/x-www-form-urlencoded
    if req.Method != http.MethodPost || req.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
        return false, "", ""
    }

    // Get form data parameters from request body
    if err := req.ParseForm(); err != nil {
        fmt.Println(err)
        return false, "", "" 
    }

    // Get request body
    body := req.PostForm.Encode()

    // For each parameter, send a new request with the parameter replaced by a payload
    for key, values := range req.PostForm {
        for _, payload := range payloads {
            // Create a new request body with the parameter replaced by a payload
            newBody := strings.Replace(body, fmt.Sprintf("%s=%s", key, url.QueryEscape(values[0])), fmt.Sprintf("%s=%s", key, url.QueryEscape(payload)), -1)

            // Set request body
            reqBody := strings.NewReader(newBody)

            // Create a new request with the updated form data
            newReq, err := http.NewRequest(req.Method, req.URL.String(), reqBody)
            if err != nil {
                fmt.Println(err)
                return false, "", ""
            }

            // Copy headers from original request to new request
            newReq.Header = req.Header

            // Get raw request
            rawReq := utils.RequestToRaw(newReq)

            // Send request
            start := time.Now()
            resp, err := client.Do(newReq)
            if err != nil {
                fmt.Println(err)
                return false, "", ""
            }
            defer resp.Body.Close()

            // Get response time
            elapsed := int(time.Since(start).Seconds())

            // Extract OOB ID
            oobID := internalUtils.ExtractOobID(payload)

            // Get URL from raw request
            url := utils.ExtractRawURL(rawReq)

            // Check if match vulnerability
            found := utils.MatchChek(matcher, resp, elapsed, oobID)
            if found {
                return true, rawReq, url
            }
        }
    }

    return false, "", ""
}

