package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
    "flag"
    "encoding/json"
    "bytes"
)


func parseFlags() (port int, jcAPI string) {
	flag.IntVar(&port, "port", 8080, "port number to use")
    flag.StringVar(&jcAPI, "jc-api", "http://127.0.0.1:5000", "jcAPI value to send to handler")
	flag.Parse()
	return
}

func main() {
    port, jaelesApi := parseFlags()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, jaelesApi)
	})

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
func handler(w http.ResponseWriter, r *http.Request, jaelesApi string) {
	// Convert the raw request to base64
	requestBytes, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	requestBase64 := base64.StdEncoding.EncodeToString(requestBytes)

	// Print the base64-encoded request
	fmt.Println(requestBase64)
    
    //Send request to jaeles API server
    sendToJaeles(requestBase64, jaelesApi)

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


// Send request to jaeles API using the base64 target raw request
func sendToJaeles(base64str string, jaelesApi string) error {
    // Create request body as JSON
    reqBody := map[string]string{
        "req": base64str,
    }
    reqBytes, err := json.Marshal(reqBody)
    if err != nil {
        return err
    }

    // Create request
    req, err := http.NewRequest("POST", jaelesApi + "/api/parse", bytes.NewReader(reqBytes))
    if err != nil {
        return err
    }

    // Add headers
    req.Header.Add("User-Agent", "Jaeles Scanner")
    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("Content-Length", string(len(reqBytes)))
    req.Close = true

    // Send request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
    }

    return nil
}