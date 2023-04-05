package helpers

import (
	"strings"
	"fmt"
)

var VulnList = []string{
	"sqli",
	"xss-script",
	"path-traversal",
}

// Check if the vulns provided are valid
func ValidateVulns(vulns string) (vulnArgs []string, err error) {
    if vulns != "" {
        vulnArgs = strings.Split(vulns, ",")
        for _, vuln := range vulnArgs {
            if !contains(VulnList, vuln) {
                return nil, fmt.Errorf("error: %s is not a valid vulnerability", vuln)
            }
        }
    }
    return vulnArgs, nil
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}