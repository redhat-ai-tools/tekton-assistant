// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package llamastackclient

import (
	"github.com/llamastack/llama-stack-client-go/internal/apierror"
	"github.com/llamastack/llama-stack-client-go/packages/param"
	"github.com/llamastack/llama-stack-client-go/shared"
)

// aliased to make [param.APIUnion] private when embedding
type paramUnion = param.APIUnion

// aliased to make [param.APIObject] private when embedding
type paramObj = param.APIObject

type Error = apierror.Error

// Configuration for an agent.
//
// This is an alias to an internal type.
type AgentConfig = shared.AgentConfig

// Whether tool use is required or automatic. This is a hint to the model which may
// not be followed. It depends on the Instruction Following capabilities of the
// model.
//
// This is an alias to an internal type.
type AgentConfigToolChoice = shared.AgentConfigToolChoice

// Equals "auto"
const AgentConfigToolChoiceAuto = shared.AgentConfigToolChoiceAuto

// Equals "required"
const AgentConfigToolChoiceRequired = shared.AgentConfigToolChoiceRequired

// Equals "none"
const AgentConfigToolChoiceNone = shared.AgentConfigToolChoiceNone

// Configuration for tool use.
//
// This is an alias to an internal type.
type AgentConfigToolConfig = shared.AgentConfigToolConfig

// Prompt format for calling custom / zero shot tools.
//
// This is an alias to an internal type.
type AgentConfigToolPromptFormat = shared.AgentConfigToolPromptFormat

// Equals "json"
const AgentConfigToolPromptFormatJson = shared.AgentConfigToolPromptFormatJson

// Equals "function_tag"
const AgentConfigToolPromptFormatFunctionTag = shared.AgentConfigToolPromptFormatFunctionTag

// Equals "python_list"
const AgentConfigToolPromptFormatPythonList = shared.AgentConfigToolPromptFormatPythonList

// This is an alias to an internal type.
type AgentConfigToolgroupUnion = shared.AgentConfigToolgroupUnion

// This is an alias to an internal type.
type AgentConfigToolgroupAgentToolGroupWithArgs = shared.AgentConfigToolgroupAgentToolGroupWithArgs

// This is an alias to an internal type.
type AgentConfigToolgroupAgentToolGroupWithArgsArgUnion = shared.AgentConfigToolgroupAgentToolGroupWithArgsArgUnion

// Configuration for an agent.
//
// This is an alias to an internal type.
type AgentConfigParam = shared.AgentConfigParam

// Configuration for tool use.
//
// This is an alias to an internal type.
type AgentConfigToolConfigParam = shared.AgentConfigToolConfigParam

// This is an alias to an internal type.
type AgentConfigToolgroupUnionParam = shared.AgentConfigToolgroupUnionParam

// This is an alias to an internal type.
type AgentConfigToolgroupAgentToolGroupWithArgsParam = shared.AgentConfigToolgroupAgentToolGroupWithArgsParam

// This is an alias to an internal type.
type AgentConfigToolgroupAgentToolGroupWithArgsArgUnionParam = shared.AgentConfigToolgroupAgentToolGroupWithArgsArgUnionParam

// Response from a batch completion request.
//
// This is an alias to an internal type.
type BatchCompletion = shared.BatchCompletion

// Response from a chat completion request.
//
// This is an alias to an internal type.
type ChatCompletionResponse = shared.ChatCompletionResponse

// Log probabilities for generated tokens.
//
// This is an alias to an internal type.
type ChatCompletionResponseLogprob = shared.ChatCompletionResponseLogprob

// A metric value included in API responses.
//
// This is an alias to an internal type.
type ChatCompletionResponseMetric = shared.ChatCompletionResponseMetric

// A message containing the model's (assistant) response in a chat conversation.
//
// This is an alias to an internal type.
type CompletionMessage = shared.CompletionMessage

// Reason why the model stopped generating. Options are: -
// `StopReason.end_of_turn`: The model finished generating the entire response. -
// `StopReason.end_of_message`: The model finished generating but generated a
// partial response -- usually, a tool call. The user may call the tool and
// continue the conversation with the tool's response. -
// `StopReason.out_of_tokens`: The model ran out of token budget.
//
// This is an alias to an internal type.
type CompletionMessageStopReason = shared.CompletionMessageStopReason

// Equals "end_of_turn"
const CompletionMessageStopReasonEndOfTurn = shared.CompletionMessageStopReasonEndOfTurn

// Equals "end_of_message"
const CompletionMessageStopReasonEndOfMessage = shared.CompletionMessageStopReasonEndOfMessage

