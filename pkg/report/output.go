package report

import (
	"fmt"
	"crypto/sha1"
	"encoding/hex"
	"net/url"
	"io/ioutil"
	"path/filepath"
	"strings"
	"os/user"
	"os"

	"github.com/fatih/color"
	"github.com/cfsdes/nucke/pkg/globals"
)

var issueReported = make(map[string]bool)

func Output(scanName, webhook, severity, url, payload, param, rawReq, rawResp, pluginDir string) {
	// Check for duplicate
	hash := createSha1Hash(scanName, url, param)
	if (!issueReported[hash]) {
		issueReported[hash] = true

		yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
		red := color.New(color.FgRed, color.Bold).SprintFunc()
		green := color.New(color.FgGreen, color.Bold).SprintFunc()
		blue := color.New(color.FgBlue, color.Bold).SprintFunc()
		magenta := color.New(color.FgMagenta, color.Bold).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()

		switch severity {
		case "Critical":
			fmt.Printf("[%s] [%s] %s \n", cyan(scanName), magenta(severity), url)
		case "High":
			fmt.Printf("[%s] [%s] %s \n", cyan(scanName), red(severity), url)
		case "Medium":
			fmt.Printf("[%s] [%s] %s \n", cyan(scanName), yellow(severity), url)
		case "Low":
			fmt.Printf("[%s] [%s] %s \n", cyan(scanName), green(severity), url)
		case "Info":
			fmt.Printf("[%s] [%s] %s \n", cyan(scanName), blue(severity), url)
		case "OOB":
			fmt.Printf("[%s] [%s] %s \n", cyan(scanName), magenta(severity), url)
		}

		// Create summary
		summary := createSummary(pluginDir, scanName, severity, rawReq, rawResp, url, payload, param)

		// Create output
		if globals.Output != "" {
			outputPath := getOutputPath(scanName, url, param)
			
			err := writeStringToFile(summary, outputPath)
			if err != nil {
				panic(err)
			}
		}

		// Send webhook request
		Notify(scanName, severity, url, summary, webhook)
	}
}

// Format output
func FormatOutput(output string) string {
	// Get the current user's home directory
    usr, err := user.Current()
    if err != nil {
        fmt.Println("Error getting current user:", err)
        os.Exit(1)
    }

    // Replace "~" with the home directory in the plugin path
    outputPath := strings.Replace(output, "~", usr.HomeDir, 1) // path/to/plugin/plugin.so
	
	return outputPath
}

// return the domain of the url
func getDomain(urlString string) string {
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		panic(err)
	}
	return parsedURL.Hostname()
}

// generate random SHA1 hash
func createSha1Hash(scanName, url, param string) string {
	// generate random SHA1 hash
	signatureString := fmt.Sprintf("%s-%s-%s", scanName, url, param)
	hash := sha1.Sum([]byte(signatureString))
	hashString := hex.EncodeToString(hash[:])

	return hashString
}

// Function to return the outputPath
func getOutputPath(scanName, url, param string) string {
	// parse domain
	domain := getDomain(url)

	// generate random SHA1 hash
	hash := createSha1Hash(scanName, url, param)

	// filename
	fileName := scanName + "-" + hash + ".md"

	// set output path
	outputPath := filepath.Join(globals.Output, domain, fileName)
	return outputPath
}

// Function to write string to file
func writeStringToFile(text string, filePath string) error {
	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		return err
	}
	
	err = ioutil.WriteFile(filePath, []byte(text), 0644)
	if err != nil {
		return err
	}
	return nil
}

func createSummary(pluginDir, scanName, severity, rawReq, rawResp, url, payload, param string) string {
    _, err := os.Stat(pluginDir + "/report-template.md")

	var reportContent string

    // Se o report-template.md existir na pasta do plugin
    if err == nil {
        reportContent = ReadFileToString("report-template.md", pluginDir)
    } else {
		reportContent = `
## Issue Details
					
- **Scan Name**: {{.scanName}}
- **Severity**: {{.severity}}
- **Vulnerable Endpoint**: {{.url}}
- **Injection Point**: {{.param}}

## Proof of Concept

The payload below was used on **{{.param}}** to trigger the vulnerability:
` + "```" + `
{{.payload}}
` + "```" + `

### HTTP Request
Request:
` + "```http" + `
{{.request}}
` +  "```" + `
					
Response:
` + "```http" + `
{{.response}}
` + "```"
    }

	summary := ParseTemplate(reportContent, map[string]interface{}{
		"request": rawReq,
		"response": rawResp,
		"url": url,
		"payload": payload,
		"param": param,
		"scanName": scanName,
		"severity": severity,
	})

	return summary
}

func isDuplicate(scanName, url, param string) {
	hash := createSha1Hash(scanName, url, param)
	if (!issueReported[hash]) {

	}
	issueReported[hash] = true
}