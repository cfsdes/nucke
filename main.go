package main

import (

    "github.com/cfsdes/nucke/internal/runner"
    "github.com/cfsdes/nucke/internal/initializers"
)

var version = "v1.0.2"

func main() {
    initializers.Start(version)

	// Start Proxy
	runner.StartProxy()
}