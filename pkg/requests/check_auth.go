package requests

import (
	"net/http"
	"fmt"

	"github.com/cfsdes/nucke/pkg/globals"
	"github.com/cfsdes/nucke/pkg/plugins/utils"
)

// Return true if request require authentication
func CheckAuth(r *http.Request, client *http.Client) bool {
	req := CloneReq(r)

	// Send Original Request
    _, originalResBody, originalStatusCode, _, _ := BasicRequest(req, client)

    // Remove auth headers
    req.Header.Del("Cookie")
    req.Header.Del("Authorization")
    
	// Send unauth request
	resp, err := client.Do(req)
    if err != nil {
        if globals.Debug {
            fmt.Println("CheckAuth Request Error:",err)
        }
    }

	// Get response body
    statusCode, unauthResBody, _, _ := ParseResponse(resp)

	// Compare if original response is less than 90% equal from the unauth Response
	if (utils.TextSimilarity(originalResBody, unauthResBody) < 0.9 && statusCode != 429) {
		return true
	} else if (originalStatusCode != statusCode && statusCode != 429) {
		return true
	} else {
		return false
	}
}