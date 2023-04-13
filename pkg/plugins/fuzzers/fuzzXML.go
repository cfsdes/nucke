package fuzzers

import (
	"net/http"
    "io/ioutil"
    "bytes"
    "regexp"
    "fmt"
    "strings"

    "github.com/cfsdes/nucke/pkg/plugins/utils"
)

func FuzzXML(r *http.Request, w http.ResponseWriter, client *http.Client, payloads []string, regexList []string) (bool, string, string, error) {
    req := utils.CloneRequest(r, w)

    // Update payloads {{.oob}} to interact url
    payloads = utils.ReplaceOob(payloads)
    
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
            resp, err := client.Do(req)
            if err != nil {
                return false, "", "", err
            }
            defer resp.Body.Close()

            // Get response body
            respBody, err := ioutil.ReadAll(resp.Body)
            if err != nil {
                return false, "", "", err
            }

            // Check if match some regex in the list (case insensitive)
            found, err := utils.MatchString(regexList, string(respBody))
            if err != nil {
                return false, "", "", err
            }
            if found {
                return true, match[1], payload, nil
            }
        }
    }

    return false, "", "", nil
}




