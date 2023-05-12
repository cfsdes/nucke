package main

import (
    "net/http"
    "fmt"
    "strings"

    "github.com/cfsdes/nucke/pkg/report"
    "github.com/cfsdes/nucke/pkg/requests"
    "github.com/cfsdes/nucke/internal/initializers"
)


func Run(r *http.Request, client *http.Client, pluginDir string) (string, string, string, bool, error) {
    // Scan
    vulnFound, rawReq, url := scan(r, client, pluginDir)

    // Report
    reportContent := report.ReadFileToString("report-template.md", pluginDir)
    summary := report.ParseTemplate(reportContent, map[string]interface{}{
        "request": rawReq,
    })
    
    return	"Medium", url, summary, vulnFound, nil
}


func scan(r *http.Request, client *http.Client, pluginDir string) (bool, string, string) {
    
    // Check if request requires authentication
    requireAuth := requests.CheckAuth(r, client)

    // Just test IDOR in authenticated endpoints
    if (requireAuth) {
        
        // Send basic request
        _, resBody, _, _, _ := requests.BasicRequest(r, client)
        originalLength := len(resBody)

        // Replace headers with user custom parameters
        if len(initializers.CustomParams) > 0 {
            for _, param := range initializers.CustomParams {
                parts := strings.SplitN(param, "=", -1)
                if len(parts) >= 2 {
                    key := fmt.Sprintf("{{.%s}}", strings.TrimSpace(parts[0]))
                    value := strings.TrimSpace(strings.Join(parts[1:], "="))
    
                    if key == "{{.idor_cookie}}" {
                        r.Header.Set("Cookie", value)
                    }

                    if key == "{{.idor_authorization}}" {
                        r.Header.Set("Authorization", value)
                    }
                }
            }
        }

        // Send request with Cookies of Account B
        _, resBody, _, _, _ = requests.BasicRequest(r, client)
        anotherLength := len(resBody)

        if (originalLength == anotherLength){
            rawReq := requests.RequestToRaw(r)
            url := requests.ExtractRawURL(rawReq)

            return true, rawReq, url
        }
    }
    
    return false, "", ""
}