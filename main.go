package main

import (
	"fmt"
    "log"

    "github.com/fatih/color"
    "github.com/cfsdes/nucke/internal/runner"
    "github.com/cfsdes/nucke/internal/utils"
)


func main() {
    // Init global variables
    utils.InitGlobals()

    // List vulns
    if utils.ListVulns {
		color.Blue("Available vulnerabilities:\n\n")
		for _, vuln := range utils.VulnList {
			fmt.Println(vuln)
		}
        fmt.Println()
		return
	}

    // Validate vulns argument
    vulnArgs, err := utils.ValidateVulns(utils.Vulns)
    if err != nil {
        log.Fatal(err)
    }

    // Initial banner
    utils.Banner()

	// Start Proxy
	runner.StartProxyHandler(vulnArgs)
}