package runner

import (
	"net/http"
	"net/url"
	"fmt"

	"github.com/cfsdes/nucke/internal/utils"
)

func ScannerHandler(r *http.Request) {
	// Create HTTP Client
	_, err := createHTTPClient()
	if err != nil {
		fmt.Println(err)
	}

	/*
		Vai usar o parser yaml (isso deve ser feito no começo na globals)
		Pegar todos os plugins .so
		Para cada plugin, chamar a função run() deles
		Analisar pela resposta da função run(), 
		se retornar true, vulnerabilidade foi encontrada
	*/

	// Handle Config Plugins
	/*
	for _, vuln := range vulnsList {
		switch vuln {
		case "sqli":
			_, err := scanners.SqliQuery(r, client)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	*/
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
