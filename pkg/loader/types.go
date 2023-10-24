package loader

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const GET = "GET"

type Loader struct {
	url                string
	requestRate        int
	duration           time.Duration
	connections        int
	insecureSkipVerify bool
	requestTimeout     time.Duration
	keepalive          bool
	results            []requestResult
	limiter            *rate.Limiter
	http2              bool
	sync.Mutex
}

type requestResult struct {
	latency   int64
	timestamp time.Time
	code      int
	timeout   bool
	readError bool
	bytesRead int64
}

type testResult struct {
	RPS           float64       `json:"rps"`
	Timeouts      int64         `json:"timeouts"`
	ReadErrors    int64         `json:"read_errors"`
	AvgThroughput int64         `json:"avg_throughput_bps"`
	AvgLatency    float64       `json:"avg_lat_us"`
	MaxLatency    float64       `json:"max_lat_us"`
	P99Latency    float64       `json:"p99_lat_us"`
	P95Latency    float64       `json:"p95_lat_us"`
	P90Latency    float64       `json:"p90_lat_us"`
	LatencyStdev  float64       `json:"stdev_lat"`
	StatusCodes   map[int]int64 `json:"status_codes"`
}
