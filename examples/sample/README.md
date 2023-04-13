# Creating plugins

A plugin can be created using the sample file as a starter code.

The plugin should return:
- **severity**: Critical, High, Medium, Low or Info
- **url**: Vulnerable endpoint identified
- **summary**: Vulnerability report (it supports markdown)
- **vulnFound**: Boolean value. If true, the plugin will report the vulnerability
- **error**

The `Run()` function is the function that will be called by the Nucke:

```go
func Run(r *http.Request, w http.ResponseWriter, client *http.Client, pluginDir string) (string, string, string, bool, error)
```

After created, the plugin should be compiled using the following command:
```bash
go build -buildmode=plugin -o plugin.so plugin.go
```

## Report & Directory Structure

The report should be placed in the same directory of plugin.

Example of structure:
```
~/Desktop/nucke-plugins/
    .. sample/
        .. sample.so
        .. report-template.txt
    
    .. sqli/
        .. sqli.so
        .. report-template.txt
```

Configuration file:
```yaml
plugins:
  - name: Example
    path: ~/Desktop/nucke-plugins/
    ids:
      - sample
      - sqli
```
