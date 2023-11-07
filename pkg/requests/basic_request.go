package requests

import (
	"net/http"
    "time"
)

/**
* Just reproduce the request provided
*/

func BasicRequest(r *http.Request, client *http.Client) (int, string, int, map[string][]string, string) {
    
    start := time.Now()

    // Send request
    responses := Do(r, client)

    // Get last array value
	lastIndex := len(responses) - 1
	lastResp := responses[lastIndex]


    // Get response body
    statusCode, responseBody, responseHeaders, rawResponse := ParseResponse(lastResp)
	
    // Get response time
    elapsed := int(time.Since(start).Seconds())

    return elapsed, responseBody, statusCode, responseHeaders, rawResponse
}
