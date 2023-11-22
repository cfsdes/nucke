package globals

// Flags
var Port string 				// Port that nucke will listen
var Threads int 				// Nucke scan threads
var Delay int 					// Delay between fuzz requests in milliseconds
var JaelesApi string			// Jaeles server API url
var Jaeles bool					// Jaeles boolean flag
var Scope string				// Regex to set the scope to be scanned
var Proxy string				// Proxy to use during scan
var PluginsConfig string		// Config.yaml file for plugins
var Output string				// Output directory for plugins
var PluginPaths []string		// Plugins paths with plugins in .so format
var InteractURL string			// Interact URL for OOB scan
var ExportCA bool				// Export PEM certificate
var Debug bool 					// Debug Error messages
var Version bool				// Return Nucke Version
var Stats bool					// Start Status server
var ListPlugins bool			// List plugins available for the path in config.yaml
var PendingScans int64 			// Number of Pending requests
var Headers []string 			// Custom Headers
var CustomParams []string 		// Custom parameters to be used during scan
var Webhook map[string][]string	// Webhooks to send scan results


// Initiate global variables
func init() {
	Port, Threads, Delay, JaelesApi, Jaeles, Proxy, PluginsConfig, Output, ExportCA, Debug, Version, Stats, ListPlugins, Headers, CustomParams = ParseFlags()
}
