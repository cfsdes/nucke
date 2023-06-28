package initializers

import (
    "log"
    "os/exec"
)

func CheckBinaries(binaries []string) {
    for _, binary := range binaries {
        cmd := exec.Command(binary, "--help")
        if err := cmd.Run(); err != nil {
            if exitError, ok := err.(*exec.ExitError); ok {
                if exitError.ExitCode() == 127 {
                    log.Fatal("command '" + binary + "' not found")
                }
            } else {
                log.Fatal("command '" + binary + "' not found")
            }
        }
    }
}
