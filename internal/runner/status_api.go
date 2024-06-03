package runner

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/fatih/color"

	"github.com/cfsdes/nucke/pkg/globals"
)

func InitStatsServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		count := globals.PendingScans
		if count > 0 {
			fmt.Fprintf(w, "Pending scans: %d\n", count)
		} else {
			fmt.Fprintf(w, "No pending scans!\n")
		}
		fmt.Fprintf(w, "Requests sent: %d\n", globals.RequestsMade)
	})

	// Start messages
	Green := color.New(color.FgGreen, color.Bold).SprintFunc()
	fmt.Printf("[%s] Status server started on port 8899\n", Green("OK"))

	go http.ListenAndServe(":8899", nil)

}
