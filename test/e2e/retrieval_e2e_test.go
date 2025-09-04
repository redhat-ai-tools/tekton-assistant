package e2e

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	llamastackclient "github.com/llamastack/llama-stack-client-go"
	"github.com/llamastack/llama-stack-client-go/option"
	"github.com/llamastack/llama-stack-client-go/shared"
)

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// TestRetrievalE2E verifies that the vector DB exists and a RAG query returns content.
// It assumes the knowledge base has been ingested separately.
func TestRetrievalE2E(t *testing.T) {
	endpoint := getEnvOrDefault("LLAMA_STACK_ENDPOINT", "http://localhost:8321")
	vectorDBID := getEnvOrDefault("VECTOR_DB_ID", "tekton_errors_db")
	modelID := getEnvOrDefault("INFERENCE_MODEL", "gemini-1.5-flash")

	client := llamastackclient.NewClient(
		option.WithBaseURL(endpoint),
		option.WithMaxRetries(2),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 1) Check vector DB exists
	list, err := client.VectorDBs.List(ctx)
	if err != nil {
		t.Fatalf("failed to list vector DBs: %v", err)
	}
	exists := false
	if list != nil {
		for _, db := range *list {
			if db.Identifier == vectorDBID {
				exists = true
				break
			}
		}
	}
	if !exists {
		t.Fatalf("vector DB %q not found at %s (ensure ingestion ran)", vectorDBID, endpoint)
	}

	// 2) Direct RAG query
	query := "create-trusted-artifact"
	res, err := client.ToolRuntime.RagTool.Query(ctx, llamastackclient.ToolRuntimeRagToolQueryParams{
		VectorDBIDs: []string{vectorDBID},
		Content:     llamastackclient.InterleavedContentUnionParam{OfString: llamastackclient.String(query)},
		QueryConfig: llamastackclient.QueryConfigParam{
			ChunkTemplate:      "Result {index}\nContent: {chunk.content}\nMetadata: {metadata}\n",
			MaxChunks:          5,
			MaxTokensInContext: 4096,
		},
	})
	if err != nil {
		t.Fatalf("rag query failed: %v", err)
	}
	raw := res.RawJSON()
	if raw == "" {
		t.Fatalf("empty RAG response JSON")
	}
	if !strings.Contains(strings.ToLower(raw), "create-trusted-artifact") {
		t.Fatalf("rag response did not include expected signal; got: %.200s", raw)
	}

	// 3) Agent-based prompt flow (streaming)
	agentResp, err := client.Agents.New(ctx, llamastackclient.AgentNewParams{
		AgentConfig: shared.AgentConfigParam{
			Instructions:             "You are a Tekton debugging expert. Prefer using the knowledge base (builtin::rag) before answering.",
			Model:                    modelID,
			EnableSessionPersistence: llamastackclient.Bool(true),
			ToolConfig:               shared.AgentConfigToolConfigParam{ToolChoice: "auto"},
			Toolgroups: []shared.AgentConfigToolgroupUnionParam{{
				OfAgentToolGroupWithArgs: &shared.AgentConfigToolgroupAgentToolGroupWithArgsParam{
					Name: "builtin::rag",
					Args: map[string]shared.AgentConfigToolgroupAgentToolGroupWithArgsArgUnionParam{
						"vector_db_ids": {OfAnyArray: []any{vectorDBID}},
					},
				},
			}},
		},
	})
	if err != nil {
		t.Fatalf("failed to create agent: %v", err)
	}

	sess, err := client.Agents.Session.New(ctx, agentResp.AgentID, llamastackclient.AgentSessionNewParams{SessionName: "tekton-e2e"})
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	stream := client.Agents.Turn.NewStreaming(ctx, sess.SessionID, llamastackclient.AgentTurnNewParams{
		AgentID: agentResp.AgentID,
		Messages: []llamastackclient.AgentTurnNewParamsMessageUnion{{
			OfUserMessage: &shared.UserMessageParam{
				Content: shared.InterleavedContentUnionParam{OfString: llamastackclient.String("Use knowledge_search on my query, then answer: Git clone step failed creating trusted artifact; how to fix?")},
			},
		}},
		ToolConfig: llamastackclient.AgentTurnNewParamsToolConfig{ToolChoice: "auto"},
		Toolgroups: []llamastackclient.AgentTurnNewParamsToolgroupUnion{{
			OfAgentToolGroupWithArgs: &llamastackclient.AgentTurnNewParamsToolgroupAgentToolGroupWithArgs{
				Name: "builtin::rag",
				Args: map[string]llamastackclient.AgentTurnNewParamsToolgroupAgentToolGroupWithArgsArgUnion{
					"vector_db_ids": {OfAnyArray: []any{vectorDBID}},
				},
			},
		}},
	})
	if err := stream.Err(); err != nil {
		t.Fatalf("agent streaming turn start failed: %v", err)
	}

	var (
		finalOut string
		ragUsed  bool
	)
	for stream.Next() {
		chunk := stream.Current()
		raw := chunk.Event.RawJSON()
		payload := chunk.Event.Payload.AsAny()
		if strings.Contains(raw, "\"step_type\":\"tool_execution\"") || strings.Contains(raw, "knowledge_search") {
			ragUsed = true
		}
		switch p := payload.(type) {
		case llamastackclient.TurnResponseEventPayloadTurnComplete:
			finalOut = p.Turn.OutputMessage.Content.AsString()
		}
	}
	if err := stream.Err(); err != nil {
		t.Fatalf("agent streaming turn failed: %v", err)
	}
	if finalOut == "" {
		t.Fatalf("empty agent output")
	}
	if !ragUsed {
		t.Fatalf("agent turn did not use RAG (knowledge_search/tool_execution not observed)")
	}
	if !strings.Contains(strings.ToLower(finalOut), "registry") && !strings.Contains(strings.ToLower(finalOut), "artifact") {
		t.Fatalf("agent output missing expected troubleshooting terms: %.200s", finalOut)
	}
}
