package requests

import (
	"net/http"
    "bytes"
    "io/ioutil"
    "net/http/httputil"
    "log"
    "strings"

    "github.com/cfsdes/nucke/pkg/globals"
)

/*
* This function is used to create a new request based on the original request
* It's useful because we can't read the body of the same request twice
*/

// Create a new request to forward
func CloneReq(req *http.Request) *http.Request {
	// Create a new request based on the original one
    // but with an empty body
    body := []byte{}
    newReq, err := http.NewRequest(req.Method, req.URL.String(), bytes.NewBuffer(body))
    if err != nil {
        log.Fatal(err)
    }
    
    // Copy the headers from the original request to the new one
    for key, values := range req.Header {
        for _, value := range values {
            newReq.Header.Set(key, value)
        }
    }

    // Add custom headers from --headers parameter
    if len(globals.Headers) > 0 {
        for _, header := range globals.Headers {
            parts := strings.SplitN(header, ":", -1)
            if len(parts) >= 2 {
                key := strings.TrimSpace(parts[0])
                value := strings.TrimSpace(strings.Join(parts[1:], ":"))
                newReq.Header.Set(key, value)
            }
        }
    }

    // Add Accept-Encoding: gzip, deflate
    newReq.Header.Set("Accept-Encoding", "gzip, deflate")

    // Delete If-None-Match
    newReq.Header.Del("If-None-Match")
    newReq.Header.Del("If-Modified-Since")

    // Add the body of the original request to the new one
    requestBytes, err := httputil.DumpRequest(req, true)
    if err != nil {
        log.Fatal(err)
    }
    bodyStart := bytes.Index(requestBytes, []byte("\r\n\r\n")) + 4
    newReq.Body = ioutil.NopCloser(bytes.NewReader(requestBytes[bodyStart:]))
    
    return newReq
}