package main

import (
	"github.com/spf13/cobra"
)

// Config mirrors Python flags
type Config struct {
	KBFile         string
	EndpointURL    string
	VectorDBID     string
	EmbeddingModel string
	EmbeddingDim   int
	ChunkSize      int
	BatchSize      int
	LogLevel       string
	MaxResults     int
	MinScore       float64
	Quiet          bool
	JSONOutput     bool
	OutputFile     string
	Force          bool
	DryRun         bool
}

var (
	rootCmd = &cobra.Command{Use: "ingester", Short: "Tekton Knowledge Base Ingester (Go)"}
	cfg     = &Config{}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&cfg.KBFile, "kb-file", "data/curated/kb.json", "Knowledge base JSON file")
	rootCmd.PersistentFlags().StringVar(&cfg.EndpointURL, "endpoint", "http://localhost:8321", "LlamaStack endpoint URL")
	rootCmd.PersistentFlags().StringVar(&cfg.VectorDBID, "vector-db", "tekton_errors_db", "Vector database identifier")
	rootCmd.PersistentFlags().StringVar(&cfg.EmbeddingModel, "embedding-model", "text-embedding-004", "Embedding model name")
	rootCmd.PersistentFlags().IntVar(&cfg.EmbeddingDim, "embedding-dim", 384, "Embedding dimension")
	rootCmd.PersistentFlags().IntVar(&cfg.ChunkSize, "chunk-size", 512, "Chunk size in tokens")
	rootCmd.PersistentFlags().IntVar(&cfg.BatchSize, "batch-size", 50, "Batch size for ingestion")
	rootCmd.PersistentFlags().StringVar(&cfg.LogLevel, "log-level", "INFO", "Logging level")
	rootCmd.PersistentFlags().BoolVar(&cfg.Quiet, "quiet", false, "Suppress colored output")
	rootCmd.PersistentFlags().BoolVar(&cfg.JSONOutput, "json-output", false, "Output results in JSON format")
	rootCmd.PersistentFlags().StringVar(&cfg.OutputFile, "output-file", "", "Save results to file")
	rootCmd.PersistentFlags().BoolVar(&cfg.Force, "force", false, "Force operations without confirmation")
	rootCmd.PersistentFlags().BoolVar(&cfg.DryRun, "dry-run", false, "Show actions without executing")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func main() { Execute() }
