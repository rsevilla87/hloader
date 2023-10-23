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
	bytesRead int64
}

type testResult struct {
	RPS           int64         `json:"rps"`
	Timeouts      int64         `json:"timeouts"`
	AvgLatency    float64       `json:"avgLatency"`
	AvgThroughput int64         `json:"avgThroughput"`
	P99Latency    float64       `json:"p99Latency"`
	P95Latency    float64       `json:"p95Latency"`
	P90Latency    float64       `json:"p90Latency"`
	P50Latency    float64       `json:"p50Latency"`
	LatencyStdev  float64       `json:"latencyStdev"`
	MaxLatency    float64       `json:"maxLatency"`
	ResponseCodes map[int]int64 `json:"responseCodes"`
}
