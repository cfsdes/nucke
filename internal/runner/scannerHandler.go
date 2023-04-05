package runner

import (
	"net/http"
	"net/url"
	"fmt"

	"github.com/cfsdes/nucke/internal/helpers"
	"github.com/cfsdes/nucke/internal/scanners"
)

func ScannerHandler(r *http.Request, vulnsList []string) {
	// Create HTTP Client
	client, err := createHTTPClient()
	if err != nil {
		fmt.Println(err)
	}

	// Loop vulnsList
	for _, vuln := range vulnsList {
		switch vuln {
		case "sqli":
			_, err := scanners.SqliQuery(r, client)
			if err != nil {
				fmt.Println(err)
			}
		case "xss-script":
			//XssScript()
		case "path-traversal":
			//PathTraversal()
		}
	}
}

// Generate HTTP Client with Proxy
func createHTTPClient() (*http.Client, error) {
    var client *http.Client
    if helpers.Proxy != "" {
        // Create HTTP client with proxy
        proxyUrl, err := url.Parse(helpers.Proxy)
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
