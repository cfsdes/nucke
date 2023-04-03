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

<div align="center"><img width="385" alt="Captura de Tela 2023-04-03 às 18 16 17" src="https://user-images.githubusercontent.com/20214824/229629405-7111d334-7957-4873-9d3e-ce46c956e588.png"></div>
