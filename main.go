package main

import (
    "github.com/cfsdes/nucke/internal/runner"
    "github.com/cfsdes/nucke/internal/utils"
)


func main() {
    // Check binaries
    binaries := []string{"interactsh-client"}
    utils.CheckBinaries(binaries)

    // Initial banner
    utils.Banner()

	// Start Proxy
	runner.StartProxyHandler()
}