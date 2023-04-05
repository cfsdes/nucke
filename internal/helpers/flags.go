package helpers

import (
	"fmt"
    "flag"
    "os"

    "github.com/fatih/color"
)

func ParseFlags() (port int, jcAPI string, jc bool, scope string, listVulns bool, vulns string) {
	flag.IntVar(&port, "port", 8080, "proxy port to use")
    flag.StringVar(&jcAPI, "jc-api", "http://127.0.0.1:5000", "jaeles API server")
    flag.BoolVar(&jc, "jc", false, "enable jaeles proxy")
	flag.StringVar(&scope, "scope", "", "regex for scope")
    flag.BoolVar(&listVulns, "list-vulns", false, "list available vulnerabilities")
    flag.StringVar(&vulns, "vulns", "", "comma-separated list of vulnerabilities to scan. Use -list-vulns to get the options")

    // Add the welcome message to the --help output
	flag.Usage = func() {
		Banner()
        color.Cyan("Usage: \n")
        fmt.Fprintf(flag.CommandLine.Output(), "  %s [flags]\n\n", os.Args[0])
        color.Cyan("Flags: \n")
		flag.PrintDefaults()
	}

    flag.Parse()

    return
}

