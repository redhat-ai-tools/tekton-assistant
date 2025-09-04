package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	llamastackclient "github.com/llamastack/llama-stack-client-go"
	"github.com/llamastack/llama-stack-client-go/option"
	"github.com/llamastack/llama-stack-client-go/shared"
)

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func main() {
	endpoint := flag.String("endpoint", getenv("LLAMA_STACK_ENDPOINT", "http://localhost:8321"), "LlamaStack endpoint URL")
	vectorDBID := flag.String("vector-db", getenv("VECTOR_DB_ID", "tekton_errors_db"), "Vector database identifier")
	model := flag.String("model", getenv("INFERENCE_MODEL", "gemini-1.5-flash"), "LLM model id")
	query := flag.String("query", "", "Prompt to ask (if empty, reads one line from stdin)")
	debug := flag.Bool("debug", false, "Print streaming events and RAG usage summary")
	flag.Parse()

	if *query == "" {
		fmt.Print("> ")
		s := bufio.NewScanner(os.Stdin)
		if s.Scan() {
			q := strings.TrimSpace(s.Text())
			query = &q
		}
		if *query == "" {
			fmt.Fprintln(os.Stderr, "no query provided")
			os.Exit(1)
		}
	}

	client := llamastackclient.NewClient(
		option.WithBaseURL(*endpoint),
		option.WithMaxRetries(2),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// Optional: warn if vector DB is missing (demo clarity)
	if dbs, err := client.VectorDBs.List(ctx); err == nil {
		exists := false
		if dbs != nil {
			for _, db := range *dbs {
				if db.Identifier == *vectorDBID {
					exists = true
					break
				}
			}
		}
		if !exists {
			fmt.Fprintf(os.Stderr, "warning: vector DB '%s' not found at %s\n", *vectorDBID, *endpoint)
		}
	}

	// Create RAG-enabled agent (builtin::rag tool) with auto tool choice
	agentResp, err := client.Agents.New(ctx, llamastackclient.AgentNewParams{
		AgentConfig: shared.AgentConfigParam{
			Instructions:             "You are a Tekton debugging expert. Prefer using the knowledge base (builtin::rag) before answering. Provide concise, actionable fixes and references.",
			Model:                    *model,
			EnableSessionPersistence: llamastackclient.Bool(true),
			ToolConfig: shared.AgentConfigToolConfigParam{
				ToolChoice: "auto",
			},
			Toolgroups: []shared.AgentConfigToolgroupUnionParam{{
				OfAgentToolGroupWithArgs: &shared.AgentConfigToolgroupAgentToolGroupWithArgsParam{
					Name: "builtin::rag",
					Args: map[string]shared.AgentConfigToolgroupAgentToolGroupWithArgsArgUnionParam{
						"vector_db_ids": {OfAnyArray: []any{*vectorDBID}},
					},
				},
			}},
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create agent: %v\n", err)
		os.Exit(1)
	}

	sess, err := client.Agents.Session.New(ctx, agentResp.AgentID, llamastackclient.AgentSessionNewParams{SessionName: "tekton-demo"})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create session: %v\n", err)
		os.Exit(1)
	}

	// First attempt: auto tool choice, suggest using KB then answer
	userMsg := "Task: Use the knowledge base via knowledge_search on my query, then produce a final answer in this turn. Do NOT ask me for more context. Provide concise, actionable fixes. Query: " + *query
	stream := client.Agents.Turn.NewStreaming(ctx, sess.SessionID, llamastackclient.AgentTurnNewParams{
		AgentID: agentResp.AgentID,
		Messages: []llamastackclient.AgentTurnNewParamsMessageUnion{{
			OfUserMessage: &shared.UserMessageParam{
				Content: shared.InterleavedContentUnionParam{OfString: llamastackclient.String(userMsg)},
			},
		}},
		ToolConfig: llamastackclient.AgentTurnNewParamsToolConfig{ToolChoice: "auto"},
		Toolgroups: []llamastackclient.AgentTurnNewParamsToolgroupUnion{
			{OfAgentToolGroupWithArgs: &llamastackclient.AgentTurnNewParamsToolgroupAgentToolGroupWithArgs{
				Name: "builtin::rag",
				Args: map[string]llamastackclient.AgentTurnNewParamsToolgroupAgentToolGroupWithArgsArgUnion{
					"vector_db_ids": {OfAnyArray: []any{*vectorDBID}},
				},
			}},
		},
	})
	if err := stream.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to start agent turn: %v\n", err)
		os.Exit(1)
	}

	var (
		finalOut   string
		ragUsed    bool
		ragSnippet string
	)
	for stream.Next() {
		chunk := stream.Current()
		raw := chunk.Event.RawJSON()
		payload := chunk.Event.Payload.AsAny()
		if *debug {
			fmt.Fprintf(os.Stderr, "event: %s\n", raw)
		}
		if strings.Contains(raw, "\"step_type\":\"tool_execution\"") || strings.Contains(raw, "knowledge_search") {
			ragUsed = true
		}
		switch p := payload.(type) {
		case llamastackclient.TurnResponseEventPayloadStepComplete:
			detail := p.StepDetails.AsAny()
			if mr, ok := detail.(llamastackclient.MemoryRetrievalStep); ok {
				ragUsed = true
				ragSnippet = mr.InsertedContext.AsString()
			}
		case llamastackclient.TurnResponseEventPayloadTurnComplete:
			finalOut = p.Turn.OutputMessage.Content.AsString()
		}
	}
	if err := stream.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "stream error: %v\n", err)
		os.Exit(1)
	}

	if strings.TrimSpace(finalOut) == "" {
		fmt.Fprintln(os.Stderr, "no answer produced")
		os.Exit(1)
	}

	if *debug {
		fmt.Fprintf(os.Stderr, "RAG used: %v\n", ragUsed)
		if ragUsed && ragSnippet != "" {
			if len(ragSnippet) > 400 {
				ragSnippet = ragSnippet[:400] + "..."
			}
			fmt.Fprintf(os.Stderr, "RAG context: %s\n", ragSnippet)
		}
	}

	fmt.Println(finalOut)
}
