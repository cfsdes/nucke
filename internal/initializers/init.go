package initializers

import (
	"os"
	"fmt"

	"github.com/cfsdes/nucke/internal/globals"
	"github.com/cfsdes/nucke/pkg/plugins/utils"
	"github.com/cfsdes/nucke/pkg/report"
)

// Esse código será o primeiro executado no projeto. Irá carregar tudo que precisamos
func Start(version string){
	// Initial banner
	Banner()

	// Print Nucke version
    if globals.Version {
		fmt.Println("\nNucke version: ",version, "\n")
		os.Exit(0)
	}

	// Set Config Plugins
	if globals.PluginsConfig != "" {
		// Parse Config.yaml and Build Plugins
		globals.Scope = ParseConfig(globals.PluginsConfig)

		// Start interact.sh
		globals.InteractURL = utils.StartInteractsh()
	}

	// Create Output Path
	if globals.Output != "" {
		globals.Output = report.FormatOutput(globals.Output)
		err := os.MkdirAll(globals.Output, 0755)
		if err != nil {
			fmt.Println("Error creating output path:",err)
		}
	}

	// Check binaries
    binaries := []string{"interactsh-client"}
    CheckBinaries(binaries)
}