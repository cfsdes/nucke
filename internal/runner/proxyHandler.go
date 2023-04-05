package runner

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
    "regexp"
    "fmt"
    "log"

    "github.com/fatih/color"
    "github.com/cfsdes/nucke/internal/scanners"
)

// Start Proxy
func StartProxyHandler(port int, jaeles bool, jaelesApi string, scope string, vulnArgs []string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, jaeles, jaelesApi, scope, vulnArgs)
	})

    server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.DefaultServeMux,
	}
	
	color.Cyan("Listening on port %d...\n", port)

    if jaeles {
        color.Cyan("Interacting with jaeles: %s\n", jaelesApi)
    }

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// Proxy Handler
func handler(w http.ResponseWriter, r *http.Request, jaeles bool, jaelesApi string, scope string, vulnArgs []string) {
	// Convert the raw request to base64
	requestBytes, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	requestBase64 := base64.StdEncoding.EncodeToString(requestBytes)

    // Send request to jaeles API server and filter if scope is specified
    if (scope != "" && regexp.MustCompile(scope).MatchString(r.URL.String()) || scope == "") {
        if jaeles {
            SendToJaeles(requestBase64, jaelesApi)
        }
        scanners.ScannerHandler(vulnArgs)
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