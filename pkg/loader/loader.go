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

func NewLoader(duration time.Duration, requestRate, connections int, url string, requestTimeout time.Duration, insecureSkipVerify bool, keepalive bool, http2 bool) Loader {
	var limit rate.Limit
	if requestRate > 0 {
		limit = rate.Limit(requestRate + 1) // We add 1 to count the main goroutine
	}
	return Loader{
		url:                url,
		requestRate:        requestRate,
		duration:           duration,
		connections:        connections,
		insecureSkipVerify: insecureSkipVerify,
		requestTimeout:     requestTimeout,
		keepalive:          keepalive,
		limiter:            rate.NewLimiter(limit, 1),
		http2:              http2,
	}
}

func (l *Loader) Run() error {
	var wg sync.WaitGroup
	signalCh := make(chan os.Signal, 1)
	stopCh := make(chan struct{})
	now := time.Now()
	for i := 0; i < l.connections; i++ {
		wg.Add(1)
		go l.load(&wg, stopCh)
	}
	signal.Notify(signalCh, os.Interrupt)
out:
	for {
		select {
		case <-signalCh:
			fmt.Println("Interrupt signal received, exiting gracefully")
			break out
		case <-time.After(l.duration):
			break out
		}
	}
	close(stopCh)
	wg.Wait()
	l.duration = time.Since(now)
	err := normaliceResults(l.results, l.duration)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return nil
}

func (l *Loader) load(wg *sync.WaitGroup, stopCh chan struct{}) {
	defer wg.Done()
	httpClient := &http.Client{
		Timeout: l.requestTimeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: l.insecureSkipVerify,
			},
			DisableKeepAlives: !l.keepalive, //false by default
			ForceAttemptHTTP2: l.http2,
		},
	}
	req, err := http.NewRequest(http.MethodGet, l.url, http.NoBody)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for {
		l.limiter.Wait(context.TODO())
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
	defer func() { // Ensure to add the result
		l.Lock()
		l.results = append(l.results, result)
		l.Unlock()
	}()
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
	bytesRead, _ := io.ReadAll(resp.Body)
	result.bytesRead = int64(len(bytesRead))
	result.code = resp.StatusCode
	result.latency = time.Since(result.timestamp).Microseconds()
}
