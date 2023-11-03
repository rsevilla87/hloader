package loader

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/montanaflynn/stats"
)

func normaliceResults(results []requestResult, duration time.Duration, csvFile string) error {
	var latencies []float64
	var totalRead int64
	var csvWriter *csv.Writer
	var err error
	var f *os.File
	if csvFile != "" {
		f, err = os.Create(csvFile)
		if err != nil {
			return err
		}
		csvWriter = csv.NewWriter(f)
	}
	result := testResult{
		StatusCodes: make(map[int]int64),
	}
	for _, r := range results {
		if csvFile != "" {
			line := []string{
				strconv.FormatInt(r.timestamp.UnixMilli(), 10),
				strconv.Itoa(r.code),
				strconv.FormatInt(r.latency, 10),
				strconv.FormatInt(r.bytesRead, 10),
				strconv.FormatBool(r.timeout),
				strconv.FormatBool(r.readError),
			}
			csvWriter.Write(line)
			csvWriter.Flush()
		}
		result.StatusCodes[r.code]++
		if r.timeout {
			result.Timeouts++
		} else if r.readError {
			result.ReadErrors++
		} else { // timeouts are not valid results to be used in latency calculations
			latencies = append(latencies, float64(r.latency))
			totalRead += r.bytesRead
		}
	}
	result.RPS = math.Floor(float64(result.StatusCodes[http.StatusOK])/duration.Seconds()*100) / 100
	result.AvgThroughput = totalRead / int64(duration.Seconds())
	result.AvgLatency, _ = stats.Mean(latencies)
	result.AvgLatency, _ = stats.Round(result.AvgLatency, 2)
	result.MaxLatency, _ = stats.Max(latencies)
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
