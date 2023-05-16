package initializers

import (
	"fmt"
    "flag"
    "os"

    "github.com/fatih/color"
)

type headersFlag []string

func (h *headersFlag) String() string {
	return fmt.Sprintf("%v", *h)
}

func (h *headersFlag) Set(value string) error {
	*h = append(*h, value)
	return nil
}

func ParseFlags() (port string, threads int, jcAPI string, jc bool, scope string, proxy string, config string, output string, updatePlugins bool, exportCA bool, debug bool, verbose bool, stats bool, headers headersFlag, parameters headersFlag) {
	flag.StringVar(&port, "port", "8888", "proxy port to use (default: 8888)")
    flag.IntVar(&threads, "threads", 8, "threads to use during plugin scan (default: 8)")
    flag.StringVar(&jcAPI, "jc-api", "http://127.0.0.1:5000", "jaeles API server (default: http://127.0.0.1:5000)")
    flag.BoolVar(&jc, "jc", false, "enable jaeles proxy")
	flag.StringVar(&scope, "scope", "", "regex for scope")
    flag.StringVar(&proxy, "proxy", "", "http proxy to use during scans")
    flag.StringVar(&config, "config", "", "yaml config file with plugins to scan")
    flag.StringVar(&output, "out", "", "output directory to save scan results")
    flag.BoolVar(&updatePlugins, "update-plugins", false, "Force the build of all plugins")
    flag.BoolVar(&exportCA, "export-ca", false, "Export proxy PEM certificate")
    flag.BoolVar(&debug, "debug", false, "Return debug error messages")
    flag.BoolVar(&verbose, "v", false, "Verbose requests made")
    flag.BoolVar(&stats, "stats", false, "Start status server on port 8899")
    flag.Var(&headers, "headers", "Set custom headers. Accept multiple flag usages.") // Accept multiple flag usages
    flag.Var(&parameters, "p", "Custom parameters to be used in templates (e.g. -p \"dest=example.com\")") // Accept multiple flag usages

    // Add the welcome message to the --help output
	flag.Usage = func() {
        Cyan := color.New(color.FgCyan, color.Bold)
        Cyan.Printf("Usage: \n")
        fmt.Fprintf(flag.CommandLine.Output(), "  %s [flags]\n\n", os.Args[0])
        //Cyan.Printf("Flags: \n")
		PrintFlagsByTopic() // Imprime as flags por tópico
	}

    flag.Parse()

    return
}

func PrintFlagsByTopic() {
    Cyan := color.New(color.FgCyan, color.Bold)

    // Define os tópicos e as flags correspondentes
    topics := map[string][]string{
        "Proxy": []string{"port", "headers"},
        "Jaeles": []string{"jc", "jc-api"},
        "Scan": []string{"config", "proxy", "threads", "scope", "out", "p"},
        "Misc": []string{"update-plugins", "export-ca", "debug", "v", "stats"},
    }

    // Imprime as flags por tópico
    for topic, flags := range topics {
        Cyan.Printf("%s:\n", topic)
        for _, name := range flags {
            f := flag.Lookup(name)
            flagText := fmt.Sprintf("-%s", name)
            fmt.Printf("  %-25s %s\n", flagText, f.Usage)
        }
        fmt.Println()
    }
}
