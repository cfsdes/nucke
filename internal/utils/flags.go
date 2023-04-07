package utils

import (
	"fmt"
    "flag"
    "os"

    "github.com/fatih/color"
)

func ParseFlags() (port int, jcAPI string, jc bool, scope string, proxy string, config string, output string) {
	flag.IntVar(&port, "port", 8080, "proxy port to use")
    flag.StringVar(&jcAPI, "jc-api", "http://127.0.0.1:5000", "jaeles API server")
    flag.BoolVar(&jc, "jc", false, "enable jaeles proxy")
	flag.StringVar(&scope, "scope", "", "regex for scope")
    flag.StringVar(&proxy, "proxy", "", "http proxy to use during scans")
    flag.StringVar(&config, "config", "", "yaml config file with plugins to scan")
    flag.StringVar(&output, "out", "", "output directory to save scan results")

    // Add the welcome message to the --help output
	flag.Usage = func() {
        color.Cyan("Usage: \n")
        fmt.Fprintf(flag.CommandLine.Output(), "  %s [flags]\n\n", os.Args[0])
        color.Cyan("Flags: \n")
		flag.PrintDefaults()
	}

    flag.Parse()

    return
}

