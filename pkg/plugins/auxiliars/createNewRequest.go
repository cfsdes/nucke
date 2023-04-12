package auxiliars

import (
	"net/http"
    "bytes"
    "io/ioutil"
    "net/http/httputil"
    "log"
)

/*
* This function is used to create a new request based on the original request
* It's necessary because we can't read the body of the same request twice
*/

// Create a new request to forward
func CreateNewRequest(req *http.Request, w http.ResponseWriter) *http.Request {
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
            newReq.Header.Add(key, value)
        }
    }

    // Add the body of the original request to the new one
    requestBytes, err := httputil.DumpRequest(req, true)
    if err != nil {
        log.Fatal(err)
    }
    bodyStart := bytes.Index(requestBytes, []byte("\r\n\r\n")) + 4
    newReq.Body = ioutil.NopCloser(bytes.NewReader(requestBytes[bodyStart:]))
    
    return newReq
}