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
    "github.com/cfsdes/nucke/internal/utils"
    "github.com/cfsdes/nucke/internal/parsers"
    pluginsUtils "github.com/cfsdes/nucke/pkg/plugins/utils"
)

// Create a channel with a buffer of threads
var ch = make(chan int, utils.Threads)

// Start Proxy
func StartProxyHandler() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	})

    server := &http.Server{
		Addr:    fmt.Sprintf(":%d", utils.Port),
		Handler: http.DefaultServeMux,
	}
	
	color.Cyan("Listening on port %d...\n", utils.Port)

    if utils.Jaeles {
        color.Cyan("Interacting with jaeles: %s\n", utils.JaelesApi)
    }

    fmt.Println()

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
    
    // Set the Host field of the request URL based on the Host header of the incoming request
    r.URL.Host = r.Host

    // Send request to jaeles API server and filter if scope is specified
    if (utils.Scope != "" && regexp.MustCompile(utils.Scope).MatchString(r.URL.String()) || utils.Scope == "") {
        
        // If jaeles scan is enabled
        if utils.Jaeles {
            parsers.SendToJaeles(requestBase64, utils.JaelesApi)
        }

        // If config with plugins is provided
        if utils.Config != "" {
            // Clone request before scanning
            req := pluginsUtils.CloneRequest(r)
            
            // executa a ScannerHandler dentro de uma goroutine
            go func() {
                // adiciona 1 ao canal para indicar que está utilizando uma goroutine
                ch <- 1

                ScannerHandler(req, w)

                // sinaliza ao canal que a goroutine está livre
                <-ch
            }()
        }
	} 

    fowardRequest(w, r)
}

func fowardRequest(w http.ResponseWriter, r *http.Request) {
    // Add default protocol scheme if missing
    if r.URL.Scheme == "" {
        if r.TLS != nil {
            r.URL.Scheme = "https"
        } else {
            r.URL.Scheme = "http"
        }
    }

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