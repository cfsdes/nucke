package main

import (
    "net/http"
    "fmt"
    "strings"

    "github.com/cfsdes/nucke/pkg/report"
    "github.com/cfsdes/nucke/pkg/requests"
)


func Run(r *http.Request, client *http.Client, pluginDir string) (string, string, string, bool, string, error) {
    // Scan
    vulnFound, rawReq, url, rawResp := scan(r, client)

    // Report
    reportContent := report.ReadFileToString("report-template.md", pluginDir)
    summary := report.ParseTemplate(reportContent, map[string]interface{}{
        "request": rawReq,
    })
    
    return	"Medium", url, summary, vulnFound, rawResp, nil
}


func scan(r *http.Request, client *http.Client) (bool, string, string, string) {

    // Check if request requires authentication
    requireAuth := requests.CheckAuth(r, client)

    if (requireAuth) {
        // Test CORS using Arbitrary Origin
        found, rawReq, url, rawResp := arbitraryOriginCheck(r, client)
        if found { return found, rawReq, url, rawResp }
        
        // Test CORS using Null Origin
        found, rawReq, url, rawResp = nullOriginCheck(r, client)
        if found { return found, rawReq, url, rawResp }
    }
    

    return false, "", "", ""
}

// Test CORS using Arbitrary Origin
func arbitraryOriginCheck(r *http.Request, client *http.Client) (bool, string, string, string) {
    
    // Format CORS URL
    hostParts := strings.Split(r.Host, ":")
    domain := hostParts[0]
    corsURL := fmt.Sprint("https://", domain, ".example.com")

    // Add CORS headers
    r.Header.Set("Origin", corsURL)

    // Send Request
    _, _, statusCode, headers, rawResp := requests.BasicRequest(r, client)

    // Get raw req and URL
    rawReq := requests.RequestToRaw(r)
    url := requests.ExtractRawURL(rawReq)

    // Verifica se o mapa de headers contém os headers necessários
    if statusCode < 300 &&
    containsHeader(headers, "Access-Control-Allow-Origin", corsURL) &&
    containsHeader(headers, "Access-Control-Allow-Credentials", "true") {
        return true, rawReq, url, rawResp
    }

    if statusCode < 300 &&
    containsHeader(headers, "Access-Control-Allow-Origin", "*") &&
    containsHeader(headers, "Access-Control-Allow-Credentials", "true") {
        return true, rawReq, url, rawResp
    }

    return false, "", "", ""
}

// Test CORS using Null Origin
func nullOriginCheck(r *http.Request, client *http.Client) (bool, string, string, string) {
    
    // Add CORS headers
    if r.Header.Get("Origin") != "" {
        r.Header.Set("Origin", "null")
    } else {
        r.Header.Add("Origin", "null")
    }

    // Send Request
    _, _, statusCode, headers, rawResp := requests.BasicRequest(r, client)

    // Get raw req and URL
    rawReq := requests.RequestToRaw(r)
    url := requests.ExtractRawURL(rawReq)

    // Verifica se o mapa de headers contém os headers necessários
    if statusCode < 300 &&
    containsHeader(headers, "Access-Control-Allow-Origin", "null") &&
    containsHeader(headers, "Access-Control-Allow-Credentials", "true") {
        return true, rawReq, url, rawResp
    }

    return false, "", "", ""
}


// Check if headers variable contains header value
func containsHeader(headers map[string][]string, headerName string, expectedValue string) bool {
	// Verifica se o mapa de headers contém o header especificado
	headerValues, ok := headers[headerName]
	if !ok {
		return false
	}

	// Verifica se o header tem o valor esperado
	return strings.Join(headerValues, ",") == expectedValue
}
