package main

import (
    "net/http"
    "github.com/cfsdes/nucke/pkg/report"
    "github.com/cfsdes/nucke/pkg/plugins/fuzzers"
    "github.com/cfsdes/nucke/pkg/plugins/detections"
)


func Run(r *http.Request, client *http.Client, pluginDir string) (string, string, string, bool, string, error) {
    // Scan
    vulnFound, rawReq, url, rawResp := scan(r, client, pluginDir)

    // Report
    reportContent := report.ReadFileToString("report-template.md", pluginDir)
    summary := report.ParseTemplate(reportContent, map[string]interface{}{
        "request": rawReq,
    })
    
    return	"Medium", url, summary, vulnFound, rawResp, nil
}


// Running all Fuzzers
func scan(r *http.Request, client *http.Client, pluginDir string) (bool, string, string, string) {
    
    // Payloads
    payloads := []string{
        "{{.original}}\x0D\x0ASet-Cookie:crlf=injection", 
        "{{.original}}\x0DSet-Cookie:crlf=injection", 
        "{{.original}}\x0ASet-Cookie:crlf=injection", 
        "{{.original}}%E5%98%8D%E5%98%8ASet-Cookie:crlf=injection",
    }
    
    // Detection
    matcher := detections.Matcher{
        Header: &detections.HeaderMatcher{
            RegexList: []string{"(?m)^Set-Cookie: crlf=injection"},
        },
    }

    if vulnFound, rawReq, url, _, _, rawResp := fuzzers.FuzzPath(r, client, payloads, matcher, "last"); vulnFound {
        return vulnFound, rawReq, url, rawResp
    }

    return false, "", "", ""
}