package main

import (
	"github.com/cfsdes/nucke/internal/initializers"
	"github.com/cfsdes/nucke/internal/runner"
)

var version = "v0.2.3"

func main() {
	initializers.Start(version)

	// Start Proxy
	runner.StartProxy()
}
