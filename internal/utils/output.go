package utils

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
)

func VulnerabilityOutput(scanName string, severity string, url string, summary string) {
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()
	//cyan := color.New(color.FgCyan).SprintFunc()
	white := color.New(color.FgWhite).SprintFunc()

	switch severity {
	case "Critical":
		fmt.Printf("[%s] [%s] %s \n", magenta(scanName), magenta(severity), white(url))
	case "High":
		fmt.Printf("[%s] [%s] %s \n", red(scanName), red(severity), white(url))
	case "Medium":
		fmt.Printf("[%s] [%s] %s \n", yellow(scanName), yellow(severity), white(url))
	case "Low":
		fmt.Printf("[%s] [%s] %s \n", green(scanName), green(severity), white(url))
	case "Info":
		fmt.Printf("[%s] [%s] %s \n", blue(scanName), blue(severity), white(url))
	}

	// Create output
	if Output != "" {
		outputPath := getOutputPath(url, scanName)

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
func getOutputPath(urlString string, scanName string) string {
	// parse domain
	domain := getDomain(urlString)

	// generate random SHA1 hash
	randomString := generateRandomString(10)
	hash := sha1.Sum([]byte(randomString))
	hashString := hex.EncodeToString(hash[:])

	// filename
	fileName := scanName + "-" + hashString

	// set output path
	outputPath := filepath.Join(Output, domain, fileName)
	return outputPath
}

// Function to write string to file
func writeStringToFile(text string, filePath string) error {
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


