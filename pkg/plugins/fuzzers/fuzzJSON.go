package fuzzers

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "time"
    "fmt"

    "github.com/cfsdes/nucke/pkg/plugins/utils"
    internalUtils "github.com/cfsdes/nucke/internal/utils"
)


func FuzzJSON(r *http.Request, w http.ResponseWriter, client *http.Client, payloads []string, matcher utils.Matcher, keepOriginalKey bool) (bool, string, string, string, string) {
    req := utils.CloneRequest(r)

    // Result channel
    resultChan := make(chan utils.Result)
    
    // Update payloads {{.oob}} to interact url
    payloads = internalUtils.ReplaceOob(payloads)
    
    // check if request is JSON
    if !(req.Method == http.MethodPost && req.Header.Get("Content-Type") == "application/json") {
        return false, "", "", "", ""
    }

    // Read request body
    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
        fmt.Println(err)
        return false, "", "", "", ""
    }

    // Create obj based on json data
    jsonData, err := unmarshalJSON(body)
    if err != nil {
        fmt.Println(err)
        return false, "", "", "", ""
    }

    for key, value := range jsonData {
        for _, payload := range payloads {
            
            // Check if value is map. If yes, recursively check it to inject payload
            addPayloadToJson(jsonData, key, value, payload, resultChan, req, client, matcher, keepOriginalKey)
        }
    }

    // Wait for any goroutine to send a result to the channel
    for i := 0; i < len(jsonData)*len(payloads); i++ {
        res := <-resultChan
        if res.Found {
            return true, res.RawReq, res.URL, res.Payload, res.Param
        }
    }

    return false, "", "", "", ""
}

// function to add payload to JSON
func addPayloadToJson(jsonData map[string]interface{}, key string, value interface{}, payload string, resultChan chan utils.Result, req *http.Request, client *http.Client, matcher utils.Matcher, keepOriginalKey bool) {
    if innerMap, ok := value.(map[string]interface{}); ok {
        // Se for um mapa, iterar sobre suas chaves e valores
        for innerKey, innerValue := range innerMap {
            addPayloadToJson(jsonData, innerKey, innerValue, payload, resultChan, req, client, matcher, keepOriginalKey)
        }
    } else {
        loopScan(jsonData, key, payload, resultChan, req, client, matcher, keepOriginalKey)
    }
}

// Scan to send request and check match
func loopScan(jsonData map[string]interface{}, key string, payload string, resultChan chan utils.Result, req *http.Request, client *http.Client, matcher utils.Matcher, keepOriginalKey bool) {
    // Iterate over each json object and add payload to it
    newJsonData := createNewJSONData(jsonData, key, payload, keepOriginalKey)

    newBody, err := json.Marshal(newJsonData)
    if err != nil {
        fmt.Println(err)
    }

    reqBody := bytes.NewReader(newBody)

    newReq, err := createNewRequest(req, reqBody)
    if err != nil {
        fmt.Println(err)
    }

    // Get raw request
    rawReq := utils.RequestToRaw(newReq)

    // Make request
    start := time.Now()
    resp, err := client.Do(newReq)
    if err != nil {
        fmt.Println(err)
    }

    // Get response time
    elapsed := int(time.Since(start).Seconds())

    // Extract OOB ID
    oobID := internalUtils.ExtractOobID(payload)

    // Check if match vulnerability
    go utils.MatchChek(matcher, resp, elapsed, oobID, rawReq, payload, key, resultChan)
}

// Convert bytes to JSON
func unmarshalJSON(body []byte) (map[string]interface{}, error) {
    var jsonData map[string]interface{}
    err := json.Unmarshal(body, &jsonData)
    return jsonData, err
}


// Create new JSON object with payload
func createNewJSONData(jsonData map[string]interface{}, key string, payload string, keepOriginalKey bool) map[string]interface{} {
    newJsonData := make(map[string]interface{})
    for k, v := range jsonData {
        if k == key {
            if m, ok := v.(map[string]interface{}); ok {
                newJsonData[k] = createNewJSONData(m, key, payload, keepOriginalKey)
            } else {
                if keepOriginalKey {
                    originalValue := fmt.Sprintf("%v", v)
                    newJsonData[k] = originalValue+payload
                } else {
                    newJsonData[k] = payload
                }
            }
        } else {
            if m, ok := v.(map[string]interface{}); ok {
                newJsonData[k] = createNewJSONData(m, key, payload, keepOriginalKey)
            } else {
                newJsonData[k] = v
            }
        }
    }
    return newJsonData
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


