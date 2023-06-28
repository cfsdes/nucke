package requests

import (
    "net/http"
    "fmt"
    "io/ioutil"
    "strings"
    "bufio"
    "strconv"

    "github.com/cfsdes/nucke/internal/globals"
)

func RequestToRaw(r *http.Request) string {
    req := CloneReq(r)

    // Write the request line
    raw := fmt.Sprintf("%s %s %s\r\n", req.Method, req.URL.RequestURI(), req.Proto)

    // Write the Host header
    raw += fmt.Sprintf("Host: %s\r\n", req.Host)

    // Write the headers
    for name, values := range req.Header {
        for _, value := range values {
            raw += fmt.Sprintf("%s: %s\r\n", name, value)
        }
    }

    // Write a blank line to end the headers
    raw += "\r\n"

    // Write the body, if present
    if req.Body != nil {
        body, err := ioutil.ReadAll(req.Body)
        if err != nil {
            // If there was an error reading the body, return the raw request up to this point
            return raw
        }
        raw += string(body)
    }

    // Replace any occurrences of "\r\n" with "\n"
    raw = strings.ReplaceAll(raw, "\r\n", "\n")

    return raw
}


// Receives a raw request and parses it to return the final URL
func ExtractRawURL(rawRequest string) string {
    // Split the raw request into its parts
    scanner := bufio.NewScanner(strings.NewReader(rawRequest))
    scanner.Split(bufio.ScanLines)

    // Extract the request line
    scanner.Scan()
    requestLine := scanner.Text()

    // Extract the headers
    headers := make(map[string]string)
    for scanner.Scan() {
        line := scanner.Text()
        if line == "" {
            break
        }
        parts := strings.SplitN(line, ":", 2)
        if len(parts) == 2 {
            name := strings.TrimSpace(parts[0])
            value := strings.TrimSpace(parts[1])
            headers[name] = value
        }
    }

    // Extract the target URL
    parts := strings.SplitN(requestLine, " ", 3)
    if len(parts) < 2 {
        return ""
    }
    targetURL := parts[1]

    // Combine the target URL and Host header to form the return URL
    if host, ok := headers["Host"]; ok {
        return "http://" + host + targetURL
    } else {
        return ""
    }
}

// Return Status Code from Raw Response
func StatusCodeFromRaw(raw string) (int) {

    // Split the raw string into lines
    lines := strings.Split(raw, "\n")

    // Check if the first line contains a valid status line
    if len(lines) < 1 {
        if globals.Debug {
            fmt.Printf("Invalid Raw Response: empty string\n")
        }
        return 0
    }

    parts := strings.SplitN(lines[0], " ", 3)
    if len(parts) != 3 {
        if globals.Debug {
            fmt.Printf("Invalid Raw Response: invalid status line: %s\n", lines[0])
        }
        return 0
    }

    // Parse the status code and return it
    statusCode, err := strconv.Atoi(parts[1])
    if err != nil {
        if globals.Debug {
            fmt.Printf("Invalid Raw Response: invalid status code: %s\n", parts[1])
        }
        return 0
    }

    return statusCode
}
