package main

import (
    "net/http"
    "github.com/cfsdes/nucke/pkg/report"
    "github.com/cfsdes/nucke/pkg/plugins/utils"
    "github.com/cfsdes/nucke/pkg/plugins/fuzzers"
    "github.com/cfsdes/nucke/pkg/requests"
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
    
    return	"High", url, summary, vulnFound, rawResp, nil
}


// Running all Fuzzers
func scan(r *http.Request, client *http.Client, pluginDir string) (bool, string, string, string) {
    
    // Make basic request
    originalResTime, _, _, _, _ := requests.BasicRequest(r, client)

    if originalResTime < 20 {
        payloads := utils.FileToSlice(pluginDir, "payloads.txt")
        matcher := detections.Matcher{
            Time: &detections.TimeMatcher{
                Operator: ">=",
                Seconds:  20,
            },
        }

        headers := []string{"User-Agent","X-Forwarded-For"}
        match, rawReq, url, _, _, rawResp := fuzzers.FuzzHeaders(r, client, payloads, headers, matcher)

        if match {
            return match, rawReq, url, rawResp
        }
    }

    return false, "", "", ""
}