package requests

import (
	"net/http"
	"github.com/cfsdes/nucke/pkg/plugins/utils"
)

// Return true if request require authentication
func CheckAuth(r *http.Request, client *http.Client) bool {
	req := CloneReq(r)

	// Send Original Request
    _, originalResBody, _, _, _ := BasicRequest(req, client)

    // Remove auth headers
    req.Header.Del("Cookie")
    req.Header.Del("Authorization")
    
	// Send unauth request
	_, unauthResBody, _, _, _ := BasicRequest(req, client)

	// Compare if original response is less than 95% equal from the unauth Response
	if (utils.TextSimilarity(originalResBody, unauthResBody) <= 0.95) {
		return true
	} else {
		return false
	}
}