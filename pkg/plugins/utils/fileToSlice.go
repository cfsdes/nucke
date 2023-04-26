package utils

import (
    "os"
    "bufio"
)

// Read file and return a slice
func FileToSlice(pluginDir string, rulesFile string) ([]string) {
	
    // if rulesFile is empty, use the default name
	if rulesFile == "" {
		rulesFile = "regex_match.txt"
	}

    // Open the file for reading
    file, err := os.Open(pluginDir + "/" + rulesFile)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    // Create an empty []string array
    var lines []string

    // Create a scanner to read from the file
    scanner := bufio.NewScanner(file)

    // Iterate over the lines in the file
    for scanner.Scan() {
        // Add the line to the []string array
        lines = append(lines, scanner.Text())
    }

    return lines
}