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
    internalUtils "github.com/cfsdes/nucke/internal/utils"
)

func FuzzXML(r *http.Request, w http.ResponseWriter, client *http.Client, payloads []string, matcher utils.Matcher) (bool, string, string, error) {
    req := utils.CloneRequest(r, w)

    // Update payloads {{.oob}} to interact url
    payloads = internalUtils.ReplaceOob(payloads)
    
    // Check if content type is XML
    if req.Header.Get("Content-Type") != "application/xml" && req.Header.Get("Content-Type") != "text/xml" {
        return false, "", "", nil
    }

    // Get request body
    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
        return false, "", "", err
    }

    // Restore request body
    req.Body = ioutil.NopCloser(bytes.NewReader(body))

    // Find XML tags in request body and replace them with payloads
    re := regexp.MustCompile(`<([^/][^>]+)>([^<]+)</([^>]+)>`)
    matches := re.FindAllStringSubmatch(string(body), -1)
    for _, match := range matches {
        for _, payload := range payloads {
            // Create a new request body with the tag replaced by a payload
            newBody := strings.Replace(string(body), match[0], fmt.Sprintf("<%s>%s</%s>", match[1], payload, match[3]), -1)

            // Set request body
            req.Body = ioutil.NopCloser(strings.NewReader(newBody))

            // Send request
            start := time.Now()
            resp, err := client.Do(req)
            if err != nil {
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
                return true, match[1], payload, nil
            }
        }
    }

    return false, "", "", nil
}




