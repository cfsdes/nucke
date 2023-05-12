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
    
    // Read rules file
    rules := utils.FileToSlice(pluginDir, "regex_match.txt")

    // Creating payload and matcher
    payloads := []string{"{{.original}}'", "{{.original}}\\"}
    matcher := detections.Matcher{
        Body: &detections.BodyMatcher{
            RegexList: rules,
        },
    }

    // Running All Fuzzers (Except fuzzHeaders)
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