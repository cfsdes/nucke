package main

import (
    "net/http"
    "github.com/cfsdes/nucke/pkg/report"
    "github.com/cfsdes/nucke/pkg/plugins/utils"
    "github.com/cfsdes/nucke/pkg/plugins/fuzzers"
    "github.com/cfsdes/nucke/pkg/plugins/detections"
)


func Run(r *http.Request, client *http.Client, pluginDir string) (string, string, string, bool, string, error) {
    // Scan
    vulnFound, _, _, _ := scan(r, client, pluginDir)
    
    // Run twice to avoid false positives
    if vulnFound {
        vulnFound, rawReq, url, rawResp := scan(r, client, pluginDir)

        // Report
        reportContent := report.ReadFileToString("report-template.md", pluginDir)
        summary := report.ParseTemplate(reportContent, map[string]interface{}{
            "request": rawReq,
        })
        
        return	"High", url, summary, vulnFound, rawResp, nil
    }

    return "", "", "", false, "", nil
}


// Running all Fuzzers
func scan(r *http.Request, client *http.Client, pluginDir string) (bool, string, string, string) {
    
    // Call All Fuzzers (Except fuzzHeaders)
    payloads := utils.FileToSlice(pluginDir, "payloads.txt")
    matcher := detections.Matcher{
        Time: &detections.TimeMatcher{
            Operator: ">=",
            Seconds:  20,
        },
    }

    fuzzers := []func(*http.Request, *http.Client, []string, detections.Matcher) (bool, string, string, string, string, string){
        fuzzers.FuzzJSON,
        fuzzers.FuzzQuery,
        fuzzers.FuzzFormData,
        fuzzers.FuzzXML,
    }

    for _, fuzzer := range fuzzers {
        if vulnFound, rawReq, url, _, _, rawResp := fuzzer(r, client, payloads, matcher); vulnFound {
            return vulnFound, rawReq, url, rawResp
        }
    }

    return false, "", "", ""
}