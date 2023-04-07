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
result := template.ParseTemplateFromFile("template.txt", map[string]interface{}{
    "description": description,
})
```