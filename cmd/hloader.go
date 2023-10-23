package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/rsevilla87/hloader/pkg/loader"
	"github.com/spf13/cobra"
)

func main() {
	var duration, requestTimeout time.Duration
	var requestRate, connections int
	var url string
	var pprof, http2, insecureSkipVerify, keepalive bool
	rootCmd := &cobra.Command{
		Use:          fmt.Sprintf(os.Args[0]),
		Short:        "Simple http loader",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if requestRate > 0 && requestRate < connections {
				return fmt.Errorf("request rate must be higher than connections")
			}
			if pprof {
				go func() {
					log.Println(http.ListenAndServe("localhost:6060", nil))
				}()
			}
			l := loader.NewLoader(duration, requestRate, connections, url, requestTimeout, insecureSkipVerify, keepalive, http2)
			return l.Run()
		},
	}
	rootCmd.Flags().StringVarP(&url, "url", "u", "", "Target URL")
	rootCmd.Flags().IntVarP(&requestRate, "rate", "r", 0, "Request rate, 0 means unlimited")
	rootCmd.Flags().IntVarP(&connections, "concurrency", "c", 1, "Number of concurrent connections")
	rootCmd.Flags().DurationVarP(&duration, "duration", "d", 10*time.Second, "Test duration")
	rootCmd.Flags().DurationVarP(&requestTimeout, "requestTimeout", "t", 1*time.Second, "Request timeout")
	rootCmd.Flags().BoolVarP(&keepalive, "keepalive", "k", true, "Enable HTTP keepalive")
	rootCmd.Flags().BoolVar(&http2, "http2", true, "Enable HTTP2")
	rootCmd.Flags().BoolVar(&pprof, "pprof", false, "Enable pprof endpoint in localhost:6060")
	rootCmd.Flags().BoolVarP(&insecureSkipVerify, "insecureSkipVerify", "i", true, "Skip server's certificate verification")
	rootCmd.MarkFlagRequired("url")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
