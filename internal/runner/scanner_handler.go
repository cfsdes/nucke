package runner

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"plugin"
	"strings"
	"time"

	"github.com/cfsdes/nucke/pkg/globals"
	"github.com/cfsdes/nucke/pkg/report"
	"github.com/cfsdes/nucke/pkg/requests"
)

func ScannerHandler(req *http.Request) {
	// Create HTTP Client
	client, err := createHTTPClient()
	if err != nil {
		fmt.Println("ScannerHandler:", err)
	}

	// Run Config Plugins
	for _, plugin := range globals.PluginPaths {
		runPlugin(plugin, req, client)
	}
}

// Generate HTTP Client with Proxy
func createHTTPClient() (*http.Client, error) {
	// Configure jar cookies to HTTP client
	jar, _ := cookiejar.New(nil)
	var client *http.Client

	if globals.Proxy != "" {
		// Create HTTP client with proxy
		proxyUrl, err := url.Parse(globals.Proxy)
		if err != nil {
			return nil, fmt.Errorf("failed to parse proxy URL: %s", err)
		}
		client = &http.Client{
			Transport: &http.Transport{
				DisableKeepAlives: true,
				TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
				Proxy:             http.ProxyURL(proxyUrl),
			},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				// Don't allow automatic redirects
				return http.ErrUseLastResponse
			},
			Timeout: time.Second * 240, // 4min timeout
			Jar:     jar,
		}
	} else {
		// Create HTTP client without proxy
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
				DisableKeepAlives: true,
			},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				// Don't allow automatic redirects
				return http.ErrUseLastResponse
			},
			Timeout: time.Second * 240, // 4min timeout
			Jar:     jar,
		}
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
	pluginPath := strings.Replace(scannerPlugin, "~", usr.HomeDir, 1) // path/to/plugin/plugin.so
	pluginDir := filepath.Dir(pluginPath)                             // path/to/plugin

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
	scanName := filepath.Base(pluginPath)
	scanName = strings.TrimSuffix(scanName, filepath.Ext(scanName))
	scanName = strings.TrimPrefix(scanName, ".") // e.g: sqli

	// Call the run() function with a req argument
	found, severity, url, payload, param, rawReq, rawResp := runFunc.(func(*http.Request, *http.Client, string) (bool, string, string, string, string, string, string))(req, client, pluginDir)

	// Get Response Status Code
	resStatusCode := requests.StatusCodeFromRaw(rawResp)

	// Parse output if vulnerability is found
	if found && resStatusCode != 429 {
		var webhook string
		webhook = report.GetWebhook(scannerPlugin)
		report.Output(scanName, webhook, severity, url, payload, param, rawReq, rawResp, pluginDir)
	}
}
