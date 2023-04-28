package requests

import (
	"io/ioutil"
	"io"
	"fmt"
	"compress/gzip"
	"net/http"
)

// Parse response
func ParseResponse(resp *http.Response) (int, string, map[string][]string) {
	// Get status code
    statusCode := resp.StatusCode

	// Get response body
	responseBody := getBody(resp)

	// Get response headers
    responseHeaders := make(map[string][]string)
    for k, v := range resp.Header {
        responseHeaders[k] = v
    }

	return statusCode, responseBody, responseHeaders
}

func getBody(resp *http.Response) string {
	
	// Get response body
	defer resp.Body.Close()
	
	var bodyReader io.ReadCloser
	var err error
	
	// Unpack/Deflate gzip response
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		bodyReader, err = gzip.NewReader(resp.Body)
		if err != nil {
			fmt.Println("Response parser gzip error: ", err)
		}
		defer bodyReader.Close()
	default:
		bodyReader = resp.Body
	}

	// Read body
	bodyBytes, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		fmt.Println("Response parser error: ", err)
	}

	return string(bodyBytes)
}