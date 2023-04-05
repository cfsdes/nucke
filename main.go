package main

import (
	"fmt"
    "log"

    "github.com/fatih/color"
    "github.com/cfsdes/nucke/internal/runner"
    "github.com/cfsdes/nucke/internal/helpers"
)


func main() {
    // Init global variables
    helpers.InitGlobals()

    // List vulns
    if helpers.ListVulns {
		color.Blue("Available vulnerabilities:\n\n")
		for _, vuln := range helpers.VulnList {
			fmt.Println(vuln)
		}
        fmt.Println()
		return
	}

    // Validate vulns argument
    vulnArgs, err := helpers.ValidateVulns(helpers.Vulns)
    if err != nil {
        log.Fatal(err)
    }

    // Initial banner
    helpers.Banner()

	// Start Proxy
	runner.StartProxyHandler(vulnArgs)
}