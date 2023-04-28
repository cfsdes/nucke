package utils

import (
	"net/http"
    "io/ioutil"
    "time"
    "fmt"
)

/**
* Just reproduce the request provided
*/

func BasicRequest(r *http.Request, client *http.Client) (int, string, int, map[string][]string) {
    req := CloneRequest(r)

    if client == nil {
        client = &http.Client{}
    }

    // Send the request
    start := time.Now()
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
        return 0, "", 0, nil
    }
    defer resp.Body.Close()

    // Read the response body
    responseBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err)
        return 0, "", 0, nil
    }

    // Get response time
    elapsed := int(time.Since(start).Seconds())

    // Get response headers
    responseHeaders := make(map[string][]string)
    for k, v := range resp.Header {
        responseHeaders[k] = v
    }

    return elapsed, string(responseBody), resp.StatusCode, responseHeaders
}
