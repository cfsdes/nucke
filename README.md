# Nucke

Nucke is a simple proxy that forwards requests to jaeles API Server, enabling scan of POST requests.

## Install

```
go install github.com/cfsdes/nucke@latest
```

## Usage

```
nucke -port 8080 -jc-api http://127.0.0.1:5000/api/parse
```


