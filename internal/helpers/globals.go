package helpers

import (
)

// Flags
var Port int
var JaelesApi string
var Jaeles bool
var Scope string
var ListVulns bool
var Vulns string
var Proxy string

// Initiate global variables
func InitGlobals() {
	Port, JaelesApi, Jaeles, Scope, ListVulns, Vulns, Proxy = ParseFlags()
}