// Equals "out_of_tokens"
const CompletionMessageStopReasonOutOfTokens = shared.CompletionMessageStopReasonOutOfTokens

// A message containing the model's (assistant) response in a chat conversation.
//
// This is an alias to an internal type.
type CompletionMessageParam = shared.CompletionMessageParam

// A text content delta for streaming responses.
//
// This is an alias to an internal type.
type ContentDeltaUnion = shared.ContentDeltaUnion

// A text content delta for streaming responses.
//
// This is an alias to an internal type.
type ContentDeltaText = shared.ContentDeltaText

// An image content delta for streaming responses.
//
// This is an alias to an internal type.
type ContentDeltaImage = shared.ContentDeltaImage

// A tool call content delta for streaming responses.
//
// This is an alias to an internal type.
type ContentDeltaToolCall = shared.ContentDeltaToolCall

// A document to be used for document ingestion in the RAG Tool.
//
// This is an alias to an internal type.
type DocumentParam = shared.DocumentParam

// The content of the document.
//
// This is an alias to an internal type.
type DocumentContentUnionParam = shared.DocumentContentUnionParam

// A image content item
//
// This is an alias to an internal type.
type DocumentContentImageContentItemParam = shared.DocumentContentImageContentItemParam

// Image as a base64 encoded string or an URL
//
// This is an alias to an internal type.
type DocumentContentImageContentItemImageParam = shared.DocumentContentImageContentItemImageParam

// A URL of the image or data URL in the format of data:image/{type};base64,{data}.
// Note that URL could have length limits.
//
// This is an alias to an internal type.
type DocumentContentImageContentItemImageURLParam = shared.DocumentContentImageContentItemImageURLParam

// A text content item
//
// This is an alias to an internal type.
type DocumentContentTextContentItemParam = shared.DocumentContentTextContentItemParam

// A URL reference to external content.
//
// This is an alias to an internal type.
type DocumentContentURLParam = shared.DocumentContentURLParam

// This is an alias to an internal type.
type DocumentMetadataUnionParam = shared.DocumentMetadataUnionParam

// A image content item
//
// This is an alias to an internal type.
type InterleavedContentUnion = shared.InterleavedContentUnion

// A image content item
//
// This is an alias to an internal type.
type InterleavedContentImageContentItem = shared.InterleavedContentImageContentItem

// Image as a base64 encoded string or an URL
//
// This is an alias to an internal type.
type InterleavedContentImageContentItemImage = shared.InterleavedContentImageContentItemImage

// A URL of the image or data URL in the format of data:image/{type};base64,{data}.
// Note that URL could have length limits.
//
// This is an alias to an internal type.
type InterleavedContentImageContentItemImageURL = shared.InterleavedContentImageContentItemImageURL

// A text content item
//
// This is an alias to an internal type.
type InterleavedContentTextContentItem = shared.InterleavedContentTextContentItem

// A image content item
//
// This is an alias to an internal type.
type InterleavedContentUnionParam = shared.InterleavedContentUnionParam

// A image content item
//
// This is an alias to an internal type.
type InterleavedContentImageContentItemParam = shared.InterleavedContentImageContentItemParam

// Image as a base64 encoded string or an URL
//
// This is an alias to an internal type.
type InterleavedContentImageContentItemImageParam = shared.InterleavedContentImageContentItemImageParam

// A URL of the image or data URL in the format of data:image/{type};base64,{data}.
// Note that URL could have length limits.
//
// This is an alias to an internal type.
type InterleavedContentImageContentItemImageURLParam = shared.InterleavedContentImageContentItemImageURLParam

// A text content item
//
// This is an alias to an internal type.
type InterleavedContentTextContentItemParam = shared.InterleavedContentTextContentItemParam

// A image content item
//
// This is an alias to an internal type.
type InterleavedContentItemUnion = shared.InterleavedContentItemUnion

// A image content item
//
// This is an alias to an internal type.
type InterleavedContentItemImage = shared.InterleavedContentItemImage

// Image as a base64 encoded string or an URL
//
// This is an alias to an internal type.
type InterleavedContentItemImageImage = shared.InterleavedContentItemImageImage

// A URL of the image or data URL in the format of data:image/{type};base64,{data}.
// Note that URL could have length limits.
//
// This is an alias to an internal type.
type InterleavedContentItemImageImageURL = shared.InterleavedContentItemImageImageURL

// A text content item
//
// This is an alias to an internal type.
type InterleavedContentItemText = shared.InterleavedContentItemText

// A image content item
//
// This is an alias to an internal type.
type InterleavedContentItemUnionParam = shared.InterleavedContentItemUnionParam

