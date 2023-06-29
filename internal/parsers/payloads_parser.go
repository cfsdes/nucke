package parsers

import (
	"strings"
	"fmt"

	"github.com/cfsdes/nucke/pkg/plugins/utils"
	"github.com/cfsdes/nucke/pkg/globals"
)

func ParsePayload(payload string) string {

	// Replace custom parameters
    if len(globals.CustomParams) > 0 {
        for _, param := range globals.CustomParams {
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
	payload = utils.ReplaceOob(payload)

	return payload
}
