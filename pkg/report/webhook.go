package report

import (
	"fmt"
	"net/http"
	"bytes"
	"encoding/json"
	
	"github.com/fatih/color"
	"github.com/cfsdes/nucke/pkg/globals"
)

// Envia request JSON POST para o webhook sobre a vulnerabilidade identificada
func Notify(scanName string, severity string, url string, summary string, webhook string) {
	Red := color.New(color.FgRed, color.Bold).SprintFunc()
	Blue := color.New(color.FgBlue, color.Bold).SprintFunc()

	// Verificar se webhook não está vazio
	if webhook != "" {
        // Criar uma estrutura de dados para representar os parâmetros JSON
		data := map[string]string{
			"plugin":  scanName,
			"severity":  severity,
			"url":       url,
			"report":   summary,
		}

        // Converter a estrutura de dados para JSON
		jsonData, err := json.Marshal(data)
		if err != nil && globals.Debug {
			fmt.Printf("[%s] Output: Error converting JSON\n", Red("ERR"))
			return
		}

        // Enviar a requisição POST JSON para o webhook
		resp, err := http.Post(webhook, "application/json", bytes.NewBuffer(jsonData))
		if err != nil && globals.Debug {
			fmt.Printf("[%s] Output: Error sending webhook request\n", Red("ERR"))
			return
		}
		defer resp.Body.Close()

        // Verificar a resposta do servidor
		if resp.StatusCode == http.StatusOK {
			fmt.Printf("[%s] Webhook: Successful Request\n", Blue("DEBUG"))
		} else {
			fmt.Printf("[%s] Webhook: Error Sending Request\n", Red("ERR"))
		}
	}
}

// Identifica o webhook para o plugin sendo escaneado e retorna ele
func GetWebhook(pluginPath string) string {   
    // Loop externo para percorrer todas as chaves do mapa
	for key, files := range globals.Webhook {
		// Loop interno para percorrer os valores associados à chave
		for _, file := range files {
			// Verificar se o valor desejado está presente
			if file == pluginPath {
				return key
			}
		}
	}

    return ""
}