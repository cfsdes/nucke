package runner

import (
	"fmt"
    "net/http"
	"github.com/fatih/color"

	"github.com/cfsdes/nucke/internal/initializers"

)


func InitStatsServer() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        count := initializers.PendingScans
		if count > 0 {
			fmt.Fprintf(w, "Pending scans: %d\n", count)
		} else {
			fmt.Fprintf(w, "No pending scans!\n")
		}
    })
    
	// Start messages
	Cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	fmt.Printf("[%s] Status server started on port 8899\n", Cyan("INF"))

	go http.ListenAndServe(":8899", nil)
	
}
