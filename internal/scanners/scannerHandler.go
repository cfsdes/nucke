package scanners

import (
	//"net/http"
	"fmt"
)

func ScannerHandler(vulnsList []string) {
	for _, vuln := range vulnsList {
		switch vuln {
		case "sqli-query":
			//SqliQuery()
		case "xss-script":
			//XssScript()
		case "path-traversal":
			//PathTraversal()
		default:
			fmt.Printf("Error: unknown vulnerability %s\n", vuln)
		}
	}
}