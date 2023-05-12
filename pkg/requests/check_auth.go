package requests

import (
	"net/http"
)

// Return true if request require authentication
func CheckAuth(r *http.Request, client *http.Client) bool {
	req := CloneReq(r)

	// Send Original Request
    _, _, originalStatusCode, _, _ := BasicRequest(req, client)

    // Remove auth headers
    req.Header.Del("Cookie")
    req.Header.Del("Authorization")
    
	// Send unauth request
	_, _, unauthStatusCode, _, _ := BasicRequest(req, client)

	if (originalStatusCode != unauthStatusCode) {
		return true
	} else {
		return false
	}
}