package utils

import (
	"net/http"
    "io/ioutil"
)

func BasicRequest(r *http.Request, w http.ResponseWriter, client *http.Client) (string, error) {
    req := CloneRequest(r, w)

    if client == nil {
        client = &http.Client{}
    }
    
    // Send the request
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    // Read the response body
    responseBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    return string(responseBody), nil
}