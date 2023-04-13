package utils

import (
	"github.com/cfsdes/nucke/internal/parsers"
	pluginsUtils "github.com/cfsdes/nucke/pkg/plugins/utils"

)

// Flags
var Port int
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
	Port, JaelesApi, Jaeles, Scope, Proxy, Config, Output = ParseFlags()

	if Config != "" {
		FilePaths = parsers.ParseConfig(Config)
		InteractURL = pluginsUtils.StartInteractsh()
	}
}
