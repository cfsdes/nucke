package plugins

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func extractFileNames(dirs []string) []string {
	var fileNames []string

	for _, dir := range dirs {
		fileName := strings.TrimSuffix(filepath.Base(dir), filepath.Ext(dir))
		fileNames = append(fileNames, fileName)
	}

	return fileNames
}

func CheckLoadedPlugins(filePaths []string, plugins []string) {
	fileNames := extractFileNames(filePaths)

	missingItems := []string{}
	for _, plugin := range plugins {
		found := false
		for _, fileName := range fileNames {
			if fileName == plugin {
				found = true
				break
			}
		}
		if !found {
			missingItems = append(missingItems, plugin)
		}
	}

	

	if len(missingItems) != 0 {
		Red := color.New(color.FgRed, color.Bold).SprintFunc()
		fmt.Printf("[%s] Error loading plugins: %v\n", Red("ERR"), missingItems)
	} 

	if (initializers.Debug) {
		Blue := color.New(color.FgBlue, color.Bold).SprintFunc()
		fmt.Printf("[%s] Plugins loaded: %v\n", Blue("DEBUG"), fileNames)
	}
}
