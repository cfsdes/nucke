package fuzzers

import (
	"net/http"
    "net/url"
    "io/ioutil"
    "bytes"
    "time"
    "fmt"
    "strings"

    "github.com/cfsdes/nucke/pkg/plugins/detections"
    "github.com/cfsdes/nucke/pkg/requests"
    "github.com/cfsdes/nucke/pkg/globals"
    "github.com/cfsdes/nucke/internal/parsers"
    "github.com/cfsdes/nucke/pkg/plugins/utils"
)

func FuzzQuery(r *http.Request, client *http.Client, payloads []string, matcher detections.Matcher) (bool, string, string, string, string, string, []detections.Result) {
    req := requests.CloneReq(r)
    
    // Extract parameters from URL
    params := req.URL.Query()

    // Result channel
    resultChan := make(chan detections.Result)

    // Array com os resultados de cada teste executado falho
    var logScans []detections.Result

    // Get request body, if method is POST
    var body []byte
    var err error
    body, err = ioutil.ReadAll(req.Body)
    if err != nil {
        // handle error
        if globals.Debug {
            fmt.Println("fuzzQuery:",err)
        }
        return false, "", "", "", "", "", nil
    }

    // For each parameter, send a new request with the parameter replaced by a payload
    for key, _ := range params {
        // Create a new query string with the parameter replaced by a payload
        for _, payload := range payloads {

            // Delay between requests
            time.Sleep(time.Duration(globals.Delay) * time.Millisecond)

            // Update payloads {{.params}}
            payload = parsers.ParsePayload(payload)

            newParams := make(url.Values)
            for k, v := range params {
                if k == key {
                    payload  = strings.Replace(payload, "{{.original}}", v[0], -1)
                    newParams.Set(k, payload)
                } else {
                    newParams.Set(k, v[0])
                }
            }

            // Copy Request
            reqCopy := requests.CloneReq(req)
            reqCopy.URL.RawQuery = newParams.Encode()

            // Add request body
            reqCopy.Body = ioutil.NopCloser(bytes.NewReader(body))

            // Get raw request
            rawReq := requests.RequestToRaw(reqCopy)

            // Send request
            start := time.Now()
            responses := requests.Do(reqCopy, client)

            // Get response time
            elapsed := int(time.Since(start).Seconds())

            // Extract OOB ID
            oobID := utils.ExtractOobID(payload)
            
            // Check if match vulnerability
            for _, resp := range responses {
                go detections.MatchCheck(matcher, resp, elapsed, oobID, rawReq, payload, key, resultChan)
            }
        }
    }

    // Wait for any goroutine to send a result to the channel
    for i := 0; i < len(params)*len(payloads); i++ {
        res := <-resultChan
        log := detections.Result{
            Found: res.Found,
            URL: res.URL,
            Payload: res.Payload,
            Param: res.Param,
            RawReq: res.RawReq,
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
