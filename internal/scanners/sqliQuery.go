package scanners

import (
    "fmt"
    "net/http"
	"io/ioutil"
)

func SqliQuery(r *http.Request, client *http.Client) (int64, error) {
    // Set empty RequestURI to avoid http: Request.RequestURI error
    r.RequestURI = ""

    // Send request
    resp, err := client.Do(r)
    if err != nil {
        return 0, fmt.Errorf("failed to send request: %s", err)
    }
    defer resp.Body.Close()

    // Read response body
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return 0, fmt.Errorf("failed to read response body: %s", err)
    }

    return int64(len(body)), nil
}


