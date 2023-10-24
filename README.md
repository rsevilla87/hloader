# hloader

```shell
 $ ./bin/hloader -h
HTTP loader

Usage:
  ./bin/hloader [flags]
  ./bin/hloader [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Get version info

Flags:
  -u, --url string          Target URL
  -r, --rate int            Request rate, 0 means unlimited
  -c, --concurrency int     Number of concurrent connections (default 1)
  -d, --duration duration   Test duration (default 10s)
  -t, --timeout duration    Request timeout (default 1s)
  -i, --insecure            Skip server's certificate verification (default true)
  -k, --keepalive           Enable HTTP keepalive (default true)
      --http2               Use HTTP2 protocol, if possible (default true)
      --pprof               Enable pprof endpoint in localhost:6060
  -h, --help                help for ./bin/hloader
```

## Compilation

```
$ make build
GOARCH=amd64 CGO_ENABLED=0 go build -v -ldflags "-X github.com/cloud-bulldozer/go-commons/version.GitCommit=a5c03b3c983255096635b872d4153c98419f8bd1 -X github.com/cloud-bulldozer/go-commons/version.Version=main -X github.com/cloud-bulldozer/go-commons/version.BuildDate=2023-10-24-12:13:18" -o bin/hloader cmd/hloader.go
```
