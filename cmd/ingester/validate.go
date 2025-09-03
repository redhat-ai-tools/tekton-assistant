package main

import (
	"context"
	"time"

	llamastackclient "github.com/llamastack/llama-stack-client-go"
	"github.com/llamastack/llama-stack-client-go/option"
	"github.com/spf13/cobra"
)

func init() { rootCmd.AddCommand(validateCmd) }

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate KB and DB connectivity",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := llamastackclient.NewClient(
			option.WithBaseURL(cfg.EndpointURL),
			option.WithMaxRetries(2),
		)
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		kbPath := resolveKBPath(cfg.KBFile)
		kbExists := fileExists(kbPath)
		list, err := client.VectorDBs.List(ctx)
		endpointReachable := err == nil
		vectorExists := false
		if list != nil {
			for _, db := range *list {
				if db.Identifier == cfg.VectorDBID {
					vectorExists = true
					break
				}
			}
		}

		return printResult(cfg, map[string]any{
			"status": "success",
			"action": "validate",
			"validation": map[string]any{
				"kb_file_exists":       kbExists,
				"kb_path":              kbPath,
				"vector_db_exists":     vectorExists,
				"vector_db_accessible": endpointReachable,
				"endpoint_reachable":   endpointReachable,
			},
			"message": "Validation completed",
		})
	},
}
