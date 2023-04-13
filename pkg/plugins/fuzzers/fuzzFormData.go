package fuzzers

import (
	"net/http"
	"fmt"
    "net/url"
    "io/ioutil"
    "strings"
    
    "github.com/cfsdes/nucke/pkg/plugins/utils"
)

func FuzzFormData(r *http.Request, w http.ResponseWriter, client *http.Client, payloads []string, regexList []string) (bool, string, string, error) {
    req := utils.CloneRequest(r, w)

    // Update payloads {{.oob}} to interact url
    payloads = utils.ReplaceOob(payloads) 
    
    // Check if method is POST and content type is application/x-www-form-urlencoded
    if req.Method != http.MethodPost || req.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
        return false, "", "", nil
    }

    // Get form data parameters from request body
    if err := req.ParseForm(); err != nil {
        return false, "", "", err
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
                return false, "", "", err
            }

            // Copy headers from original request to new request
            newReq.Header = req.Header

            // Send request
            resp, err := client.Do(newReq)
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
                return true, key, payload, nil
            }
        }
    }

    return false, "", "", nil
}

