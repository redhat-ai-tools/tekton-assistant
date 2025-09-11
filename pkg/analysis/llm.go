package analysis

import (
	"context"
	"fmt"
	"os"
	"time"

	openai "github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

type LLM interface {
	Analyze(ctx context.Context, input string) (string, error)
}

// OpenAIConfig holds configuration for the OpenAI-backed LLM.
type OpenAIConfig struct {
	APIKey         string
	Model          string
	Temperature    float32
	MaxTokens      int
	BaseURL        string
	RequestTimeout time.Duration
}

type openAILLM struct {
	client    openai.Client
	model     string
	temp      float32
	maxTokens int
}

// NewOpenAILLM constructs an LLM that talks to OpenAI's chat completions.
func NewOpenAILLM(cfg OpenAIConfig) (LLM, error) {
	apiKey := cfg.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}

	// Build client options
	opts := []option.RequestOption{option.WithAPIKey(apiKey)}
	if cfg.BaseURL != "" {
		opts = append(opts, option.WithBaseURL(cfg.BaseURL))
	}
	c := openai.NewClient(opts...)
	model := cfg.Model
	if model == "" {
		model = "gpt-4o-mini"
	}
	return &openAILLM{
		client:    c,
		model:     model,
		temp:      cfg.Temperature,
		maxTokens: cfg.MaxTokens,
	}, nil
}

func (o *openAILLM) Analyze(ctx context.Context, input string) (string, error) {
	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a Tekton TaskRun debugging expert. Provide concise, actionable diagnosis and fixes."),
			openai.UserMessage(input),
		},
		Model: openai.ChatModel(o.model),
	}
	// Note: temperature and max tokens omitted for now to avoid param.Opt types

	resp, err := o.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("empty completion choices")
	}
	return resp.Choices[0].Message.Content, nil
}
