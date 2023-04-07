package template

import (
	"io/ioutil"
	"strings"
	"text/template"
	"path/filepath"
	"runtime"
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

func ReadFileToString(report string) string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		panic("Failed to get path of currently executing Go file")
	}
	
	// Get the directory that contains the Go file
	dir := filepath.Dir(filename)

	// Get the name of the currently executing Go file without the extension
	//fileWithoutExt := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	
	// Join the directory, file name, and relative file path
	absPath := filepath.Join(dir, report)

	// Read the file content
	content, err := ioutil.ReadFile(absPath)
	if err != nil {
		panic(err)
	}

	return string(content)
}

