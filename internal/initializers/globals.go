package initializers

import (
	"os"
	"fmt"

	"github.com/cfsdes/nucke/internal/initializers/plugins"
)

// Flags
var Port string 			// Port that nucke will listen
var Threads int 			// Nucke scan threads
var JaelesApi string		// Jaeles server API url
var Jaeles bool				// Jaeles boolean flag
var Scope string			// Regex to set the scope to be scanned
var Proxy string			// Proxy to use during scan
var Config string			// Config.yaml file for plugins
var Output string			// Output directory for plugins
var FilePaths []string		// File paths with plugins in golang format
var PluginPaths []string	// Plugins paths with plugins in .so format
var InteractURL string		// Interact URL for OOB scan
var UpdatePlugins bool		// Force the update of all plugins
var ExportCA bool			// Export PEM certificate
var Debug bool 				// Debug Error messages
var Verbose bool			// Verbose Output
var Stats bool				// Start Status server
var PendingScans int64 		// Number of Pending requests
var Headers []string 		// Custom Headers
var CustomParams []string 	// Custom parameters to be used during scan


// Initiate global variables
func init() {
	Port, Threads, JaelesApi, Jaeles, Scope, Proxy, Config, Output, UpdatePlugins, ExportCA, Debug, Verbose, Stats, Headers, CustomParams = ParseFlags()

	// Initial banner
	Banner()

	if Config != "" {
		// Parse Config.yaml
		FilePaths = plugins.ParseConfig(Config)

		// Build plugins
		PluginPaths = plugins.BuildPlugins(FilePaths, UpdatePlugins)

		// Start interact.sh
		StartInteractsh()
	}

	if Output != "" {
		Output = FormatOutput(Output)
		err := os.MkdirAll(Output, 0755)
		if err != nil {
			fmt.Println("Error creating output path:",err)
		}
	}
}
