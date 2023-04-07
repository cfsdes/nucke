package parsers

import (
	"fmt"
	"github.com/fatih/color"
)

func VulnerabilityOutput(scanName string, severity string, url string, rawReq string, desc string) {
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	//blue := color.New(color.FgBlue).SprintFunc()
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
	}

	// TODO: Salvar rawReq em algum lugar com a descrição. Montar um template de report em markdown

}