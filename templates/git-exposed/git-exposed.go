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
    payloads := []string{".git/config?nil="}
    
    // Detection
    matcher := detections.Matcher{
        Body: &detections.BodyMatcher{
            RegexList: []string{"[core]"},
        },
        StatusCode: &detections.StatusCodeMatcher{
            Operator: "==",
            Code: 200,
        },
        Operator: "AND",
    }

    if vulnFound, rawReq, url, _, _, rawResp := fuzzers.FuzzPath(r, client, payloads, matcher, "*"); vulnFound {
        return vulnFound, rawReq, url, rawResp
    }

    return false, "", "", ""
}