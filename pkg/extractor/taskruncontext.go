package extractor

import (
	"context"
	"tekton-assistant/pkg/types"
)

// ExtractTaskRunContext is a placeholder implementation.
func ExtractTaskRunContext(ctx context.Context, taskrunID, namespace string) (*types.TaskRunContext, error) {
	return &types.TaskRunContext{
		TaskRunID: taskrunID,
		Namespace: namespace,
		Succeeded: false,
		FailedStep: &types.FailedStep{
			Name:     "unknown",
			ExitCode: 1,
		},
		Error: &types.ErrorContext{
			Type:    "NotImplemented",
			Message: "Extractor not implemented",
			Context: "",
		},
	}, nil
}
