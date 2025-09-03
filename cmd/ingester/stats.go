package main

import (
	"context"
	"time"

	llamastackclient "github.com/llamastack/llama-stack-client-go"
	"github.com/llamastack/llama-stack-client-go/option"
	"github.com/spf13/cobra"
)

func init() { rootCmd.AddCommand(statsCmd) }

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show DB stats",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := llamastackclient.NewClient(
			option.WithBaseURL(cfg.EndpointURL),
			option.WithMaxRetries(2),
		)
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		list, _ := client.VectorDBs.List(ctx)
		exists := false
		if list != nil {
			for _, db := range *list {
				if db.Identifier == cfg.VectorDBID {
					exists = true
					break
				}
			}
		}

		return printResult(cfg, map[string]any{
			"status": "success",
			"action": "stats",
			"stats": map[string]any{
				"vector_db_id":     cfg.VectorDBID,
				"endpoint_url":     cfg.EndpointURL,
				"embedding_model":  cfg.EmbeddingModel,
				"embedding_dim":    cfg.EmbeddingDim,
				"chunk_size":       cfg.ChunkSize,
				"batch_size":       cfg.BatchSize,
				"vector_db_exists": exists,
			},
			"message": "Statistics retrieved successfully",
		})
	},
}
