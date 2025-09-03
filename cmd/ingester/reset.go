package main

import (
	"context"
	"fmt"
	"time"

	"github.com/fatih/color"
	llamastackclient "github.com/llamastack/llama-stack-client-go"
	"github.com/llamastack/llama-stack-client-go/option"
	"github.com/spf13/cobra"
)

func init() { rootCmd.AddCommand(resetCmd) }

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset/clear vector database",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !cfg.Force {
			c := color.New(color.FgYellow)
			c.Printf("This will delete vector database '%s'. Re-run with --force to proceed.\n", cfg.VectorDBID)
			return nil
		}
		client := llamastackclient.NewClient(
			option.WithBaseURL(cfg.EndpointURL),
			option.WithMaxRetries(2),
		)
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		if err := client.VectorDBs.Unregister(ctx, cfg.VectorDBID); err != nil {
			return printResult(cfg, map[string]any{
				"status":  "error",
				"action":  "reset",
				"message": fmt.Sprintf("Failed to reset vector database: %v", err),
			})
		}
		return printResult(cfg, map[string]any{
			"status":  "success",
			"action":  "reset",
			"message": fmt.Sprintf("Vector database '%s' has been reset", cfg.VectorDBID),
		})
	},
}
