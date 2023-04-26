package plugins

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
	"strings"
)

func BuildPlugins(filePaths []string, updatePlugins bool) ([]string) {
    // cria o diretório ~/.nucke/compiled-plugins/ se ele não existir
    compileDir := filepath.Join(os.Getenv("HOME"), ".nucke", "compiled-plugins")
    if _, err := os.Stat(compileDir); os.IsNotExist(err) {
        if err := os.MkdirAll(compileDir, 0755); err != nil {
            panic(fmt.Sprintf("Error creating directory %s: %v\n", compileDir, err))
        }
    }

    // compila cada arquivo .go e salva o arquivo .so em ~/.nucke/compiled-plugins/
    for _, dir := range filePaths {
        soFile := filepath.Join(compileDir, filepath.Base(dir[:len(dir)-3])+".so")
        if _, err := os.Stat(soFile); os.IsNotExist(err) || updatePlugins {
            cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", soFile, dir)
            cmd.Stdout = os.Stdout
            cmd.Stderr = os.Stderr
            err := cmd.Run()
            if err != nil {
                panic(fmt.Sprintf("Error compiling plugin %s: %v\n", dir, err))
            }
        }
    }

    // return plugin .so paths
    return createPluginsPath(filePaths)
}


func createPluginsPath(filePaths []string) ([]string) {
	var pluginsPath []string
	for _, path := range filePaths {
		// obtem apenas o nome do arquivo sem a extensão
		fileName := strings.TrimSuffix(filepath.Base(path), ".go")

		// concatena com o caminho para a pasta de plugins compilados
		pluginPath := fmt.Sprintf("~/.nucke/compiled-plugins/%s.so", fileName)

		// adiciona o caminho do plugin à lista
		pluginsPath = append(pluginsPath, pluginPath)
	}

	return pluginsPath
}
