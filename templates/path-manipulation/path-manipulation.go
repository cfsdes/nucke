package main

import (
    "net/http"
    "github.com/cfsdes/nucke/pkg/report"
    "github.com/cfsdes/nucke/pkg/plugins/fuzzers"
    "github.com/cfsdes/nucke/pkg/plugins/detections"
    "github.com/cfsdes/nucke/pkg/requests"
)


func Run(r *http.Request, client *http.Client, pluginDir string) (string, string, string, bool, error) {
    // Scan
    vulnFound, rawReq, url := scan(r, client, pluginDir)

    // Report
    reportContent := report.ReadFileToString("report-template.md", pluginDir)
    summary := report.ParseTemplate(reportContent, map[string]interface{}{
        "request": rawReq,
    })
    
    return	"Info", url, summary, vulnFound, nil
}


// Running all Fuzzers
func scan(r *http.Request, client *http.Client, pluginDir string) (bool, string, string) {

    // Make basic request
    _, resBody, _, _ := requests.BasicRequest(r, client)
    originalLength := len(resBody)

    // Compare the original length with the length with payload
    payload1 := []string{".../{{.original}}"}
    matcher := detections.Matcher{
        ContentLength: &detections.ContentLengthMatcher{
            Operator: "!=",
            Length: originalLength,
        },
    }

    fuzzers := []func(*http.Request, *http.Client, []string, detections.Matcher) (bool, string, string, string, string){
        fuzzers.FuzzJSON,
        fuzzers.FuzzQuery,
        fuzzers.FuzzFormData,
        fuzzers.FuzzXML,
    }

    for _, fuzzer := range fuzzers {
        if match1, _, _, _, param1 := fuzzer(r, client, payload1, matcher); match1 {
            /*
                If length with payload is different of original length,
                Try to "fix" the query. If the payload with query fixing in the same parameter
                return a response equal to the original, it's vulnerable.   
            */
            
            payload2 := []string{"./myw/../{{.original}}"}
            matcher2 := detections.Matcher{
                ContentLength: &detections.ContentLengthMatcher{
                    Operator: "==",
                    Length: originalLength,
                },
            }

            for _, fuzzer := range fuzzers {
                if vulnFound, rawReq, url, _, param2 := fuzzer(r, client, payload2, matcher2); vulnFound {
                    if param2 == param1 {
                        return vulnFound, rawReq, url
                    }
                }
            }
        }
    }
    
    return false, "", ""
    
}