// A image content item
//
// This is an alias to an internal type.
type InterleavedContentItemImageParam = shared.InterleavedContentItemImageParam

// Image as a base64 encoded string or an URL
//
// This is an alias to an internal type.
type InterleavedContentItemImageImageParam = shared.InterleavedContentItemImageImageParam

// A URL of the image or data URL in the format of data:image/{type};base64,{data}.
// Note that URL could have length limits.
//
// This is an alias to an internal type.
type InterleavedContentItemImageImageURLParam = shared.InterleavedContentItemImageImageURLParam

// A text content item
//
// This is an alias to an internal type.
type InterleavedContentItemTextParam = shared.InterleavedContentItemTextParam

// A message from the user in a chat conversation.
//
// This is an alias to an internal type.
type MessageUnionParam = shared.MessageUnionParam

// Configuration for the RAG query generation.
//
// This is an alias to an internal type.
type QueryConfigParam = shared.QueryConfigParam

// Search mode for retrieval—either "vector", "keyword", or "hybrid". Default
// "vector".
//
// This is an alias to an internal type.
type QueryConfigMode = shared.QueryConfigMode

// Equals "vector"
const QueryConfigModeVector = shared.QueryConfigModeVector

// Equals "keyword"
const QueryConfigModeKeyword = shared.QueryConfigModeKeyword

// Equals "hybrid"
const QueryConfigModeHybrid = shared.QueryConfigModeHybrid

// Configuration for the ranker to use in hybrid search. Defaults to RRF ranker.
//
// This is an alias to an internal type.
type QueryConfigRankerUnionParam = shared.QueryConfigRankerUnionParam

// Reciprocal Rank Fusion (RRF) ranker configuration.
//
// This is an alias to an internal type.
type QueryConfigRankerRrfParam = shared.QueryConfigRankerRrfParam

// Weighted ranker configuration that combines vector and keyword scores.
//
// This is an alias to an internal type.
type QueryConfigRankerWeightedParam = shared.QueryConfigRankerWeightedParam

// Configuration for the default RAG query generator.
//
// This is an alias to an internal type.
type QueryGeneratorConfigUnionParam = shared.QueryGeneratorConfigUnionParam

// Configuration for the default RAG query generator.
//
// This is an alias to an internal type.
type QueryGeneratorConfigDefaultParam = shared.QueryGeneratorConfigDefaultParam

// Configuration for the LLM-based RAG query generator.
//
// This is an alias to an internal type.
type QueryGeneratorConfigLlmParam = shared.QueryGeneratorConfigLlmParam

// Result of a RAG query containing retrieved content and metadata.
//
// This is an alias to an internal type.
type QueryResult = shared.QueryResult

// This is an alias to an internal type.
type QueryResultMetadataUnion = shared.QueryResultMetadataUnion

// Configuration for JSON schema-guided response generation.
//
// This is an alias to an internal type.
type ResponseFormatUnion = shared.ResponseFormatUnion

// Configuration for JSON schema-guided response generation.
//
// This is an alias to an internal type.
type ResponseFormatJsonSchema = shared.ResponseFormatJsonSchema

// This is an alias to an internal type.
type ResponseFormatJsonSchemaJsonSchemaUnion = shared.ResponseFormatJsonSchemaJsonSchemaUnion

// Configuration for grammar-guided response generation.
//
// This is an alias to an internal type.
type ResponseFormatGrammar = shared.ResponseFormatGrammar

// This is an alias to an internal type.
type ResponseFormatGrammarBnfUnion = shared.ResponseFormatGrammarBnfUnion

// Configuration for JSON schema-guided response generation.
//
// This is an alias to an internal type.
type ResponseFormatUnionParam = shared.ResponseFormatUnionParam

// Configuration for JSON schema-guided response generation.
//
// This is an alias to an internal type.
type ResponseFormatJsonSchemaParam = shared.ResponseFormatJsonSchemaParam

// This is an alias to an internal type.
type ResponseFormatJsonSchemaJsonSchemaUnionParam = shared.ResponseFormatJsonSchemaJsonSchemaUnionParam

// Configuration for grammar-guided response generation.
//
// This is an alias to an internal type.
type ResponseFormatGrammarParam = shared.ResponseFormatGrammarParam

// This is an alias to an internal type.
type ResponseFormatGrammarBnfUnionParam = shared.ResponseFormatGrammarBnfUnionParam

// This is an alias to an internal type.
type ReturnType = shared.ReturnType

// This is an alias to an internal type.
type ReturnTypeType = shared.ReturnTypeType

// Equals "string"
const ReturnTypeTypeString = shared.ReturnTypeTypeString

