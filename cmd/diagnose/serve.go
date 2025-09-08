package main

import (
	"log"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the context extractor HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.New(os.Stdout, "context-extractor ", log.LstdFlags|log.Lshortfile)
		srv := NewHTTPServer(cfg.Addr, logger)

		var wg sync.WaitGroup
		srv.startListener(&wg)
		wg.Wait()
	},
}
