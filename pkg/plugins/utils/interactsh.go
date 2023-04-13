package utils

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "regexp"
    "time"
    "strings"

	"github.com/fatih/color"
)

var interactOutput string = "/tmp/nucke-interact"
var interactSession string = "/tmp/nucke-interact-session"
var interactURL string

func StartInteractsh() string {
	// Initial Message
	color.Magenta("Starting interactsh...\n")

	// Start interactsh client and save session file
    cmd := exec.Command("interactsh-client", "-sf", interactSession)
    file, err := os.Create(interactOutput)
    if err != nil {
		fmt.Println("Error creating interact file:", err)
		os.Exit(1)
	}
	defer file.Close()
    cmd.Stdout = file
	cmd.Stderr = file

    err = cmd.Start()
    if err != nil {
        fmt.Println("Error executing interactsh command:", err)
        os.Exit(1)
    }

	// Wait for 5 seconds before interrupt the interactsh process
    time.Sleep(5 * time.Second)

    if err := cmd.Process.Signal(os.Interrupt); err != nil {
        fmt.Println("Error to interrupt interactsh process:", err)
        os.Exit(1)
    }

	// Read the interactsh initial output
    file, err = os.Open(interactOutput)
    if err != nil {
        fmt.Println("Error opening file:", err)
        os.Exit(1)
    }
    defer file.Close()

	// Parse the initial output to grep the OOB URL
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        re := regexp.MustCompile(`[a-zA-Z0-9]+\.oast\.[a-zA-Z0-9]+`)
        match := re.FindString(line)
        if match != "" {
            interactURL = match
			return interactURL
        }
    }

    if err := scanner.Err(); err != nil {
        fmt.Println("Error reading file:", err)
        os.Exit(1)
    }

	return ""
}

// Replace {{.oob}} with interactSH URL
func ReplaceOob(arr []string) []string {
	for i, s := range arr {
		// Verify if element has string "{{.oob}}"
		if strings.Contains(s, "{{.oob}}") {
			// replace string "{{.oob}}" with interactURL
			arr[i] = strings.ReplaceAll(s, "{{.oob}}", interactURL)
		}
	}

	return arr
}