// Equals "number"
const ReturnTypeTypeNumber = shared.ReturnTypeTypeNumber

// Equals "boolean"
const ReturnTypeTypeBoolean = shared.ReturnTypeTypeBoolean

// Equals "array"
const ReturnTypeTypeArray = shared.ReturnTypeTypeArray

// Equals "object"
const ReturnTypeTypeObject = shared.ReturnTypeTypeObject

// Equals "json"
const ReturnTypeTypeJson = shared.ReturnTypeTypeJson

// Equals "union"
const ReturnTypeTypeUnion = shared.ReturnTypeTypeUnion

// Equals "chat_completion_input"
const ReturnTypeTypeChatCompletionInput = shared.ReturnTypeTypeChatCompletionInput

// Equals "completion_input"
const ReturnTypeTypeCompletionInput = shared.ReturnTypeTypeCompletionInput

// Equals "agent_turn_input"
const ReturnTypeTypeAgentTurnInput = shared.ReturnTypeTypeAgentTurnInput

// This is an alias to an internal type.
type ReturnTypeParam = shared.ReturnTypeParam

// Details of a safety violation detected by content moderation.
//
// This is an alias to an internal type.
type SafetyViolation = shared.SafetyViolation

// This is an alias to an internal type.
type SafetyViolationMetadataUnion = shared.SafetyViolationMetadataUnion

// Severity level of the violation
//
// This is an alias to an internal type.
type SafetyViolationViolationLevel = shared.SafetyViolationViolationLevel

// Equals "info"
const SafetyViolationViolationLevelInfo = shared.SafetyViolationViolationLevelInfo

// Equals "warn"
const SafetyViolationViolationLevelWarn = shared.SafetyViolationViolationLevelWarn

// Equals "error"
const SafetyViolationViolationLevelError = shared.SafetyViolationViolationLevelError

// Sampling parameters.
//
// This is an alias to an internal type.
type SamplingParamsResp = shared.SamplingParamsResp

// The sampling strategy.
//
// This is an alias to an internal type.
type SamplingParamsStrategyUnionResp = shared.SamplingParamsStrategyUnionResp

// Greedy sampling strategy that selects the highest probability token at each
// step.
//
// This is an alias to an internal type.
type SamplingParamsStrategyGreedyResp = shared.SamplingParamsStrategyGreedyResp

// Top-p (nucleus) sampling strategy that samples from the smallest set of tokens
// with cumulative probability >= p.
//
// This is an alias to an internal type.
type SamplingParamsStrategyTopPResp = shared.SamplingParamsStrategyTopPResp

// Top-k sampling strategy that restricts sampling to the k most likely tokens.
//
// This is an alias to an internal type.
type SamplingParamsStrategyTopKResp = shared.SamplingParamsStrategyTopKResp

// Sampling parameters.
//
// This is an alias to an internal type.
type SamplingParams = shared.SamplingParams

// The sampling strategy.
//
// This is an alias to an internal type.
type SamplingParamsStrategyUnion = shared.SamplingParamsStrategyUnion

// Greedy sampling strategy that selects the highest probability token at each
// step.
//
// This is an alias to an internal type.
type SamplingParamsStrategyGreedy = shared.SamplingParamsStrategyGreedy

// Top-p (nucleus) sampling strategy that samples from the smallest set of tokens
// with cumulative probability >= p.
//
// This is an alias to an internal type.
type SamplingParamsStrategyTopP = shared.SamplingParamsStrategyTopP

// Top-k sampling strategy that restricts sampling to the k most likely tokens.
//
// This is an alias to an internal type.
type SamplingParamsStrategyTopK = shared.SamplingParamsStrategyTopK

// A scoring result for a single row.
//
// This is an alias to an internal type.
type ScoringResult = shared.ScoringResult

// This is an alias to an internal type.
type ScoringResultAggregatedResultUnion = shared.ScoringResultAggregatedResultUnion

// This is an alias to an internal type.
type ScoringResultScoreRowUnion = shared.ScoringResultScoreRowUnion

// Response from a completion request.
//
// This is an alias to an internal type.
type SharedCompletionResponse = shared.SharedCompletionResponse

// Reason why generation stopped
//
// This is an alias to an internal type.
type SharedCompletionResponseStopReason = shared.SharedCompletionResponseStopReason

// Equals "end_of_turn"
const SharedCompletionResponseStopReasonEndOfTurn = shared.SharedCompletionResponseStopReasonEndOfTurn

// Equals "end_of_message"
const SharedCompletionResponseStopReasonEndOfMessage = shared.SharedCompletionResponseStopReasonEndOfMessage

