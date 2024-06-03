package fuzzers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/cfsdes/nucke/internal/parsers"
	"github.com/cfsdes/nucke/pkg/globals"
	"github.com/cfsdes/nucke/pkg/plugins/detections"
	"github.com/cfsdes/nucke/pkg/plugins/utils"
	"github.com/cfsdes/nucke/pkg/requests"
	"github.com/fatih/color"
)

// Counter of channels opened
var channelsOpened int

func FuzzJSON(r *http.Request, client *http.Client, pluginDir string, payloads []string, matcher detections.Matcher) (bool, string, string, string, string, string, []detections.Result) {
	req := requests.CloneReq(r)

	// Result channel
	resultChan := make(chan detections.Result)

	// Array com os resultados de cada teste executado falho
	var logScans []detections.Result

	// check if request is JSON
	if !(req.Method == http.MethodPost && req.Header.Get("Content-Type") == "application/json") {
		return false, "", "", "", "", "", nil
	}

	// Read request body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		if globals.Debug {
			fmt.Println("fuzzJSON:", err)
		}
		return false, "", "", "", "", "", nil
	}

	// Create obj based on json data
	jsonData, err := unmarshalJSON(body)
	if err != nil {
		if globals.Debug {
			fmt.Println("fuzzJSON:", err)
		}
		return false, "", "", "", "", "", nil
	}

	for key, value := range jsonData {
		for _, payload := range payloads {

			// Check if value is map. If yes, recursively check it to inject payload
			addPayloadToJson(jsonData, key, value, payload, resultChan, req, client, matcher, pluginDir)
		}
	}

	// Wait for any goroutine to send a result to the channel
	for i := 0; i < channelsOpened; i++ {
		res := <-resultChan
		log := detections.Result{
			Found:   res.Found,
			URL:     res.URL,
			Payload: res.Payload,
			Param:   res.Param,
			RawReq:  res.RawReq,
			RawResp: res.RawResp,
			ResBody: res.ResBody,
		}
		logScans = append(logScans, log)
	}

	for _, res := range logScans {
		if res.Found {
			return true, res.URL, res.Payload, res.Param, res.RawReq, res.RawResp, logScans
		}
	}

	return false, "", "", "", "", "", logScans
}

// function to add payload to JSON
func addPayloadToJson(jsonData map[string]interface{}, key string, value interface{}, payload string, resultChan chan detections.Result, req *http.Request, client *http.Client, matcher detections.Matcher, pluginDir string) {
	if innerMap, ok := value.(map[string]interface{}); ok {
		// Se for um mapa, iterar sobre suas chaves e valores
		for innerKey, innerValue := range innerMap {
			addPayloadToJson(jsonData, innerKey, innerValue, payload, resultChan, req, client, matcher, pluginDir)
		}
	} else if innerArray, ok := value.([]interface{}); ok {
		// Se for um array, não faz nada
		payload = parsers.ParsePayload(payload)
		newInnerArray := append(innerArray, payload)
		updatedArrayStr, err := json.Marshal(newInnerArray)
		if err != nil && globals.Debug {
			// Print error
			Red := color.New(color.FgRed, color.Bold).SprintFunc()
			fmt.Printf("[%s] Error FuzzJSON: %v\n", Red("ERR"), err)
		}

		loopScan(jsonData, key, string(updatedArrayStr), resultChan, req, client, matcher, pluginDir)

	} else {
		// Update payloads {{.params}}
		payload = parsers.ParsePayload(payload)
		loopScan(jsonData, key, payload, resultChan, req, client, matcher, pluginDir)
	}
}

// Scan to send request and check match
func loopScan(jsonData map[string]interface{}, key string, payload string, resultChan chan detections.Result, req *http.Request, client *http.Client, matcher detections.Matcher, pluginDir string) {

	// Delay between requests
	time.Sleep(time.Duration(globals.Delay) * time.Millisecond)

	// Iterate over each json object and add payload to it
	newJsonData := createNewJSONData(jsonData, key, payload)

	newBody, err := json.Marshal(newJsonData)
	if err != nil {
		if globals.Debug {
			fmt.Println("fuzzJSON:", err)
		}
	}

	reqBody := bytes.NewReader(newBody)

	newReq, err := createNewRequest(req, reqBody)
	if err != nil {
		if globals.Debug {
			fmt.Println("fuzzJSON:", err)
		}
	}

	// Get raw request
	rawReq := requests.RequestToRaw(newReq)

	// Make request
	start := time.Now()
	responses := requests.Do(newReq, client)

	// Get response time
	elapsed := int(time.Since(start).Seconds())

	// Extract OOB ID
	oobID := utils.ExtractOobID(payload)

	// Check if match vulnerability
	for _, resp := range responses {
		channelsOpened++
		go detections.MatchCheck(pluginDir, matcher, resp, elapsed, oobID, rawReq, payload, key, resultChan)
	}
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
			// Se o valor da key do json for um obj, iterar novamente
			if m, ok := v.(map[string]interface{}); ok {
				newJsonData[k] = createNewJSONData(m, key, payload)

				// Se o valor da key do json for um array, adicionar o payload no array
			} else if _, isArray := v.([]interface{}); isArray {
				var result []interface{}
				payload = strings.Replace(payload, "{{.original}}", "", -1)
				if err := json.Unmarshal([]byte(payload), &result); err != nil {
					fmt.Println("Erro ao decodificar a string JSON:", err)
				}
				newJsonData[k] = result

				// Se o valor da key do json for uma string, substituir com o payload
			} else {
				originalValue := fmt.Sprintf("%v", v)
				payload = strings.Replace(payload, "{{.original}}", originalValue, -1)
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

// Count the total number of parameters recursively in the JSON data
func countParams(jsonData map[string]interface{}) int {
	count := 0
	for _, value := range jsonData {
		if innerMap, ok := value.(map[string]interface{}); ok {
			count += countParams(innerMap)
		} else {
			count++
		}
	}
	return count
}
