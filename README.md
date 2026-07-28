# hloader

[![Go-build](https://github.com/rsevilla87/hloader/actions/workflows/go-build.yml/badge.svg?branch=main&event=push)](https://github.com/rsevilla87/hloader/actions/workflows/go-build.yml)

hloader is a lightweight HTTP/HTTPS load testing tool designed for generating concurrent HTTP requests with configurable rate limiting and connection management.

## Scope and Limitations

hloader is specifically designed for **HTTP/HTTPS load testing only**. It supports:
- HTTP/1.1 and HTTP/2 protocols
- Configurable request rates and concurrency
- TLS/SSL connections with optional certificate verification
- HTTP keepalive connections
- Optional Prometheus metrics endpoint for continuous load generation

**Out of Scope**: hloader does not support other transport protocols such as TCP, UDP, or SCTP. For protocol-level load testing including SCTP support for telco use-cases, please refer to [k8s-netperf](https://github.com/cloud-bulldozer/k8s-netperf).

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
  -u, --url string                   Target URL
  -r, --rate int                     Request rate, 0 means unlimited
  -c, --concurrency int              Number of concurrent connections (default 1)
  -d, --duration duration            Test duration (default 10s)
  -t, --timeout duration             Request timeout (default 1s)
  -i, --insecure                     Skip server's certificate verification (default true)
  -k, --keepalive                    Enable HTTP keepalive (default true)
      --http2                        Use HTTP2 protocol, if possible (default true)
      --pprof                        Enable pprof endpoint in localhost:6060
  -o, --output string                Dump request outputs in the given CSV file
  -m, --enable-prometheus-metrics    Enable Prometheus metrics; run until interrupted (ignores --duration)
  -p, --metrics-port int             Metrics port (default 9000)
  -h, --help                         help for ./bin/hloader

```

## Compilation

```shell
$ make build
GOARCH=amd64 CGO_ENABLED=0 go build -v -ldflags "-X github.com/cloud-bulldozer/go-commons/version.GitCommit=a5c03b3c983255096635b872d4153c98419f8bd1 -X github.com/cloud-bulldozer/go-commons/version.Version=main -X github.com/cloud-bulldozer/go-commons/version.BuildDate=2023-10-24-12:13:18" -o bin/hloader cmd/hloader.go
```

## CSV output

The csv output has the following format:

```csv
1699016948650,200,276028,537,false,false
<timestamp in milliseconds>,<status code>,<latency in microseconds>,<bytes read>,<timeout>,<http error>
```

Like for example:

```csv
1699016948640,200,283085,537,false,false
1699016948650,200,276028,537,false,false
1699016948670,200,255849,537,false,false
1699016948700,200,225850,537,false,false
1699016948709,200,216579,537,false,false
```

## Prometheus metrics

When `--enable-prometheus-metrics` (`-m`) is set, hloader runs until interrupted (`--duration` is ignored) and exposes a Prometheus scrape endpoint at `http://localhost:<metrics-port>/metrics` (default port `9000`). Results are streamed as metrics instead of being aggregated into the JSON summary or CSV output. Stop the load with `SIGINT` (Ctrl+C).

Example:

```shell
$ ./bin/hloader -u https://example.com -c 10 -r 100 -m
$ curl localhost:9000/metrics
```

Exposed metrics:

| Metric | Type | Description |
| --- | --- | --- |
| `hloader_http_requests_total` | Counter (label: `code`) | Total HTTP requests by status code |
| `hloader_http_timeouts_total` | Counter | Total HTTP request timeouts |
| `hloader_http_read_errors_total` | Counter | Total response body read errors |
| `hloader_request_latency_seconds` | Histogram | Request latency in seconds |
| `hloader_bytes_read_total` | Counter | Total bytes read from responses |
