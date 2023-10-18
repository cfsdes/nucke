package requests

import (
	"net/http"
    "fmt"

    "github.com/cfsdes/nucke/pkg/globals"
)

/**
* Make a request based on the http.Request object received.
* Return all responses in a slice format. One response for each redirect + the final response
*/

func Do(req *http.Request, client *http.Client) ([]*http.Response) {
	responses := []*http.Response{}
	redirectLimit := 10 // Define o limite de redirecionamentos

	for {
		response, err := client.Do(req)
		if err != nil {
			if globals.Debug {
				fmt.Println("do_request:",err)
			}
			return responses
		}

		responses = append(responses, response)

		if len(responses) > redirectLimit || response.StatusCode < 300 || response.StatusCode >= 400 {
			break
		}

		nextURL, err := response.Location()
		if err != nil {
			break
		}

		// Crie uma nova solicitação com o URL do redirecionamento
		req, err = http.NewRequest("GET", nextURL.String(), nil)
		if err != nil {
			break
		}

		// Defina manualmente o cabeçalho Host para o domínio do URL do redirecionamento
		req.Host = nextURL.Host

		// Configure os cookies na nova solicitação
		cookies := client.Jar.Cookies(req.URL)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
	}

	return responses
}