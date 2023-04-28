package utils

import (
	"net/http"
    "io/ioutil"
    "time"
    "io"
    "compress/gzip"
    "fmt"
)

/**
* Just reproduce the request provided
*/

func BasicRequest(r *http.Request, client *http.Client) (int, string, int, map[string][]string) {
    req := CloneRequest(r)

    if client == nil {
        client = &http.Client{}
    }

    // Send the request
    start := time.Now()
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
        return 0, "", 0, nil
    }
    defer resp.Body.Close()

    // Get response body
	defer resp.Body.Close()
	
	var bodyReader io.ReadCloser
	
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

	bodyBytes, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		fmt.Println("Response parser error: ", err)
	}

    // Get response time
    elapsed := int(time.Since(start).Seconds())

    // Get response headers
    responseHeaders := make(map[string][]string)
    for k, v := range resp.Header {
        responseHeaders[k] = v
    }

    return elapsed, string(bodyBytes), resp.StatusCode, responseHeaders
}
