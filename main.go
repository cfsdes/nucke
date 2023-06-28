package main

import (
    "fmt"

    "github.com/cfsdes/nucke/internal/runner"
    "github.com/cfsdes/nucke/internal/initializers"
)

var version = "v0.0.1"

func main() {
    // Check binaries
    binaries := []string{"interactsh-client"}
    initializers.CheckBinaries(binaries)

    // Print Nucke version
    if initializers.Version {
		fmt.Println("\nNucke version: ",version, "\n")
        return
	}

	// Start Proxy
	runner.StartProxy()
}