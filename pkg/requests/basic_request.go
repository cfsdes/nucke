package requests

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cfsdes/nucke/pkg/globals"
)

/**
* Just reproduce the request provided
 */

func BasicRequest(r *http.Request, client *http.Client) (int, string, int, map[string][]string, string, error) {

	// Delay between requests
	time.Sleep(time.Duration(globals.Delay) * time.Millisecond)

	start := time.Now()

	// Send request
	responses := Do(r, client)
	if len(responses) == 0 {
		return 0, "", 0, nil, "", fmt.Errorf("Basic Request failed: Response nil")
	}

	// Get last array value
	lastIndex := len(responses) - 1
	lastResp := responses[lastIndex]

	// Get response body
	statusCode, responseBody, responseHeaders, rawResponse := ParseResponse(lastResp)

	// Get response time
	elapsed := int(time.Since(start).Seconds())

	return elapsed, responseBody, statusCode, responseHeaders, rawResponse, nil
}
