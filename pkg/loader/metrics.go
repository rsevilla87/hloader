package loader

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "hloader_http_requests_total",
		Help: "Total number of HTTP requests by status code",
	}, []string{"code"})

	httpTimeoutsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hloader_http_timeouts_total",
		Help: "Total number of HTTP request timeouts",
	})

	httpReadErrorsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hloader_http_read_errors_total",
		Help: "Total number of HTTP response body read errors",
	})

	requestLatency = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "hloader_request_latency_seconds",
		Help: "HTTP request latency in seconds",
		// 100µs .. 10s — aligned with typical sub-ms to low-ms HTTP latencies
		Buckets: []float64{
			0.0001, 0.00025, 0.0005, 0.001, 0.0025, 0.005, 0.01, 0.025,
			0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10,
		},
	})

	bytesReadTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hloader_bytes_read_total",
		Help: "Total number of bytes read from HTTP responses",
	})
)

func startMetricsServer(port int) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Prometheus metrics listening on %s/metrics\n", addr)
	go func() {
		if err := http.ListenAndServe(addr, mux); err != nil {
			log.Printf("metrics server error: %v\n", err)
		}
	}()
}

func (r requestResult) observeMetrics() {
	if r.timeout {
		httpTimeoutsTotal.Inc()
		return
	}
	if r.readError {
		httpReadErrorsTotal.Inc()
	}
	httpRequestsTotal.WithLabelValues(strconv.Itoa(r.code)).Inc()
	requestLatency.Observe(float64(r.latency) / 1e6)
	bytesReadTotal.Add(float64(r.bytesRead))
}
