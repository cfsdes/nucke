package utils

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "regexp"
    "time"
    "strings"
    "math/rand"

	"github.com/fatih/color"
    "github.com/cfsdes/nucke/pkg/globals"
)

var interactOutput string = "/tmp/nucke-interact"
var interactSession string = "/tmp/nucke-interact-session"
var cmd *exec.Cmd

func StartInteractsh() {
    // Removing old session files
    deleteFileIfExists(interactSession)

	// Initial Message
    Cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	fmt.Printf("[%s] Starting interactsh...\n", Cyan("INF"))

    // Start interact first time
    startInteractSession()

    // restart interact in background every 1 hour
    go restartInteract()
}

// Start Interactsh
func startInteractSession() {

    // Start interactsh client and save session file
    cmd = exec.Command("interactsh-client", "-sf", interactSession)

    // Create output file
    deleteFileIfExists(interactOutput)
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

    // Wait for 5 seconds
    time.Sleep(5 * time.Second)

    // Read the interactsh initial output
    for {
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
                globals.InteractURL = match// Exit the loop if match is found
                return
            }
        }

        if err := scanner.Err(); err != nil {
            fmt.Println("Error reading file:", err)
            os.Exit(1)
        }

        if globals.Debug {
            Blue := color.New(color.FgBlue, color.Bold).SprintFunc()
            fmt.Printf("[%s] Error getting interactsh URL, trying again...\n", Blue("DEBUG"))
        }

        // Wait for 5 seconds before retrying
        time.Sleep(5 * time.Second)
    }
}

// restart interactsh process every 1 hour
func restartInteract() {
    for {
        time.Sleep(1 * time.Hour)

        // Send SIGINT signal to interrupt interactsh process
        err := cmd.Process.Signal(os.Interrupt)
        if err != nil {
            fmt.Println("Error sending SIGINT signal to interactsh process:", err)
            continue
        }

        // espera 5 segundos para o processo interactsh terminar
        time.Sleep(5 * time.Second)

        // kill interactsh process
        err = cmd.Process.Kill()
        if err != nil {
            fmt.Println("Error killing interactsh process:", err)
        }

        // Restart interactsh
        startInteractSession()
    }
}

// Replace {{.oob}} with interactSH URL
func ReplaceOob(payload string) string {
    if strings.Contains(payload, "{{.oob}}") {
        // Replace "{{.oob}}" with random ID + InteractURL
        id := fmt.Sprintf("%08d", rand.Intn(100000000))
        payload = strings.ReplaceAll(payload, "{{.oob}}", id+"."+globals.InteractURL)
    }

	return payload
}

// Extract OOB ID
func ExtractOobID(url string) string {
    re := regexp.MustCompile(`\d{8}\.`)
    match := re.FindString(url)
    if match == "" {
        return ""
    }
    id := match[:8]
    return id
}

// Function to check interactsh interaction
func CheckOobInteraction(oobID string) bool {

    // Wait some seconds before analyze
    time.Sleep(15 * time.Second)

    // Read the interactsh output
    file, err := os.Open(interactOutput)
    if err != nil {
        fmt.Println("Error opening file:", err)
        os.Exit(1)
    }
    defer file.Close()

    // Parse the initial output to grep the OOB URL
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if strings.Contains(line, oobID) {
            return true
        }
    }

    if err := scanner.Err(); err != nil {
        fmt.Println("Error reading file:", err)
        os.Exit(1)
    }

    return false
}


func deleteFileIfExists(filePath string) {
    if _, err := os.Stat(filePath); err == nil {
        if err := os.Remove(filePath); err != nil {
            fmt.Println("Error deleting interactsh file:", err)
        }
    }
}
