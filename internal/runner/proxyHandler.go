package runner

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
    "regexp"
    "fmt"
)

// Proxy Handler
func Handler(w http.ResponseWriter, r *http.Request, jaelesApi string, scope string) {
	// Convert the raw request to base64
	requestBytes, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	requestBase64 := base64.StdEncoding.EncodeToString(requestBytes)

    // Send request to jaeles API server and filter if scope is specified
    if scope != "" && regexp.MustCompile(scope).MatchString(r.URL.String()) {
		SendToJaeles(requestBase64, jaelesApi)
	} else if scope == "" {
        SendToJaeles(requestBase64, jaelesApi)
    }

    fowardRequest(w, r)
	
}

func fowardRequest(w http.ResponseWriter, r *http.Request) {
    // Forward the request to the destination server
    resp, err := http.DefaultTransport.RoundTrip(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadGateway)
        return
    }
    defer resp.Body.Close()

    // Copy the response headers
    for key, values := range resp.Header {
        for _, value := range values {
            w.Header().Add(key, value)
        }
    }

    // Copy the response body and return it to the client
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Write(body)
}