package requests

import (
	"net/http"
    "time"
    "fmt"

    "github.com/cfsdes/nucke/internal/initializers"
)

/**
* Just reproduce the request provided
*/

func BasicRequest(r *http.Request, client *http.Client) (int, string, int, map[string][]string) {
    req := CloneReq(r)

    if client == nil {
        client = &http.Client{}
    }

    // Send the request
    start := time.Now()
    resp, err := client.Do(req)
    if err != nil {
        if initializers.Debug {
            fmt.Println("Basic Request Error:",err)
        }
        return 0, "", 0, nil
    }

    // Get response body
    statusCode, responseBody, responseHeaders, _ := ParseResponse(resp)
	
    // Get response time
    elapsed := int(time.Since(start).Seconds())

    return elapsed, responseBody, statusCode, responseHeaders
}
