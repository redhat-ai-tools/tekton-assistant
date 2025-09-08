package types

import "context"

// FailedStep represents a failed step in a TaskRun
type FailedStep struct {
	Name     string `json:"name"`
	ExitCode int    `json:"exit_code"`
}

// ErrorContext represents error details from a TaskRun
type ErrorContext struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Context string `json:"context"`
}

// TaskRunContext represents the output of TaskRun extraction
type TaskRunContext struct {
	TaskRunID  string        `json:"taskrun_id"`
	Namespace  string        `json:"namespace"`
	Succeeded  bool          `json:"succeeded"`
	FailedStep *FailedStep   `json:"failed_step,omitempty"`
	Error      *ErrorContext `json:"error,omitempty"`
}

// Extractor defines interface for TaskRun extraction
type Extractor interface {
	Extract(ctx context.Context, taskrunID, namespace string) (*TaskRunContext, error)
}
