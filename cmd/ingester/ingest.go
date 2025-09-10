package main

import (
	"context"
	"fmt"
	"time"

	loaderpkg "tekton-assist/pkg/ingester"

	llamastackclient "github.com/llamastack/llama-stack-client-go"
	"github.com/llamastack/llama-stack-client-go/option"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(ingestCmd)
}

var ingestCmd = &cobra.Command{
	Use:   "ingest",
	Short: "Ingest knowledge base into vector DB",
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.DryRun {
			return printResult(cfg, map[string]any{
				"status":  "success",
				"action":  "ingest",
				"message": fmt.Sprintf("DRY RUN: would ingest from %s", cfg.KBFile),
			})
		}
		kbPath := resolveKBPath(cfg.KBFile)
		l := loaderpkg.NewLoader()
		entries, err := l.LoadKB(kbPath)
		if err != nil {
			return printResult(cfg, map[string]any{
				"status":  "error",
				"action":  "ingest",
				"message": fmt.Sprintf("Failed to load KB: %v", err),
			})
		}

		client := llamastackclient.NewClient(
			option.WithBaseURL(cfg.EndpointURL),
			option.WithMaxRetries(2),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		// Ensure vector DB exists (best-effort, ignore errors)
		_, _ = client.VectorDBs.Register(ctx, llamastackclient.VectorDBRegisterParams{
			EmbeddingModel:     cfg.EmbeddingModel,
			VectorDBID:         cfg.VectorDBID,
			EmbeddingDimension: llamastackclient.Int(int64(cfg.EmbeddingDim)),
			ProviderID:         llamastackclient.String("faiss"),
		})

		// Build typed documents
		documents := make([]llamastackclient.DocumentParam, 0, len(entries))
		now := float64(time.Now().Unix())
		for _, e := range entries {
			meta := map[string]llamastackclient.DocumentMetadataUnionParam{
				"error":      {OfString: llamastackclient.String(trunc(e.Error, 200))},
				"error_type": {OfString: llamastackclient.String(e.ErrorType)},
				"severity":   {OfString: llamastackclient.String(e.Severity)},
				"source":     {OfString: llamastackclient.String(e.Source)},
				"reference":  {OfString: llamastackclient.String(e.Reference)},
				"task_name":  {OfString: llamastackclient.String(fmt.Sprint(e.Metadata["task_name"]))},
				"repository": {OfString: llamastackclient.String(fmt.Sprint(e.Metadata["repository"]))},
				"timestamp":  {OfFloat: llamastackclient.Float(now)},
			}
			documents = append(documents, llamastackclient.DocumentParam{
				DocumentID: e.ID,
				Content:    llamastackclient.DocumentContentUnionParam{OfString: llamastackclient.String(e.CombinedText)},
				Metadata:   meta,
			})
		}

		// Batch insert via typed RAG tool
		success := 0
		for i := 0; i < len(documents); i += cfg.BatchSize {
			end := i + cfg.BatchSize
			if end > len(documents) {
				end = len(documents)
			}
			batch := documents[i:end]

			if err := client.ToolRuntime.RagTool.Insert(ctx, llamastackclient.ToolRuntimeRagToolInsertParams{
				ChunkSizeInTokens: int64(cfg.ChunkSize),
				Documents:         batch,
				VectorDBID:        cfg.VectorDBID,
			}); err != nil {
				return printResult(cfg, map[string]any{
					"status":  "error",
					"action":  "ingest",
					"message": fmt.Sprintf("Batch %d insert failed: %v", (i/cfg.BatchSize)+1, err),
				})
			}
			success += end - i
		}

		return printResult(cfg, map[string]any{
			"status":  "success",
			"action":  "ingest",
			"message": fmt.Sprintf("Ingested %d entries successfully", success),
		})
	},
}
