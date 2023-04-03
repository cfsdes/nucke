package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
    "flag"
)

func parseFlags() (port int, timeout int) {
	flag.IntVar(&port, "port", 8080, "port number to use")
	flag.IntVar(&timeout, "timeout", 60, "timeout in seconds")
	flag.Parse()
	return
}

func main() {
    port, _ := parseFlags()

	http.HandleFunc("/", handler)
    server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.DefaultServeMux,
	}

	log.Printf("Listening on port %d...\n", port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// Proxy Handler
func handler(w http.ResponseWriter, r *http.Request) {
	// Convert the raw request to base64
	requestBytes, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	requestBase64 := base64.StdEncoding.EncodeToString(requestBytes)

	// Print the base64-encoded request
	fmt.Println(requestBase64)

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