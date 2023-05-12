package fuzzers

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "time"
    "fmt"
    "strings"

    "github.com/cfsdes/nucke/pkg/plugins/detections"
    "github.com/cfsdes/nucke/pkg/requests"
    "github.com/cfsdes/nucke/internal/initializers"
    "github.com/cfsdes/nucke/internal/parsers"
)


func FuzzJSON(r *http.Request, client *http.Client, payloads []string, matcher detections.Matcher) (bool, string, string, string, string) {
    req := requests.CloneReq(r)

    // Result channel
    resultChan := make(chan detections.Result)
    
    // check if request is JSON
    if !(req.Method == http.MethodPost && req.Header.Get("Content-Type") == "application/json") {
        return false, "", "", "", ""
    }

    // Read request body
    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
        if initializers.Debug {
            fmt.Println("fuzzJSON:",err)
        }
        return false, "", "", "", ""
    }

    // Create obj based on json data
    jsonData, err := unmarshalJSON(body)
    if err != nil {
        if initializers.Debug {
            fmt.Println("fuzzJSON:",err)
        }
        return false, "", "", "", ""
    }

    for key, value := range jsonData {
        for _, payload := range payloads {
            // Update payloads {{.params}}
            payload = parsers.ParsePayload(payload)
            
            // Check if value is map. If yes, recursively check it to inject payload
            addPayloadToJson(jsonData, key, value, payload, resultChan, req, client, matcher)
        }
    }

    // Wait for any goroutine to send a result to the channel
    for i := 0; i < len(jsonData)*len(payloads); i++ {
        res := <-resultChan
        if res.Found {
            return true, res.RawReq, res.URL, res.Payload, res.Param, res.RawResp
        }
    }

    return false, "", "", "", ""
}

// function to add payload to JSON
func addPayloadToJson(jsonData map[string]interface{}, key string, value interface{}, payload string, resultChan chan detections.Result, req *http.Request, client *http.Client, matcher detections.Matcher) {
    if innerMap, ok := value.(map[string]interface{}); ok {
        // Se for um mapa, iterar sobre suas chaves e valores
        for innerKey, innerValue := range innerMap {
            addPayloadToJson(jsonData, innerKey, innerValue, payload, resultChan, req, client, matcher)
        }
    } else {
        loopScan(jsonData, key, payload, resultChan, req, client, matcher)
    }
}

// Scan to send request and check match
func loopScan(jsonData map[string]interface{}, key string, payload string, resultChan chan detections.Result, req *http.Request, client *http.Client, matcher detections.Matcher) {
    // Iterate over each json object and add payload to it
    newJsonData := createNewJSONData(jsonData, key, payload)

    newBody, err := json.Marshal(newJsonData)
    if err != nil {
        if initializers.Debug {
            fmt.Println("fuzzJSON:",err)
        }
    }

    reqBody := bytes.NewReader(newBody)

    newReq, err := createNewRequest(req, reqBody)
    if err != nil {
        if initializers.Debug {
            fmt.Println("fuzzJSON:",err)
        }
    }

    // Get raw request
    rawReq := requests.RequestToRaw(newReq)

    // Make request
    start := time.Now()
    resp, err := client.Do(newReq)
    if err != nil {
        if initializers.Debug {
            fmt.Println("fuzzJSON:",err)
        }
    }

    // Get response time
    elapsed := int(time.Since(start).Seconds())

    // Extract OOB ID
    oobID := initializers.ExtractOobID(payload)

    // Check if match vulnerability
    go detections.MatchCheck(matcher, resp, elapsed, oobID, rawReq, payload, key, resultChan)
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
                newJsonData[k] = createNewJSONData(m, key, payload)
            } else {
                originalValue := fmt.Sprintf("%v", v)
                payload  = strings.Replace(payload, "{{.original}}", originalValue, -1)
                newJsonData[k] = payload
            }
        } else {
            if m, ok := v.(map[string]interface{}); ok {
                newJsonData[k] = createNewJSONData(m, key, payload)
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


