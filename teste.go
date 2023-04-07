package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"os/user"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"
)

type Plugin struct {
	Name    string   `yaml:"name"`
	Path    string   `yaml:"path"`
	Ids     []string `yaml:"ids"`
	Exclude []string `yaml:"exclude,omitempty"`
}

type Config struct {
	Plugins []Plugin `yaml:"plugins"`
}

func main() {
	// Read config file
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	// Create an array to store all filePaths (Alterar depois para retornar na func)
	var filePaths []string

	// Parse yaml
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Error parsing YAML file: %v", err)
	}

	// Iterate over plugins and get the file path
	for _, plugin := range config.Plugins {
		fmt.Printf("Loading plugin %s...\n", plugin.Name)

		if strings.HasPrefix(plugin.Path, "github.com/") {
			downloadRepository(&plugin)
		}
		
		for _, id := range plugin.Ids {
			if id == "*" {
				for _, file := range listFiles(plugin.Path, ".so") {
					name := file[:len(file)-3]
					if contains(plugin.Exclude, name) {
						//Skipping excluded ID
						continue
					}
					filePath := filepath.Join(plugin.Path, file)
					filePaths = append(filePaths, filePath)
				}
				continue
			}

			if contains(plugin.Exclude, id) {
				//Skipping excluded ID
				continue
			}

			filePath := filepath.Join(plugin.Path, id+".so")
			filePaths = append(filePaths, filePath)
		}
	}
}

// Auxiliar function (check if slice contains string)
func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

// Auxiliar function (list files in the directory)
func listFiles(dirPath string, ext string) []string {
	if strings.HasPrefix(dirPath, "~/") {
		usr, err := user.Current()
		if err != nil {
			log.Fatalf("Error getting current user: %v", err)
		}
		dirPath = filepath.Join(usr.HomeDir, dirPath[2:])
	}

	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatalf("Error listing files: %v", err)
	}
	var result []string
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ext {
			continue
		}
		result = append(result, file.Name())
	}
	return result
}

// Auxiliar function (download GitHub repository and update plugin path)
func downloadRepository(plugin *Plugin) {
	// Create ~/.nucke/repositories directory if it doesn't exist
	repoDir := filepath.Join(os.Getenv("HOME"), ".nucke", "repositories")
	err := os.MkdirAll(repoDir, 0755)
	if err != nil {
		log.Fatalf("Error creating repository directory: %v", err)
	}

	// Modify plugin.Path to use SSH pattern
	sshPath := fmt.Sprintf("git@github.com:%s/%s.git", strings.Split(plugin.Path, "/")[1], strings.Split(plugin.Path, "/")[2])
	destinationPath := filepath.Join(repoDir, strings.Split(plugin.Path, "/")[2])

	// Check if repository is already cloned
	if _, err := os.Stat(destinationPath); err == nil {
		// Repository already cloned, update it
		pullCmd := exec.Command("git", "pull")
		pullCmd.Dir = destinationPath
		err = pullCmd.Run()
		if err != nil {
			log.Fatalf("Error updating repository: %v", err)
		}
	} else {
		// Repository not cloned, clone it
		cloneCmd := exec.Command("git", "clone", sshPath, destinationPath)
		err = cloneCmd.Run()
		if err != nil {
			log.Fatalf("Error cloning repository: %v", err)
		}
	}

	// Update plugin.Path to the local directory
	plugin.Path = filepath.Join(repoDir, strings.Join(strings.Split(plugin.Path, "/")[2:], "/"))
}


