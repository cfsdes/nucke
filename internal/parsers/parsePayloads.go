package parsers

import (
	"strings"
	"fmt"

	"github.com/cfsdes/nucke/internal/initializers"
)

func ParsePayload(payload string) string {

	// Replace custom parameters
    if len(initializers.CustomParams) > 0 {
        for _, param := range initializers.CustomParams {
            parts := strings.SplitN(param, "=", -1)
            if len(parts) >= 2 {
                key := fmt.Sprintf("{{.%s}}", strings.TrimSpace(parts[0]))
                value := strings.TrimSpace(strings.Join(parts[1:], "="))

				// Verify if element has string "{{.oob}}"
				if strings.Contains(payload, key) {
					payload = strings.ReplaceAll(payload, key, value)
				}
            }
        }
    }

	// Replace {{.oob}} with interactsh url
	payload = initializers.ReplaceOob(payload)

	return payload
}
