package template

import (
	"io/ioutil"
	"strings"
	"text/template"
	"path/filepath"
	"fmt"

	"github.com/cfsdes/nucke/internal/utils"
)

func ParseTemplate(templateString string, data interface{}) (string) {
	tmpl, err := template.New("file").Parse(templateString)
	if err != nil {
		panic(err)
	}
	var output strings.Builder
	err = tmpl.Execute(&output, data)
	if err != nil {
		panic(err)
	}
	return output.String()
}

func ReadFileToString(report string, pluginName string) string {
	var reportFilePath string

	for _, filePath := range utils.FilePaths {
        fileName := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
        if fileName == pluginName {
            reportFilePath = filepath.Dir(filePath) + "/" + report
        }
    }

	if (reportFilePath == ""){
		fmt.Println("Error while reading template report. Plugin directory not found")
		return ""
	}

    content, err := ioutil.ReadFile(reportFilePath)
    if err != nil {
        panic(err)
    }

    return string(content)
}