// Equals "out_of_tokens"
const SharedCompletionResponseStopReasonOutOfTokens = shared.SharedCompletionResponseStopReasonOutOfTokens

// Log probabilities for generated tokens.
//
// This is an alias to an internal type.
type SharedCompletionResponseLogprob = shared.SharedCompletionResponseLogprob

// A metric value included in API responses.
//
// This is an alias to an internal type.
type SharedCompletionResponseMetric = shared.SharedCompletionResponseMetric

// Tool definition used in runtime contexts.
//
// This is an alias to an internal type.
type SharedToolDef = shared.SharedToolDef

// This is an alias to an internal type.
type SharedToolDefMetadataUnion = shared.SharedToolDefMetadataUnion

// Parameter definition for a tool.
//
// This is an alias to an internal type.
type SharedToolDefParameter = shared.SharedToolDefParameter

// (Optional) Default value for the parameter if not provided
//
// This is an alias to an internal type.
type SharedToolDefParameterDefaultUnion = shared.SharedToolDefParameterDefaultUnion

// Tool definition used in runtime contexts.
//
// This is an alias to an internal type.
type SharedToolDefParam = shared.SharedToolDefParam

// This is an alias to an internal type.
type SharedToolDefMetadataUnionParam = shared.SharedToolDefMetadataUnionParam

// Parameter definition for a tool.
//
// This is an alias to an internal type.
type SharedToolDefParameterParam = shared.SharedToolDefParameterParam

// (Optional) Default value for the parameter if not provided
//
// This is an alias to an internal type.
type SharedToolDefParameterDefaultUnionParam = shared.SharedToolDefParameterDefaultUnionParam

// A system message providing instructions or context to the model.
//
// This is an alias to an internal type.
type SystemMessageParam = shared.SystemMessageParam

// This is an alias to an internal type.
type ToolCall = shared.ToolCall

// This is an alias to an internal type.
type ToolCallArgumentsUnion = shared.ToolCallArgumentsUnion

// This is an alias to an internal type.
type ToolCallArgumentsMapItemUnion = shared.ToolCallArgumentsMapItemUnion

// This is an alias to an internal type.
type ToolCallArgumentsMapItemArrayItemUnion = shared.ToolCallArgumentsMapItemArrayItemUnion

// This is an alias to an internal type.
type ToolCallArgumentsMapItemMapItemUnion = shared.ToolCallArgumentsMapItemMapItemUnion

// This is an alias to an internal type.
type ToolCallToolName = shared.ToolCallToolName

// Equals "brave_search"
const ToolCallToolNameBraveSearch = shared.ToolCallToolNameBraveSearch

// Equals "wolfram_alpha"
const ToolCallToolNameWolframAlpha = shared.ToolCallToolNameWolframAlpha

// Equals "photogen"
const ToolCallToolNamePhotogen = shared.ToolCallToolNamePhotogen

// Equals "code_interpreter"
const ToolCallToolNameCodeInterpreter = shared.ToolCallToolNameCodeInterpreter

// This is an alias to an internal type.
type ToolCallParam = shared.ToolCallParam

// This is an alias to an internal type.
type ToolCallArgumentsUnionParam = shared.ToolCallArgumentsUnionParam

// This is an alias to an internal type.
type ToolCallArgumentsMapItemUnionParam = shared.ToolCallArgumentsMapItemUnionParam

// This is an alias to an internal type.
type ToolCallArgumentsMapItemArrayItemUnionParam = shared.ToolCallArgumentsMapItemArrayItemUnionParam

// This is an alias to an internal type.
type ToolCallArgumentsMapItemMapItemUnionParam = shared.ToolCallArgumentsMapItemMapItemUnionParam

// Either an in-progress tool call string or the final parsed tool call
//
// This is an alias to an internal type.
type ToolCallOrStringUnion = shared.ToolCallOrStringUnion

// This is an alias to an internal type.
type ToolParamDefinition = shared.ToolParamDefinition

// This is an alias to an internal type.
type ToolParamDefinitionDefaultUnion = shared.ToolParamDefinitionDefaultUnion

// A message representing the result of a tool invocation.
//
// This is an alias to an internal type.
type ToolResponseMessage = shared.ToolResponseMessage

// A message representing the result of a tool invocation.
//
// This is an alias to an internal type.
type ToolResponseMessageParam = shared.ToolResponseMessageParam

// A message from the user in a chat conversation.
//
// This is an alias to an internal type.
type UserMessage = shared.UserMessage

// A message from the user in a chat conversation.
//
// This is an alias to an internal type.
type UserMessageParam = shared.UserMessageParam
