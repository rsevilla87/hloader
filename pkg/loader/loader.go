package loader

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func NewLoader(duration, requestTimeout time.Duration, requestRate, connections int, url string, insecure, keepalive, http2 bool, csv string, metrics bool, port int) Loader {
	var limit rate.Limit
	if requestRate > 0 {
		limit = rate.Limit(requestRate + 1) // We add 1 to count the main goroutine
	}
	return Loader{
		url:                url,
		requestRate:        requestRate,
		duration:           duration,
		connections:        connections,
		insecureSkipVerify: insecure,
		requestTimeout:     requestTimeout,
		keepalive:          keepalive,
		limiter:            rate.NewLimiter(limit, 1),
		http2:              http2,
		csv:                csv,
		metrics:            metrics,
		port:               port,
	}
}

func (l *Loader) Run() error {
	var wg sync.WaitGroup
	signalCh := make(chan os.Signal, 1)
	stopCh := make(chan struct{})
	now := time.Now()
	if l.metrics {
		startMetricsServer(l.port)
	}
	httpClient := &http.Client{
		Timeout: l.requestTimeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: l.insecureSkipVerify,
			},
			DisableKeepAlives:   !l.keepalive,
			ForceAttemptHTTP2:   l.http2,
			MaxIdleConns:        l.connections,
			MaxConnsPerHost:     l.connections,
			MaxIdleConnsPerHost: l.connections,
		},
	}
	for i := 0; i < l.connections; i++ {
		req, err := http.NewRequest(http.MethodGet, l.url, http.NoBody)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		wg.Go(func() {
			l.load(httpClient, req, stopCh)
		})
	}
	var durationCh <-chan time.Time
	if l.duration > 0 {
		durationCh = time.After(l.duration)
	}
	signal.Notify(signalCh, os.Interrupt)
out:
	for {
		select {
		case <-signalCh:
			os.Stderr.Write([]byte("Interrupt signal received, exiting gracefully\n"))
			break out
		case <-durationCh:
			break out
		}
	}
	close(stopCh)
	wg.Wait()
	// If metrics are enabled, we don't need to normalize the results as no output will be generated
	if l.metrics {
		return nil
	}
	l.duration = time.Since(now)
	err := normaliceResults(l.results, l.duration, l.csv)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return nil
}

func (l *Loader) load(httpClient *http.Client, req *http.Request, stopCh chan struct{}) {
	for {
		select {
		case <-stopCh:
			httpClient.CloseIdleConnections()
			return
		default:
			l.sendRequest(httpClient, req)
		}
	}
}

func (l *Loader) sendRequest(httpClient *http.Client, req *http.Request) {
	result := requestResult{}
	defer func() {
		if l.metrics {
			result.observeMetrics()
			return
		} else {
			l.Lock()
			l.results = append(l.results, result)
			l.Unlock()
		}
	}()
	l.limiter.Wait(context.TODO())
	result.timestamp = time.Now()
	resp, err := httpClient.Do(req)
	if err != nil {
		if os.IsTimeout(err) {
			result.timeout = true
		} else {
			fmt.Println(err)
			os.Exit(1)
		}
		return
	}
	defer resp.Body.Close()
	result.bytesRead, err = io.Copy(io.Discard, resp.Body)
	result.latency = time.Since(result.timestamp).Microseconds()
	if err != nil {
		result.readError = true
	}
	result.code = resp.StatusCode
}
