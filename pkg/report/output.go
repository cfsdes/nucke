package report

import (
	"fmt"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"net/url"
	"io/ioutil"
	"path/filepath"
	"strings"
	"os/user"
	"os"

	"github.com/fatih/color"
	"github.com/cfsdes/nucke/pkg/globals"
)

func VulnerabilityOutput(scanName string, severity string, url string, summary string) {
	yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
	red := color.New(color.FgRed, color.Bold).SprintFunc()
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	blue := color.New(color.FgBlue, color.Bold).SprintFunc()
	magenta := color.New(color.FgMagenta, color.Bold).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	switch severity {
	case "Critical":
		fmt.Printf("[%s] [%s] %s \n", cyan(scanName), magenta(severity), url)
	case "High":
		fmt.Printf("[%s] [%s] %s \n", cyan(scanName), red(severity), url)
	case "Medium":
		fmt.Printf("[%s] [%s] %s \n", cyan(scanName), yellow(severity), url)
	case "Low":
		fmt.Printf("[%s] [%s] %s \n", cyan(scanName), green(severity), url)
	case "Info":
		fmt.Printf("[%s] [%s] %s \n", cyan(scanName), blue(severity), url)
	}

	// Create output
	if globals.Output != "" {
		outputPath := getOutputPath(url, scanName, summary)

		err := writeStringToFile(summary, outputPath)
		if err != nil {
			panic(err)
		}
	}
}

// Format output
func FormatOutput(output string) string {
	// Get the current user's home directory
    usr, err := user.Current()
    if err != nil {
        fmt.Println("Error getting current user:", err)
        os.Exit(1)
    }

    // Replace "~" with the home directory in the plugin path
    outputPath := strings.Replace(output, "~", usr.HomeDir, 1) // path/to/plugin/plugin.so
	
	return outputPath
}

// Func to generate random string
func generateRandomString(length int) string {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(randomBytes)
}

// return the domain of the url
func getDomain(urlString string) string {
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		panic(err)
	}
	return parsedURL.Hostname()
}

// Function to return the outputPath
func getOutputPath(urlString string, scanName string, summary string) string {
	// parse domain
	domain := getDomain(urlString)

	// generate random SHA1 hash
	signatureString := fmt.Sprintf("%s-%s-%s", len(summary), urlString, scanName)
	hash := sha1.Sum([]byte(signatureString))
	hashString := hex.EncodeToString(hash[:])

	// filename
	fileName := scanName + "-" + hashString

	// set output path
	outputPath := filepath.Join(globals.Output, domain, fileName)
	return outputPath
}

// Function to write string to file
func writeStringToFile(text string, filePath string) error {
	// Check if file already exists
	if _, err := os.Stat(filePath); err == nil {
		// File already exists, return without writing to it
		return nil
	}

	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		return err
	}
	
	err = ioutil.WriteFile(filePath, []byte(text), 0644)
	if err != nil {
		return err
	}
	return nil
}


