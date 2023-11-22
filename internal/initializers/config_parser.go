package initializers

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
	"github.com/cfsdes/nucke/pkg/globals"
)

type Plugin struct {
	Name    string   `yaml:"name"`
	Path    string   `yaml:"path"`
	Webhook string   `yaml:"webhook"`
	Ids     []string `yaml:"ids"`
	Exclude []string `yaml:"exclude,omitempty"`
}

type Config struct {
	Scope string `yaml:"scope"`
	Plugins []Plugin `yaml:"plugins"`
}

var filePaths []string

func ParseConfig(configFile string) (scope string, pluginPaths []string){
	// Initial message
	if globals.Debug {
		Blue := color.New(color.FgBlue, color.Bold).SprintFunc()
		fmt.Printf("[%s] Parsing config file...\n", Blue("DEBUG"))
	}

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

	// Get Scope from Config.yaml
	scope = config.Scope

	// Initialize map array
	globals.Webhook = make(map[string][]string)

	// Iterate over plugins and get the file path
	for _, plugin := range config.Plugins {

		// Atualizando o path para o diretorio correto
		plugin = formatPluginPath(plugin)
		
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

		// List Plugins
		if globals.ListPlugins {
			ListPlugins(plugin.Path)
		}

		// Check Plugins
		CheckLoadedPlugins(filePaths, plugin.Ids)

		// Build Plugins
		pluginPaths = BuildPlugins(filePaths)

		// Criar array para webhook
		globals.Webhook[plugin.Webhook] = pluginPaths
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
func downloadRepository(plugin Plugin) (Plugin) {
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

	return plugin
}


func formatPluginPath(plugin Plugin) (Plugin) {

	// Atualizando o path para o diretorio correto
	if strings.HasPrefix(plugin.Path, "github.com/") {
		return downloadRepository(plugin)
	}

	// Ajustar ~ para HOME
	if strings.HasPrefix(plugin.Path, "~/") {
		usr, err := user.Current()
		if err != nil {
			log.Fatalf("Error getting current user: %v", err)
		}
		plugin.Path = filepath.Join(usr.HomeDir, plugin.Path[2:])
	}

	// Obter o caminho absoluto completo
	if !filepath.IsAbs(plugin.Path) {
		absDir, err := filepath.Abs(plugin.Path)
		if err != nil && globals.Debug {
			fmt.Printf("Error getting absolute plugin path: %v\n", err)
		}
		plugin.Path = absDir
	}

	return plugin
}