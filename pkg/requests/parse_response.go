package requests

import (
	"io/ioutil"
	"io"
	"fmt"
	"compress/gzip"
	"compress/flate"
	"net/http"
	"strings"

	"github.com/cfsdes/nucke/internal/initializers"
)

// Parse response
func ParseResponse(resp *http.Response) (int, string, map[string][]string, string) {
	
	// Check if resp is nil
	if resp == nil {
        return 0, "", nil, ""
    }

	// Get status code
    statusCode := resp.StatusCode

	// Get response body
	responseBody := getBody(resp)

	// Get response headers
    responseHeaders := make(map[string][]string)
    for k, v := range resp.Header {
        responseHeaders[k] = v
    }

	// Get Raw Response
	rawResp := createRawResponse(resp, responseBody)

	return statusCode, responseBody, responseHeaders, rawResp
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
	case "deflate":
		bodyReader = flate.NewReader(resp.Body)
	default:
		bodyReader = resp.Body
	}

	// Read body
	bodyBytes, err := ioutil.ReadAll(bodyReader)
	if err != nil && initializers.Debug {
		fmt.Println("Response parser error: ", err)
	}

	return string(bodyBytes)
}

func createRawResponse(resp *http.Response, body string) string {
	
	// Write the response status line
    raw := fmt.Sprintf("%s %d %s\r\n", resp.Proto, resp.StatusCode, http.StatusText(resp.StatusCode))

    // Write the headers
    for name, values := range resp.Header {
        for _, value := range values {
            raw += fmt.Sprintf("%s: %s\r\n", name, value)
        }
    }

    // Write a blank line to end the headers
    raw += "\r\n"

    // Write the body, if present
    raw += body

    // Replace any occurrences of "\r\n" with "\n"
    raw = strings.ReplaceAll(raw, "\r\n", "\n")

    return raw
}