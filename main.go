package main

import (
    "github.com/cfsdes/nucke/internal/runner"
    "github.com/cfsdes/nucke/internal/utils"
)


func main() {
    // Initial banner
    utils.Banner()

	// Start Proxy
	runner.StartProxyHandler()
}