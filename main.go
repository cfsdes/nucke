package main

import (

    "github.com/cfsdes/nucke/internal/runner"
    "github.com/cfsdes/nucke/internal/initializers"
)

var version = "v2.0.1"

func main() {
    initializers.Start(version)

	// Start Proxy
	runner.StartProxy()
}