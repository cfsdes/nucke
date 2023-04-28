package runner

import (
	"encoding/base64"
	"net/http"
	"net/http/httputil"
    "regexp"
    "fmt"
    "log"
    "github.com/elazarl/goproxy"
    "os"
    "encoding/pem"

    "github.com/fatih/color"
    "github.com/cfsdes/nucke/internal/utils"
    "github.com/cfsdes/nucke/internal/parsers"
    "github.com/cfsdes/nucke/pkg/requests"
)

// Create a channel with a buffer of threads
var ch = make(chan int, utils.Threads)

// Start Proxy
func StartProxy() {

    // Export CA certificate
    if utils.ExportCA {
        exportCA()
        return
    }

    // Cria um proxy com a função de roteamento personalizada
    proxy := goproxy.NewProxyHttpServer()
    //proxy.Verbose = true
    proxy.OnRequest().DoFunc(requestHandler)
    proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

    // Start messages
    color.Cyan("Listening on port %s...\n", utils.Port)

    if utils.Jaeles {
        color.Cyan("Interacting with jaeles: %s\n", utils.JaelesApi)
    }

    fmt.Println()

    // Start to listen
    log.Fatal(http.ListenAndServe(":"+utils.Port, proxy))

}

// Proxy Handler
func requestHandler(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
    // Convert the raw request to base64
	requestBytes, err := httputil.DumpRequest(req, true)
    if err != nil {
        fmt.Println("requestHandler: Error converting rawRequest: ",err)
    }
    requestBase64 := base64.StdEncoding.EncodeToString(requestBytes)

    
    // Send request to jaeles API server and filter if scope is specified
    if (utils.Scope != "" && regexp.MustCompile(utils.Scope).MatchString(req.URL.String()) || utils.Scope == "") {
        
        // If jaeles scan is enabled
        if utils.Jaeles {
            parsers.SendToJaeles(requestBase64, utils.JaelesApi)
        }

        // If config with plugins is provided
        if utils.Config != "" {
            // Clone request before scanning
            req2 := requests.CloneReq(req)
            
            // executa a ScannerHandler dentro de uma goroutine
            go func() {
                // adiciona 1 ao canal para indicar que está utilizando uma goroutine
                ch <- 1

                ScannerHandler(req2)

                // sinaliza ao canal que a goroutine está livre
                <-ch
            }()
        }
	} 

    // Defina a lógica para lidar com a requisição aqui
    return req, nil
}

// Function to export CA certificates
func exportCA() {
    color.Cyan("CA certificate exported to local path: nucke-cert.pem\n\n")

    // Criar o arquivo cert.pem
    file, err := os.Create("nucke-cert.pem")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer file.Close()

    // Obter o certificado X.509 da propriedade Certificate
    cert := goproxy.GoproxyCa.Certificate[0]

    // Codificar o certificado em formato PEM e escrever no arquivo
    err = pem.Encode(file, &pem.Block{Type: "CERTIFICATE", Bytes: cert})
    if err != nil {
        fmt.Println(err)
        return
    }
}