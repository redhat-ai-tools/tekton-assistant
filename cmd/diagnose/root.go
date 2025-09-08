package main

import (
	"github.com/spf13/cobra"
)

type Config struct {
	Addr string
}

var (
	rootCmd = &cobra.Command{Use: "context-extractor", Short: "Tekton TaskRun context extractor"}
	cfg     = &Config{}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&cfg.Addr, "addr", ":8080", "HTTP listen address")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func main() { Execute() }
