# Template Package

## Parser

Import:
```go
import "github.com/cfsdes/nucke/pkg/template"
```

Parse template from string:
```go
result := template.ParseTemplateFromFile("Hello {{.msg}}", map[string]interface{}{
    "msg": "World",
})
```

Parse template from file:
```go
templateString, err := template.ReadFileToString("template-report.txt")
summary := template.ParseTemplate(templateString, map[string]interface{}{
    "msg": "World",
})
```