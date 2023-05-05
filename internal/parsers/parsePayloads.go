package parsers

import (
	"strings"
	"fmt"

	"github.com/cfsdes/nucke/internal/initializers"
)

func ParsePayloads(payloads []string) []string {

	// Replace custom parameters
    if len(initializers.CustomParams) > 0 {
        for _, param := range initializers.CustomParams {
            parts := strings.SplitN(param, "=", -1)
            if len(parts) >= 2 {
                key := fmt.Sprintf("{{.%s}}", strings.TrimSpace(parts[0]))
                value := strings.TrimSpace(strings.Join(parts[1:], "="))

				// ...
				for i, s := range payloads {
					// Verify if element has string "{{.oob}}"
					if strings.Contains(s, key) {
						payloads[i] = strings.ReplaceAll(s, key, value)
					}
				}
            }
        }
    }

	// Replace {{.oob}} with interactsh url
	payloads = initializers.ReplaceOob(payloads)

	return payloads
}
