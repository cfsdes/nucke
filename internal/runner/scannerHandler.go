package runner

import (
	"net/http"
	"net/url"
    "net/http/httputil"
	"fmt"
	"os"
    "plugin"
	"path/filepath"
    "os/user"
    "strings"
    "bytes"

	"github.com/cfsdes/nucke/internal/utils"
)

func ScannerHandler(req *http.Request, w http.ResponseWriter) {
	// Create HTTP Client
	client, err := createHTTPClient()
	if err != nil {
		fmt.Println(err)
	}

    // Create New Request based on Original Request
    newReq := createNewRequest(req, w)

	// Run Config Plugins
	for _, plugin := range utils.FilePaths {
		runPlugin(plugin, newReq, client)
	}
}

// Create a new request to forward
func createNewRequest(r *http.Request, w http.ResponseWriter) *http.Request {
    // Get request bytes
	requestBytes, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}
    
    // Generate new request
    newReq, err := http.NewRequest(r.Method, r.URL.String(), bytes.NewReader(requestBytes))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return nil
    }

    // Copy headers from original request to new request
    for key, values := range r.Header {
        for _, value := range values {
            newReq.Header.Add(key, value)
        }
    }
    
    return newReq
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
	severity, url, summary, found, err := runFunc.(func(*http.Request, *http.Client) (string, string, string, bool, error))(req, client)
    if err != nil {
        fmt.Println("Error running plugin:", err)
        os.Exit(1)
    }

	// Parse output if vulnerability is found
	if found {
		fileExt := filepath.Ext(scannerPlugin)
		scanName := filepath.Base(scannerPlugin[:len(scannerPlugin)-len(fileExt)])
		
		VulnerabilityOutput(scanName, severity, url, summary)
	}
}