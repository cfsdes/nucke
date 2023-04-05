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
    "github.com/cfsdes/nucke/internal/helpers"
)

// Start Proxy
func StartProxyHandler(vulnArgs []string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, vulnArgs)
	})

    server := &http.Server{
		Addr:    fmt.Sprintf(":%d", helpers.Port),
		Handler: http.DefaultServeMux,
	}
	
	color.Cyan("Listening on port %d...\n", helpers.Port)

    if helpers.Jaeles {
        color.Cyan("Interacting with jaeles: %s\n", helpers.JaelesApi)
    }

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// Proxy Handler
func handler(w http.ResponseWriter, r *http.Request, vulnArgs []string) {
	// Convert the raw request to base64
	requestBytes, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	requestBase64 := base64.StdEncoding.EncodeToString(requestBytes)

    // Send request to jaeles API server and filter if scope is specified
    if (helpers.Scope != "" && regexp.MustCompile(helpers.Scope).MatchString(r.URL.String()) || helpers.Scope == "") {
        if helpers.Jaeles {
            SendToJaeles(requestBase64, helpers.JaelesApi)
        }
        scanners.ScannerHandler(r, vulnArgs)
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