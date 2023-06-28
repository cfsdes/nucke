package plugins

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
)

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


