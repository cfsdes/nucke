package template

import (
	"io/ioutil"
	"strings"
	"text/template"
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
    content, err := ioutil.ReadFile(pluginDir + "/" + report)
    if err != nil {
        panic("Error while reading template report. Plugin directory not found")
    }

    return string(content)
}


