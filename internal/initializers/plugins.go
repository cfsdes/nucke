package initializers

import (
    "fmt"
    "os"
    "os/exec"
	"os/user"
    "path/filepath"
    "strings"

    "github.com/fatih/color"
	"github.com/cfsdes/nucke/pkg/globals"

)

// Compila cada plugin contido no filePaths
func BuildPlugins(filePaths []string) ([]string) {

	Cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	fmt.Printf("[%s] Building plugins...\n", Cyan("INF"))

    // Array onde ficarão os caminhos para os arquivos .so
    var pluginsPath []string

    // compila cada arquivo .go e salva o arquivo .so em ~/.nucke/compiled-plugins/
    for _, dir := range filePaths {
        // Pegando diretório do plugin.go
        compileDir := filepath.Dir(dir)

		// Alterando para o diretório mainDir
		os.Chdir(compileDir)

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
		if plugin == "*" {
			found = true
		}
		if !found {
			missingItems = append(missingItems, plugin)
		}
	}


	if len(missingItems) != 0 {
		Red := color.New(color.FgRed, color.Bold).SprintFunc()
		fmt.Printf("[%s] Error loading plugins: %v\n", Red("ERR"), missingItems)
	} 

	if globals.Debug {
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

// List all plugins in the filePath (dir)
func ListPlugins(dir string){
	expandedDir, err := expandUser(dir)
	if err != nil && globals.Debug {
		fmt.Println("List plugins error:", err)
		return
	}

	Cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	fmt.Printf("[%s] Plugins available on %s:\n\n", Cyan("INF"), expandedDir)
	err = filepath.Walk(expandedDir, func(path string, info os.FileInfo, err error) error {
		if err != nil && globals.Debug {
			fmt.Println("List plugins error:", err)
			return err
		}

		if !info.IsDir() && filepath.Ext(info.Name()) == ".go" {
			fileName := filepath.Base(info.Name())
			fileName = fileName[:len(fileName)-len(filepath.Ext(fileName))]
			fmt.Println(fileName)
		}

		return nil
	})
	fmt.Println("\n")

	if err != nil && globals.Debug {
		fmt.Println("List plugins error:", err)
	}

	os.Exit(0)
}

func expandUser(path string) (string, error) {
	if len(path) < 2 || path[:2] != "~/" {
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return filepath.Join(usr.HomeDir, path[2:]), nil
}