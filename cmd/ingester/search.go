package main

import (
	"context"
	"fmt"
	"time"

	llamastackclient "github.com/llamastack/llama-stack-client-go"
	"github.com/llamastack/llama-stack-client-go/option"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(searchCmd)
}

var (
	query      string
	maxResults int
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search the vector DB",
	RunE: func(cmd *cobra.Command, args []string) error {
		if query == "" {
			return fmt.Errorf("--query is required")
		}
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		client := llamastackclient.NewClient(
			option.WithBaseURL(cfg.EndpointURL),
			option.WithMaxRetries(2),
		)

		res, err := client.ToolRuntime.RagTool.Query(ctx, llamastackclient.ToolRuntimeRagToolQueryParams{
			VectorDBIDs: []string{cfg.VectorDBID},
			Content:     llamastackclient.InterleavedContentUnionParam{OfString: llamastackclient.String(query)},
			QueryConfig: llamastackclient.QueryConfigParam{
				ChunkTemplate:      "Result {index}\nContent: {chunk.content}\nMetadata: {metadata}\n",
				MaxChunks:          int64(maxResults),
				MaxTokensInContext: 4096,
			},
		})
		if err != nil {
			return printResult(cfg, map[string]any{
				"status":  "error",
				"action":  "search",
				"message": fmt.Sprintf("Query failed: %v", err),
			})
		}

		// The response content is interleaved; for now, just return raw JSON
		return printResult(cfg, map[string]any{
			"status":  "success",
			"action":  "search",
			"results": map[string]any{"raw": res.RawJSON()},
			"message": fmt.Sprintf("Search succeeded for query: '%s'", query),
		})
	},
}

func init() {
	searchCmd.Flags().StringVar(&query, "query", "", "Search query")
	searchCmd.Flags().IntVar(&maxResults, "max-results", 5, "Maximum search results")
}
