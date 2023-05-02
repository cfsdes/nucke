package initializers

import (
	"fmt"
    "flag"
    "os"

    "github.com/fatih/color"
)

func ParseFlags() (port string, threads int, jcAPI string, jc bool, scope string, proxy string, config string, output string, updatePlugins bool, exportCA bool, debug bool, verbose bool) {
	flag.StringVar(&port, "port", "8888", "proxy port to use")
    flag.IntVar(&threads, "threads", 8, "threads to use during plugin scan")
    flag.StringVar(&jcAPI, "jc-api", "http://127.0.0.1:5000", "jaeles API server")
    flag.BoolVar(&jc, "jc", false, "enable jaeles proxy")
	flag.StringVar(&scope, "scope", "", "regex for scope")
    flag.StringVar(&proxy, "proxy", "", "http proxy to use during scans")
    flag.StringVar(&config, "config", "", "yaml config file with plugins to scan")
    flag.StringVar(&output, "out", "", "output directory to save scan results")
    flag.BoolVar(&updatePlugins, "update-plugins", false, "Force the build of all plugins")
    flag.BoolVar(&exportCA, "export-ca", false, "Export proxy PEM certificate")
    flag.BoolVar(&debug, "debug", false, "Return debug error messages")
    flag.BoolVar(&verbose, "v", false, "Verbose output")

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
        "Proxy": []string{"port"},
        "Jaeles": []string{"jc", "jc-api"},
        "Scan": []string{"config", "proxy", "threads", "out"},
        "Misc": []string{"update-plugins", "export-ca", "debug", "v"},
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
