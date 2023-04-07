package main

import (
    "net/http"
    "github.com/cfsdes/nucke/pkg/template"
)


func Run(r *http.Request, client *http.Client) (string, string, string, bool, error) {
	// Code here

    var severity string = "High"
    var description string = "vulnerability description"
    var url string = "http://example.com"
    var vulnFound bool = true
    
    reportContent := template.ReadFileToString("report-template.txt")
    summary := template.ParseTemplate(reportContent, map[string]interface{}{
        "description": description,
    })

    return	severity, url, summary, vulnFound, nil
}