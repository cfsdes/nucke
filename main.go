package main

import (
	"fmt"
    "log"

    "github.com/fatih/color"
    "github.com/cfsdes/nucke/internal/runner"
    "github.com/cfsdes/nucke/internal/helpers"
)


func main() {
    port, jaelesApi, jc, scope, listVulns, vulns := helpers.ParseFlags()

    // List vulns
    if listVulns {
		color.Blue("Available vulnerabilities:\n\n")
		for _, vuln := range helpers.VulnList {
			fmt.Println(vuln)
		}
        fmt.Println()
		return
	}

    // Validate vulns argument
    vulnArgs, err := helpers.ValidateVulns(vulns)
    if err != nil {
        log.Fatal(err)
    }

    // Initial banner
    helpers.Banner()

	// Start Proxy
	runner.StartProxyHandler(port, jc, jaelesApi, scope, vulnArgs)
}