package main

import (
	"fmt"
	"log"
	"net/http"
    "flag"
    "os"

    "github.com/fatih/color"
    "github.com/cfsdes/nucke/internal/runner"
)


func parseFlags() (port int, jcAPI string) {
	flag.IntVar(&port, "port", 8080, "port number to use")
    flag.StringVar(&jcAPI, "jc-api", "http://127.0.0.1:5000", "jcAPI value to send to handler")
	
    // Add the welcome message to the --help output
	flag.Usage = func() {
		initialMessage()
        color.Cyan("Usage: \n")
        fmt.Fprintf(flag.CommandLine.Output(), "  %s [flags]\n\n", os.Args[0])
        color.Cyan("Flags: \n")
		flag.PrintDefaults()
	}

    flag.Parse()
	return
}

func initialMessage() {
    // Print a colorful welcome message
    fmt.Println()
    color.Blue("Welcome to Nucke Server!")
    fmt.Println()

    color.Yellow(`
        ,--,
  _ ___/ /\|
 ;( )__, )
; //   '--;
  \     |
   ^    ^`)

    fmt.Println()
}

func main() {
    port, jaelesApi := parseFlags()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		runner.Handler(w, r, jaelesApi)
	})

    server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.DefaultServeMux,
	}

	
    initialMessage()
	color.Cyan("Listening on port %d...\n", port)
    color.Cyan("Interacting with jaeles: %s\n", jaelesApi)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}