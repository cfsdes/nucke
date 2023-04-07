# Sample

A scanner can be created using the sample file as a starter code.

The scanner should return:
- **severity**: Critical, High, Medium, Low or Info
- **url**: Vulnerable endpoint identified
- **summary**: Vulnerability report (it supports markdown)
- **vulnFound**: Boolean value. If true, the scanner will report the vulnerability
- **error**

The `Run()` function is the function that will be called by the Nucke:

```go
func Run(r *http.Request, client *http.Client) (string, string, string, bool, error)
```

After created, the scanner should be compiled as a plugin:
```bash
go build -buildmode=plugin -o scanner.so scanner.go
```
