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
	"github.com/cfsdes/nucke/internal/parsers"
)

func ScannerHandler(req *http.Request) {
	// Create HTTP Client
	client, err := createHTTPClient()
	if err != nil {
		fmt.Println(err)
	}

	// Handle Config Plugins
	for _, plugin := range utils.FilePaths {
		runPlugin(plugin, req, client)
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
func runPlugin(scannerPlugin string, req *http.Request, client *http.Client) {
	// Get the current user's home directory
    usr, err := user.Current()
    if err != nil {
        fmt.Println("Error getting current user:", err)
        os.Exit(1)
    }

    // Replace "~" with the home directory in the plugin path
    pluginPath := strings.Replace(scannerPlugin, "~", usr.HomeDir, 1)

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

    // Call the run() function with a req argument
	severity, url, rawReq, desc, found, err := runFunc.(func(*http.Request, *http.Client) (string, string, string, string, bool, error))(req, client)
    if err != nil {
        fmt.Println("Error running plugin:", err)
        os.Exit(1)
    }

	// Parse output if vulnerability is found
	if found {
        // TODO: Arrumar o filename, está vindo o diretório e nao o nome do template
		fileExt := filepath.Ext(scannerPlugin)
		scanName := scannerPlugin[:len(scannerPlugin)-len(fileExt)]
		
		parsers.VulnerabilityOutput(scanName, severity, url, rawReq, desc)
	}
}