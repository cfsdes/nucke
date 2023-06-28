package initializers

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"

    "github.com/fatih/color"

)

// Compila cada plugin contido no filePaths
func BuildPlugins(filePaths []string) ([]string) {

    // Array onde ficarão os caminhos para os arquivos .so
    var pluginsPath []string

    // compila cada arquivo .go e salva o arquivo .so em ~/.nucke/compiled-plugins/
    for _, dir := range filePaths {
        // Pegando diretório do plugin.go
        compileDir := filepath.Dir(dir)

        // Compilando os plugins no diretório de cada um
        soFile := filepath.Join(compileDir, "."+filepath.Base(dir[:len(dir)-3])+".so")
        cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", soFile, dir)
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        err := cmd.Run()
        if err != nil {
            panic(fmt.Sprintf("Error compiling plugin %s: %v\n", dir, err))
        }

        pluginsPath = append(pluginsPath, soFile)
    }

    // return plugin .so paths
    return pluginsPath
}


// Check what plugins were loaded and return debug or error message
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

	if Debug {
		Blue := color.New(color.FgBlue, color.Bold).SprintFunc()
		fmt.Printf("[%s] Plugins loaded: %v\n", Blue("DEBUG"), fileNames)
	}
}

func extractFileNames(dirs []string) []string {
	var fileNames []string

	for _, dir := range dirs {
		fileName := strings.TrimSuffix(filepath.Base(dir), filepath.Ext(dir))
		fileNames = append(fileNames, fileName)
	}

	return fileNames
}