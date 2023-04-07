package utils

import (
	"github.com/cfsdes/nucke/internal/parsers"
)

// Flags
var Port int
var JaelesApi string
var Jaeles bool
var Scope string
var Proxy string
var Config string
var FilePaths []string

// Initiate global variables
func InitGlobals() {
	Port, JaelesApi, Jaeles, Scope, Proxy, Config = ParseFlags()

	if Config != "" {
		FilePaths = parsers.ParseConfig(Config)
	}
}
