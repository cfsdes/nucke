package runner

import (
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync/atomic"

	"github.com/elazarl/goproxy"

	"github.com/cfsdes/nucke/pkg/globals"
	"github.com/cfsdes/nucke/pkg/requests"
	"github.com/fatih/color"
)

// Create a channel with a buffer of threads
var ch = make(chan int, globals.Threads)

// Start Proxy
func StartProxy() {

	// Export CA certificate
	if globals.ExportCA {
		exportCA()
		return
	}

	// Cria um proxy com a função de roteamento personalizada
	proxy := goproxy.NewProxyHttpServer()

	// Desabilitar logs da biblioteca goproxy
	logger := log.New(ioutil.Discard, "", 0)
	proxy.Logger = logger
	proxy.Verbose = false

	proxy.OnRequest().DoFunc(requestHandler)
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	// Start messages
	Green := color.New(color.FgGreen, color.Bold).SprintFunc()
	fmt.Printf("[%s] Nucke running on port %s!\n", Green("OK"), globals.Port)

	// Start Status Server
	if globals.Stats {
		InitStatsServer()
	}

	fmt.Println()

	// Start to listen
	log.Fatal(http.ListenAndServe(":"+globals.Port, proxy))

}

// Proxy Handler
func requestHandler(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {

	// Send request and run scan. Filter scope if needed
	if globals.Scope != "" && regexp.MustCompile(globals.Scope).MatchString(req.URL.String()) || globals.Scope == "" {

		// If verbose
		if globals.Debug {
			Blue := color.New(color.FgBlue, color.Bold).SprintFunc()
			fmt.Printf("[%s] Scanning: %s %s\n", Blue("DEBUG"), req.Method, req.URL.String())
		}

		// If config with plugins is provided
		if globals.PluginsConfig != "" {
			// Clone request before scanning
			reqScan := requests.CloneReq(req)

			// Add request to pendingRequests
			atomic.AddInt64(&globals.PendingScans, 1)

			// executa a ScannerHandler dentro de uma goroutine
			go func() {
				// adiciona 1 ao canal para indicar que está utilizando uma goroutine
				ch <- 1

				// Ensure the decrement is executed even if an error occurs during scan
				defer func() {
					// Remove request from pendingRequests
					atomic.AddInt64(&globals.PendingScans, -1)

					// sinaliza ao canal que a goroutine está livre
					<-ch
				}()

				ScannerHandler(reqScan)
			}()
		}
	}

	// Clone original request
	reqCopy := requests.CloneReq(req)

	// Defina a lógica para lidar com a requisição aqui
	return reqCopy, nil
}

// Function to export CA certificates
func exportCA() {
	color.Cyan("CA certificate exported to local path: nucke-cert.crt\n\n")

	// Criar o arquivo cert.pem
	file, err := os.Create("nucke-cert.crt")
	if err != nil {
		fmt.Println("Export CA:", err)
		return
	}
	defer file.Close()

	// Obter o certificado X.509 da propriedade Certificate
	cert := goproxy.GoproxyCa.Certificate[0]

	// Codificar o certificado em formato PEM e escrever no arquivo
	err = pem.Encode(file, &pem.Block{Type: "CERTIFICATE", Bytes: cert})
	if err != nil {
		fmt.Println("Export CA:", err)
		return
	}
}
