// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package llamastackclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/llamastack/llama-stack-client-go/internal/apijson"
	"github.com/llamastack/llama-stack-client-go/internal/requestconfig"
	"github.com/llamastack/llama-stack-client-go/option"
	"github.com/llamastack/llama-stack-client-go/packages/respjson"
	"github.com/llamastack/llama-stack-client-go/shared"
)

// AgentStepService contains methods and other services that help with interacting
// with the llama-stack-client API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAgentStepService] method instead.
type AgentStepService struct {
	Options []option.RequestOption
}

// NewAgentStepService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewAgentStepService(opts ...option.RequestOption) (r AgentStepService) {
	r = AgentStepService{}
	r.Options = opts
	return
}

// Retrieve an agent step by its ID.
func (r *AgentStepService) Get(ctx context.Context, stepID string, query AgentStepGetParams, opts ...option.RequestOption) (res *AgentStepGetResponse, err error) {
	opts = append(r.Options[:], opts...)
	if query.AgentID == "" {
		err = errors.New("missing required agent_id parameter")
		return
	}
	if query.SessionID == "" {
		err = errors.New("missing required session_id parameter")
		return
	}
	if query.TurnID == "" {
		err = errors.New("missing required turn_id parameter")
		return
	}
	if stepID == "" {
		err = errors.New("missing required step_id parameter")
		return
	}
	path := fmt.Sprintf("v1/agents/%s/session/%s/turn/%s/step/%s", query.AgentID, query.SessionID, query.TurnID, stepID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Response containing details of a specific agent step.
type AgentStepGetResponse struct {
	// The complete step data and execution details
	Step AgentStepGetResponseStepUnion `json:"step,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Step        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentStepGetResponse) RawJSON() string { return r.JSON.raw }
func (r *AgentStepGetResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AgentStepGetResponseStepUnion contains all possible properties and values from
// [InferenceStep], [ToolExecutionStep], [ShieldCallStep], [MemoryRetrievalStep].
//
// Use the [AgentStepGetResponseStepUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type AgentStepGetResponseStepUnion struct {
	// This field is from variant [InferenceStep].
	ModelResponse shared.CompletionMessage `json:"model_response"`
	StepID        string                   `json:"step_id"`
	// Any of "inference", "tool_execution", "shield_call", "memory_retrieval".
	StepType    string    `json:"step_type"`
	TurnID      string    `json:"turn_id"`
	CompletedAt time.Time `json:"completed_at"`
	StartedAt   time.Time `json:"started_at"`
	// This field is from variant [ToolExecutionStep].
	ToolCalls []shared.ToolCall `json:"tool_calls"`
	// This field is from variant [ToolExecutionStep].
	ToolResponses []ToolResponse `json:"tool_responses"`
	// This field is from variant [ShieldCallStep].
	Violation shared.SafetyViolation `json:"violation"`
	// This field is from variant [MemoryRetrievalStep].
	InsertedContext shared.InterleavedContentUnion `json:"inserted_context"`
	// This field is from variant [MemoryRetrievalStep].
	VectorDBIDs string `json:"vector_db_ids"`
	JSON        struct {
		ModelResponse   respjson.Field
		StepID          respjson.Field
		StepType        respjson.Field
		TurnID          respjson.Field
		CompletedAt     respjson.Field
		StartedAt       respjson.Field
		ToolCalls       respjson.Field
		ToolResponses   respjson.Field
		Violation       respjson.Field
		InsertedContext respjson.Field
		VectorDBIDs     respjson.Field
		raw             string
	} `json:"-"`
}

// anyAgentStepGetResponseStep is implemented by each variant of
// [AgentStepGetResponseStepUnion] to add type safety for the return type of
// [AgentStepGetResponseStepUnion.AsAny]
type anyAgentStepGetResponseStep interface {
	implAgentStepGetResponseStepUnion()
}

func (InferenceStep) implAgentStepGetResponseStepUnion()       {}
func (ToolExecutionStep) implAgentStepGetResponseStepUnion()   {}
func (ShieldCallStep) implAgentStepGetResponseStepUnion()      {}
func (MemoryRetrievalStep) implAgentStepGetResponseStepUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := AgentStepGetResponseStepUnion.AsAny().(type) {
//	case llamastackclient.InferenceStep:
//	case llamastackclient.ToolExecutionStep:
//	case llamastackclient.ShieldCallStep:
//	case llamastackclient.MemoryRetrievalStep:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u AgentStepGetResponseStepUnion) AsAny() anyAgentStepGetResponseStep {
	switch u.StepType {
	case "inference":
		return u.AsInference()
	case "tool_execution":
		return u.AsToolExecution()
	case "shield_call":
		return u.AsShieldCall()
	case "memory_retrieval":
		return u.AsMemoryRetrieval()
	}
	return nil
}

func (u AgentStepGetResponseStepUnion) AsInference() (v InferenceStep) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentStepGetResponseStepUnion) AsToolExecution() (v ToolExecutionStep) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentStepGetResponseStepUnion) AsShieldCall() (v ShieldCallStep) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentStepGetResponseStepUnion) AsMemoryRetrieval() (v MemoryRetrievalStep) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u AgentStepGetResponseStepUnion) RawJSON() string { return u.JSON.raw }

func (r *AgentStepGetResponseStepUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AgentStepGetParams struct {
	AgentID   string `path:"agent_id,required" json:"-"`
	SessionID string `path:"session_id,required" json:"-"`
	TurnID    string `path:"turn_id,required" json:"-"`
	paramObj
}
