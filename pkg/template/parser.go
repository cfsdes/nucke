package template

import (
	"io/ioutil"
	"strings"
	"text/template"
	"path/filepath"
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

func ReadFileToString(report string, pluginDir string) string {
    absPath := filepath.Join(pluginDir, report)

    content, err := ioutil.ReadFile(absPath)
    if err != nil {
        panic(err)
    }

    return string(content)
}


