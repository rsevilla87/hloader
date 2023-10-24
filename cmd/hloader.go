package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/cloud-bulldozer/go-commons/version"
	"github.com/rsevilla87/hloader/pkg/loader"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get version info",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version:", version.Version)
		fmt.Println("Git Commit:", version.GitCommit)
		fmt.Println("Build Date:", version.BuildDate)
		fmt.Println("Go Version:", version.GoVersion)
		fmt.Println("OS/Arch:", version.OsArch)
	},
}

func main() {
	var duration, requestTimeout time.Duration
	var requestRate, connections int
	var url string
	var pprof, http2, insecure, keepalive bool
	rootCmd := &cobra.Command{
		Use:   fmt.Sprintf(os.Args[0]),
		Short: "Simple http loader",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			if requestRate > 0 && requestRate < connections {
				return fmt.Errorf("request rate must be higher than connections")
			}
			if pprof {
				go func() {
					log.Println(http.ListenAndServe("localhost:6060", nil))
				}()
			}
			l := loader.NewLoader(duration, requestTimeout, requestRate, connections, url, insecure, keepalive, http2)
			return l.Run()
		},
	}
	rootCmd.Flags().StringVarP(&url, "url", "u", "", "Target URL")
	rootCmd.Flags().IntVarP(&requestRate, "rate", "r", 0, "Request rate, 0 means unlimited")
	rootCmd.Flags().IntVarP(&connections, "concurrency", "c", 1, "Number of concurrent connections")
	rootCmd.Flags().DurationVarP(&duration, "duration", "d", 10*time.Second, "Test duration")
	rootCmd.Flags().DurationVarP(&requestTimeout, "timeout", "t", 1*time.Second, "Request timeout")
	rootCmd.Flags().BoolVarP(&insecure, "insecure", "i", true, "Skip server's certificate verification")
	rootCmd.Flags().BoolVarP(&keepalive, "keepalive", "k", true, "Enable HTTP keepalive")
	rootCmd.Flags().BoolVar(&http2, "http2", true, "Use HTTP2 protocol, if possible")
	rootCmd.Flags().BoolVar(&pprof, "pprof", false, "Enable pprof endpoint in localhost:6060")
	rootCmd.Flags().SortFlags = false
	rootCmd.MarkFlagRequired("url")
	rootCmd.AddCommand(versionCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
