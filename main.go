package main

import (

    "github.com/cfsdes/nucke/internal/runner"
    "github.com/cfsdes/nucke/internal/initializers"
)

var version = "v0.0.5"

func main() {
    initializers.Start(version)

	// Start Proxy
	runner.StartProxy()
}