package main

import (
	//"fmt"
    //"log"

    //"github.com/fatih/color"
    "github.com/cfsdes/nucke/internal/runner"
    "github.com/cfsdes/nucke/internal/utils"
)


func main() {
    // Initial banner
    utils.Banner()

    // Init global variables
    utils.InitGlobals()

	// Start Proxy
	runner.StartProxyHandler()
}