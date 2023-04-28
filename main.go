package main

import (
    "github.com/cfsdes/nucke/internal/runner"
    "github.com/cfsdes/nucke/internal/initializers"
)


func main() {
    // Check binaries
    binaries := []string{"interactsh-client"}
    initializers.CheckBinaries(binaries)

    // Initial banner
    initializers.Banner()

	// Start Proxy
	runner.StartProxy()
}