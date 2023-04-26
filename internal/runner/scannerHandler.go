package runner

import (
	"net/http"
	"net/url"
	"fmt"
	"os"
    "plugin"
	"path/filepath"
    "os/user"
    "strings"

	"github.com/cfsdes/nucke/internal/utils"
)

func ScannerHandler(req *http.Request, w http.ResponseWriter) {
	// Create HTTP Client
	client, err := createHTTPClient()
	if err != nil {
		fmt.Println(err)
	}

	// Run Config Plugins
	for _, plugin := range utils.PluginPaths {
		runPlugin(plugin, req, w, client)
	}
}

// Generate HTTP Client with Proxy
func createHTTPClient() (*http.Client, error) {
    var client *http.Client
    if utils.Proxy != "" {
        // Create HTTP client with proxy
        proxyUrl, err := url.Parse(utils.Proxy)
        if err != nil {
            return nil, fmt.Errorf("failed to parse proxy URL: %s", err)
        }
        client = &http.Client{
            Transport: &http.Transport{
                DisableKeepAlives: true,
                Proxy: http.ProxyURL(proxyUrl),
            },
        }
    } else {
        // Create HTTP client without proxy
        client = &http.Client{}
    }
    return client, nil
}


// Run plugin
func runPlugin(scannerPlugin string, req *http.Request, w http.ResponseWriter, client *http.Client) {
	// Get the current user's home directory
    usr, err := user.Current()
    if err != nil {
        fmt.Println("Error getting current user:", err)
        os.Exit(1)
    }

    // Replace "~" with the home directory in the plugin path
    pluginPath := strings.Replace(scannerPlugin, "~", usr.HomeDir, 1) // path/to/plugin/plugin.so
    pluginDir := filepath.Dir(pluginPath)   // path/to/plugin

    // Load the plugin file
    plug, err := plugin.Open(pluginPath)
    if err != nil {
        fmt.Println("Error loading plugin:", err)
        os.Exit(1)
    }

    // Look up the run() function
    runFunc, err := plug.Lookup("Run")
    if err != nil {
        fmt.Println("Error looking up run() function:", err)
        os.Exit(1)
    }

    // Get plugin name (without extension and "." in the start of the name)
    scanName := filepath.Base(pluginDir)
    scanName = strings.TrimSuffix(scanName, filepath.Ext(scanName))
    scanName = strings.TrimPrefix(scanName, ".") // e.g: sqli

    // Call the run() function with a req argument
	severity, url, summary, found, err := runFunc.(func(*http.Request, http.ResponseWriter, *http.Client, string) (string, string, string, bool, error))(req, w, client, pluginDir)
    if err != nil {
        fmt.Println("Error running plugin:", err)
        os.Exit(1)
    }

	// Parse output if vulnerability is found
	if found {
		
		
		utils.VulnerabilityOutput(scanName, severity, url, summary)
	}
}