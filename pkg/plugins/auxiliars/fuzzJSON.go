package auxiliars

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "regexp"
)


func FuzzJSON(r *http.Request, w http.ResponseWriter, client *http.Client, payloads []string, regexList []string) (bool, string, string, error) {
    req := CreateNewRequest(r, w)
    
    // Check if method is POST and content type is application/json
    if req.Method != http.MethodPost || req.Header.Get("Content-Type") != "application/json" {
        return false, "", "", nil
    }

    // Get request body
    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
        return false, "", "", err
    }

    // For each key in the JSON object, send a new request with the key replaced by a payload
    var jsonData map[string]interface{}
    err = json.Unmarshal(body, &jsonData)
    if err != nil {
        return false, "", "", err
    }
    for key := range jsonData {
        for _, payload := range payloads {
            // Create a new JSON object with the key replaced by a payload
            newJsonData := make(map[string]interface{})
            for k, v := range jsonData {
                if k == key {
                    newJsonData[k] = payload
                } else {
                    newJsonData[k] = v
                }
            }

            // Encode the new JSON object to a byte array
            newBody, err := json.Marshal(newJsonData)
            if err != nil {
                return false, "", "", err
            }

            // Set request body
            reqBody := bytes.NewReader(newBody)

            // Create a new request with the updated JSON object
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

            // Check if response body matches any regex in the list (case insensitive)
            for _, regex := range regexList {
                match, err := regexp.MatchString("(?i)"+regex, string(respBody))
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

