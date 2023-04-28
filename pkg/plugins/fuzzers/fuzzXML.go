package fuzzers

import (
	"net/http"
    "io/ioutil"
    "bytes"
    "regexp"
    "fmt"
    "strings"
    "time"

    "github.com/cfsdes/nucke/pkg/plugins/utils"
    "github.com/cfsdes/nucke/pkg/requests"
    internalUtils "github.com/cfsdes/nucke/internal/utils"
)

func FuzzXML(r *http.Request, client *http.Client, payloads []string, matcher utils.Matcher, keepOriginalKey bool) (bool, string, string, string, string) {
    req := requests.CloneReq(r)

    // Result channel
    resultChan := make(chan utils.Result)

    // Update payloads {{.oob}} to interact url
    payloads = internalUtils.ReplaceOob(payloads)
    
    // Check if content type is XML
    if req.Header.Get("Content-Type") != "application/xml" && req.Header.Get("Content-Type") != "text/xml" {
        return false, "", "", "", ""
    }

    // Get request body
    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
        fmt.Println(err)
        return false, "", "", "", ""
    }

    // Restore request body
    req.Body = ioutil.NopCloser(bytes.NewReader(body))

    // Find XML tags in request body and replace them with payloads
    re := regexp.MustCompile(`<([^/][^>]+)>([^<]+)</([^>]+)>`)
    matches := re.FindAllStringSubmatch(string(body), -1)
    for _, match := range matches {
        for _, payload := range payloads {
            // Copy Request
            reqCopy := requests.CloneReq(req)

            // Create a new request body with the tag replaced by a payload
            var newBody string
            if keepOriginalKey {
                newBody = strings.Replace(string(body), match[0], fmt.Sprintf("<%s>%s</%s>", match[1], match[2]+payload, match[3]), -1)
            } else {
                newBody = strings.Replace(string(body), match[0], fmt.Sprintf("<%s>%s</%s>", match[1], payload, match[3]), -1)
            }
            
            

            // Set request body
            reqCopy.Body = ioutil.NopCloser(strings.NewReader(newBody))

            // Get raw request
            rawReq := requests.RequestToRaw(reqCopy)

            // Send request
            start := time.Now()
            resp, err := client.Do(reqCopy)
            if err != nil {
                fmt.Println(err)
                return false, "", "", "", ""
            }

            // Get response time
            elapsed := int(time.Since(start).Seconds())

            // Extract OOB ID
            oobID := internalUtils.ExtractOobID(payload)

            // Check if match vulnerability
            go utils.MatchChek(matcher, resp, elapsed, oobID, rawReq, payload, match[0], resultChan)
        }
    }

    // Wait for any goroutine to send a result to the channel
    for i := 0; i < len(matches)*len(payloads); i++ {
        res := <-resultChan
        if res.Found {
            return true, res.RawReq, res.URL, res.Payload, res.Param
        }
    }

    return false, "", "", "", ""
}




