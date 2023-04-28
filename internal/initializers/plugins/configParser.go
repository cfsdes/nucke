package plugins

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
	"github.com/fatih/color"
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

func ParseConfig(configFile string) (filePaths []string){
	// Initial message
	Cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	fmt.Printf("[%s] Loading plugins...\n", Cyan("INF"))

	// Read config file
	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	// Parse yaml
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Error parsing YAML file: %v", err)
	}

	// Iterate over plugins and get the file path
	for _, plugin := range config.Plugins {
		if strings.HasPrefix(plugin.Path, "github.com/") {
			downloadRepository(&plugin)
		}
		
		for _, id := range plugin.Ids {
			if id == "*" {
				for _, path := range listFiles(plugin.Path, ".go") {
					name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
					if contains(plugin.Exclude, name) {
						//Skipping excluded ID
						continue
					}

					filePaths = append(filePaths, path)
				}
				continue
			}

			if contains(plugin.Exclude, id) {
				//Skipping excluded ID
				continue
			}

			if id != "*" {
				for _, path := range listFiles(plugin.Path, ".go") {
					name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
					if name == id {
						filePaths = append(filePaths, path)
						continue
					}
				}
				continue
			}
		}
	}

	return
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

	var result []string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing file %q: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ext {
			result = append(result, path)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error listing files: %v", err)
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


