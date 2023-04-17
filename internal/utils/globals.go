package utils

import (
	"os"
	"fmt"

	"github.com/cfsdes/nucke/internal/parsers"
)

// Flags
var Port int
var Threads int
var JaelesApi string
var Jaeles bool
var Scope string
var Proxy string
var Config string
var Output string
var FilePaths []string
var InteractURL string

// Initiate global variables
func init() {
	Port, Threads, JaelesApi, Jaeles, Scope, Proxy, Config, Output = ParseFlags()

	if Config != "" {
		FilePaths = parsers.ParseConfig(Config)
		InteractURL = StartInteractsh()
	}

	if Output != "" {
		Output = FormatOutput(Output)
		err := os.MkdirAll(Output, 0755)
		if err != nil {
			fmt.Println(err)
		}
	}
}
