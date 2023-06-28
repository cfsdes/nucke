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
    "io/ioutil"
    "sync/atomic"

    "github.com/fatih/color"
    "github.com/cfsdes/nucke/internal/globals"
    "github.com/cfsdes/nucke/internal/parsers"
    "github.com/cfsdes/nucke/pkg/requests"
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
    Cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
    fmt.Printf("[%s] Listening on port %s...\n", Cyan("INF"), globals.Port)

    if globals.Jaeles {
        fmt.Printf("[%s] Interacting with jaeles: %s\n", Cyan("INF"), globals.JaelesApi)
    }

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

    // Convert the raw request to base64
	requestBytes, err := httputil.DumpRequest(req, true)
    if err != nil {
        fmt.Println("requestHandler: Error converting rawRequest: ",err)
    }
    requestBase64 := base64.StdEncoding.EncodeToString(requestBytes)


    // Send request to jaeles API server and filter if scope is specified
    if (globals.Scope != "" && regexp.MustCompile(globals.Scope).MatchString(req.URL.String()) || globals.Scope == "") {
        
        // If verbose
        if globals.Debug {
            Blue := color.New(color.FgBlue, color.Bold).SprintFunc()
            fmt.Printf("[%s] Scanning: %s %s\n", Blue("DEBUG"), req.Method, req.URL.String())
        }

        // If jaeles scan is enabled
        if globals.Jaeles {
            parsers.SendToJaeles(requestBase64, globals.JaelesApi)
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

                ScannerHandler(reqScan)

                // Remove request from pendingRequests
                atomic.AddInt64(&globals.PendingScans, -1)
                
                // sinaliza ao canal que a goroutine está livre
                <-ch
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
        fmt.Println("Export CA:",err)
        return
    }
    defer file.Close()

    // Obter o certificado X.509 da propriedade Certificate
    cert := goproxy.GoproxyCa.Certificate[0]

    // Codificar o certificado em formato PEM e escrever no arquivo
    err = pem.Encode(file, &pem.Block{Type: "CERTIFICATE", Bytes: cert})
    if err != nil {
        fmt.Println("Export CA:",err)
        return
    }
}