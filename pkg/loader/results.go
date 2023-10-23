package loader

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/montanaflynn/stats"
)

func normaliceResults(results []requestResult, duration time.Duration) error {
	var latencies []float64
	var totalRead int64
	result := testResult{
		ResponseCodes: make(map[int]int64),
	}
	for _, r := range results {
		result.ResponseCodes[r.code]++
		if r.timeout {
			result.Timeouts++
		} else if r.readError {
			result.ReadErrors++
		} else { // timeouts are not valid results to be used in latency calculations
			latencies = append(latencies, float64(r.latency))
			totalRead += r.bytesRead
		}
	}
	result.RPS = int64(float64(result.ResponseCodes[http.StatusOK]) / float64(duration.Seconds()))
	result.AvgThroughput = totalRead / int64(duration.Seconds())
	result.AvgLatency, _ = stats.Mean(latencies)
	result.AvgLatency, _ = stats.Round(result.AvgLatency, 2)
	result.MaxLatency, _ = stats.Max(latencies)
	result.P50Latency, _ = stats.Percentile(latencies, 50)
	result.P90Latency, _ = stats.Percentile(latencies, 90)
	result.P95Latency, _ = stats.Percentile(latencies, 95)
	result.P99Latency, _ = stats.Percentile(latencies, 99)
	result.LatencyStdev, _ = stats.StandardDeviation(latencies)
	result.LatencyStdev, _ = stats.Round(result.LatencyStdev, 2)
	jsonResult, err := json.MarshalIndent(&result, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(jsonResult))
	return nil
}
