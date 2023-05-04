package main

import (
    "net/http"
    "github.com/cfsdes/nucke/pkg/report"
    "github.com/cfsdes/nucke/pkg/plugins/utils"
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
    
    return	"High", url, summary, vulnFound, nil
}


// Running all Fuzzers
func scan(r *http.Request, client *http.Client, pluginDir string) (bool, string, string) {

    // Scan query fixing content length >= & ,=
    vulnFound, rawReq, url := queryFixingContentLengthBased(r, client, pluginDir)
    if vulnFound {
        return vulnFound, rawReq, url
    }

    // Scan query fixing status code
    vulnFound, rawReq, url := queryFixingStatusCodeBased(r, client, pluginDir)
    if vulnFound {
        return vulnFound, rawReq, url
    }

    return false, "", ""
    
}

func queryFixingContentLengthBased(r *http.Request, client *http.Client, pluginDir string) (bool, string, string) {
    
    // Make basic request
    _, resBody, _, _ := requests.BasicRequest(r, client)
    originalLength := len(resBody)

    // Payloads
    payload1 := []string{"{{.original}}'"}
    
    // Compare the original length with the length with payload
    matcher := detections.Matcher{
        ContentLength: &detections.ContentLengthMatcher{
            Operator: ">=",
            Length: originalLength+200,
        },
    }
    matcher2 := detections.Matcher{
        ContentLength: &detections.ContentLengthMatcher{
            Operator: "<=",
            Length: originalLength-200,
        },
    }

    // Fuzzing
    fuzzers := []func(*http.Request, *http.Client, []string, detections.Matcher) (bool, string, string, string, string){
        fuzzers.FuzzJSON,
        fuzzers.FuzzQuery,
        fuzzers.FuzzFormData,
        fuzzers.FuzzXML,
    }

    // Compare the original length with the length with payload
    for _, fuzzer := range fuzzers {
        if match1, _, _, _, param1 := fuzzer(r, client, payload1, matcher); match1 {
            /*
                If length with payload is >= original length + 200,
                Try to "fix" the query. If the payload with query fixing in the same parameter
                return a response equal to the original, it's vulnerable.   
            */
            
            payload2 := utils.FileToSlice(pluginDir, "payloads.txt")
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

    for _, fuzzer := range fuzzers {
        if match1, _, _, _, param1 := fuzzer(r, client, payload1, matcher2); match1 {
            /*
                If length with payload is <= original length - 200,
                Try to "fix" the query. If the payload with query fixing in the same parameter
                return a response equal to the original, it's vulnerable.  
            */
            
            payload2 := utils.FileToSlice(pluginDir, "payloads.txt")
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

func queryFixingStatusCodeBased(r *http.Request, client *http.Client, pluginDir string) (bool, string, string) {
    
    // Make basic request
    _, _, originalStatusCode, _ := requests.BasicRequest(r, client)

    // Compare the original length with the length with payload
    payload1 := []string{"{{.original}}'"}
    matcher := detections.Matcher{
        StatusCode: &detections.StatusCodeMatcher{
            Operator: "!=",
            Code: originalStatusCode,
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
                If length with payload is different of original status code,
                Try to "fix" the query. If the payload with query fixing in the same parameter
                return a status equal to the original, it's vulnerable.   
            */
            
            payload2 := utils.FileToSlice(pluginDir, "payloads.txt")
            matcher2 := detections.Matcher{
                StatusCode: &detections.StatusCodeMatcher{
                    Operator: "==",
                    Code: originalStatusCode,
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