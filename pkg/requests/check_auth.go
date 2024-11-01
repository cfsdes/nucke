package requests

import (
	"net/http"

	"github.com/cfsdes/nucke/pkg/plugins/utils"
)

// Return true if request require authentication
func CheckAuth(r *http.Request, client *http.Client) bool {
	req := CloneReq(r)

	// Send Original Request
	_, originalResBody, originalStatusCode, _, _, err := BasicRequest(req, client)
	if err != nil {
		return false
	}

	// Remove auth headers
	req.Header.Del("Cookie")
	req.Header.Del("Authorization")

	// Send unauth request
	_, unauthResBody, statusCode, _, _, err := BasicRequest(req, client)
	if err != nil {
		return false
	}

	// Compare if original response is less than 90% equal from the unauth Response
	if utils.ResponseSimilarity(originalResBody, unauthResBody) < 0.9 && statusCode != 429 && originalStatusCode < 300 {
		return true
	} else if originalStatusCode != statusCode && statusCode != 429 && originalStatusCode < 300 {
		return true
	} else {
		return false
	}
}
