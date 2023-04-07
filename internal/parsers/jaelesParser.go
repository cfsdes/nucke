package parsers

import (
	"fmt"
	"net/http"
    "encoding/json"
    "bytes"
)

// Send request to jaeles API using the base64 target raw request
func SendToJaeles(base64str string, jaelesApi string) error {
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