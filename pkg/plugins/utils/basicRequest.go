package utils

import (
	"net/http"
    "io/ioutil"
    "time"
)

/**
* Just reproduce the request provided
*/

func BasicRequest(r *http.Request, w http.ResponseWriter, client *http.Client) (int, string, int, map[string][]string, error) {
    req := CloneRequest(r)

    if client == nil {
        client = &http.Client{}
    }

    // Send the request
    start := time.Now()
    resp, err := client.Do(req)
    if err != nil {
        return 0, "", 0, nil, err
    }
    defer resp.Body.Close()

    // Read the response body
    responseBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return 0, "", 0, nil, err
    }

    // Get response time
    elapsed := int(time.Since(start).Seconds())

    // Get response headers
    responseHeaders := make(map[string][]string)
    for k, v := range resp.Header {
        responseHeaders[k] = v
    }

    return elapsed, string(responseBody), resp.StatusCode, responseHeaders, nil
}
