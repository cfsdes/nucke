package fuzzers

import (
	"net/http"
    "io/ioutil"
    "bytes"

    "github.com/cfsdes/nucke/pkg/plugins/utils"
)

func FuzzHeaders(r *http.Request, w http.ResponseWriter, client *http.Client, payloads []string, headers []string, regexList []string) (bool, string, string, error) {
    req := utils.CloneRequest(r, w)

    // Update payloads {{.oob}} to interact url
    payloads = utils.ReplaceOob(payloads)
    
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
            resp, err := client.Do(req2)
            if err != nil {
                // handle error
                return false, "", "", err
            }
            defer resp.Body.Close()

            // Get response body
            body, err := ioutil.ReadAll(resp.Body)
            if err != nil {
                // handle error
                return false, "", "", err
            }

            // Check if match some regex in the list (case insensitive)
            found, err := utils.MatchString(regexList, string(body))
            if err != nil {
                return false, "", "", err
            }
            if found {
                return true, header, payload, nil
            }
        }
    }

    return false, "", "", nil
}



