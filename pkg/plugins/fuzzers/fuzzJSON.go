package fuzzers

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "regexp"

    "github.com/cfsdes/nucke/plugins/utils"
)


func FuzzJSON(r *http.Request, w http.ResponseWriter, client *http.Client, payloads []string, regexList []string) (bool, string, string, error) {
    req := utils.CloneRequest(r, w)
    
    // check if request is JSON
    if !(req.Method == http.MethodPost && req.Header.Get("Content-Type") == "application/json") {
        return false, "", "", nil
    }

    // Read request body
    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
        return false, "", "", err
    }

    // Create obj based on json data
    jsonData, err := unmarshalJSON(body)
    if err != nil {
        return false, "", "", err
    }

    // Iterate over each json object and add payload to it
    for key := range jsonData {
        for _, payload := range payloads {
            newJsonData := createNewJSONData(jsonData, key, payload)

            newBody, err := json.Marshal(newJsonData)
            if err != nil {
                return false, "", "", err
            }

            reqBody := bytes.NewReader(newBody)

            newReq, err := createNewRequest(req, reqBody)
            if err != nil {
                return false, "", "", err
            }

            resp, err := client.Do(newReq)
            if err != nil {
                return false, "", "", err
            }
            defer resp.Body.Close()

            respBody, err := ioutil.ReadAll(resp.Body)
            if err != nil {
                return false, "", "", err
            }

            if isRegexMatch(respBody, regexList) {
                return true, key, payload, nil
            }
        }
    }

    return false, "", "", nil
}


// Convert bytes to JSON
func unmarshalJSON(body []byte) (map[string]interface{}, error) {
    var jsonData map[string]interface{}
    err := json.Unmarshal(body, &jsonData)
    return jsonData, err
}


// Create new JSON object with payload
func createNewJSONData(jsonData map[string]interface{}, key string, payload string) map[string]interface{} {
    newJsonData := make(map[string]interface{})
    for k, v := range jsonData {
        if k == key {
            if m, ok := v.(map[string]interface{}); ok {
                newJsonData[k] = addPayloadToMap(m, payload)
            } else {
                newJsonData[k] = payload
            }
        } else {
            newJsonData[k] = v
        }
    }
    return newJsonData
}

// Add payload to JSON object
func addPayloadToMap(m map[string]interface{}, payload string) map[string]interface{} {
    newMap := make(map[string]interface{})
    for k, v := range m {
        if m, ok := v.(map[string]interface{}); ok {
            newMap[k] = addPayloadToMap(m, payload)
        } else {
            newMap[k] = payload
        }
    }
    return newMap
}

// Create new HTTP Request
func createNewRequest(req *http.Request, reqBody *bytes.Reader) (*http.Request, error) {
    newReq, err := http.NewRequest(req.Method, req.URL.String(), reqBody)
    if err != nil {
        return nil, err
    }
    newReq.Header = req.Header
    return newReq, nil
}

// Check if response match regex
func isRegexMatch(respBody []byte, regexList []string) bool {
    for _, regex := range regexList {
        match, err := regexp.MatchString("(?i)"+regex, string(respBody))
        if err != nil {
            return false
        }
        if match {
            return true
        }
    }
    return false
}

