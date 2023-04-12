package auxiliars

import (
	"net/http"
    "net/url"
    "io/ioutil"
    "bytes"
    "regexp"
)

func FuzzQuery(r *http.Request, w http.ResponseWriter, client *http.Client, payloads []string, regexList []string) (bool, string, string, error) {
    req := CreateNewRequest(r, w)
    
    // Extract parameters from URL
    params := req.URL.Query()

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

            // Send request
            resp, err := client.Do(req)
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
            for _, regex := range regexList {
                match, err := regexp.MatchString("(?i)"+regex, string(body))
                if err != nil {
                    return false, "", "", err
                }
                if match {
                    return true, key, payload, nil
                }
            }
        }
    }

    return false, "", "", nil
}
