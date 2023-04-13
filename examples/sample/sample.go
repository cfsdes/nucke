package main

import (
    "net/http"
    "github.com/cfsdes/nucke/pkg/template"
)


func Run(r *http.Request, w http.ResponseWriter, client *http.Client, pluginDir string) (string, string, string, bool, error) {
	// Scan
    vulnFound := scan(r, w, client)
    
    // Report details
    var severity string = "High"
    var description string = "vulnerability description"
    var url string = "http://example.com"
    
    reportContent := template.ReadFileToString("report-template.txt", pluginDir)
    summary := template.ParseTemplate(reportContent, map[string]interface{}{
        "description": description,
    })

    return	severity, url, summary, vulnFound, nil
}


func scan(r *http.Request, w http.ResponseWriter, client *http.Client) bool {}