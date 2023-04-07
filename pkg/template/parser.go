package template

import (
	"io/ioutil"
	"strings"
	"text/template"
)

func ParseTemplateFromFile(templateFile string, data interface{}) (string) {
	templateString, err := readFileToString(templateFile)
    if err != nil {
        panic(err)
    }

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

func ParseTemplateFromString(templateString string, data interface{}) (string) {
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

func readFileToString(filepath string) (string, error) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

