package ingester

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"tekton-assist/pkg/types"
)

type Loader struct{}

func NewLoader() *Loader { return &Loader{} }

func (l *Loader) LoadKB(path string) ([]types.KnowledgeEntry, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var raw []map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, err
	}
	entries := make([]types.KnowledgeEntry, 0, len(raw))
	for i, e := range raw {
		ke := types.KnowledgeEntry{
			ID:        generateID(e, i),
			Error:     getString(e, "error"),
			Context:   getString(e, "context"),
			Solution:  getString(e, "solution"),
			Reference: getString(e, "reference"),
			Source:    getString(e, "source"),
			Metadata:  getMap(e, "metadata"),
		}
		ke.ErrorType = inferErrorType(ke.Error)
		ke.Severity = inferSeverity(ke.Error, ke.ErrorType)
		ke.CombinedText = buildCombinedText(ke)
		entries = append(entries, ke)
	}
	return entries, nil
}

func buildCombinedText(e types.KnowledgeEntry) string {
	parts := []string{}
	if e.Error != "" {
		parts = append(parts, fmt.Sprintf("ERROR: %s", e.Error))
	}
	if e.Context != "" {
		parts = append(parts, fmt.Sprintf("CONTEXT: %s", e.Context))
	}
	if e.Solution != "" {
		parts = append(parts, fmt.Sprintf("SOLUTION: %s", e.Solution))
	}
	if tn := getStringFromMap(e.Metadata, "task_name"); tn != "" {
		parts = append(parts, fmt.Sprintf("TEKTON_TASK: %s", tn))
	}
	if fs := getStringFromMap(e.Metadata, "failed_step"); fs != "" {
		parts = append(parts, fmt.Sprintf("FAILED_STEP: %s", fs))
	}
	if repo := getStringFromMap(e.Metadata, "repository"); repo != "" {
		parts = append(parts, fmt.Sprintf("REPOSITORY: %s", repo))
	}
	if related, ok := e.Metadata["related_errors"].([]any); ok && len(related) > 0 {
		vals := make([]string, 0, len(related))
		for _, v := range related {
			if s, ok := v.(string); ok {
				vals = append(vals, s)
			}
		}
		sort.Strings(vals)
		parts = append(parts, fmt.Sprintf("RELATED_ERRORS: %s", strings.Join(vals, ", ")))
	}
	return strings.Join(parts, "\n")
}

func generateID(e map[string]any, idx int) string {
	if md, ok := e["metadata"].(map[string]any); ok {
		et, ok1 := md["error_type"].(string)
		ts, ok2 := md["timestamp"].(string)
		if ok1 && ok2 && len(ts) >= 10 {
			return fmt.Sprintf("%s-%s-%03d", strings.ReplaceAll(strings.ToLower(et), "_", "-"), ts[:10], idx)
		}
	}
	h := sha1.Sum([]byte(getString(e, "error")))
	return fmt.Sprintf("kb-entry-%03d-%s", idx, hex.EncodeToString(h[:])[:6])
}

func inferErrorType(errorText string) string {
	el := strings.ToLower(errorText)
	switch {
	case containsAny(el, []string{"oom", "exit code 137", "memory"}):
		return "OOM"
	case containsAny(el, []string{"registry", "push", "pull", "image"}):
		return "Registry_Push_Failure"
	case containsAny(el, []string{"build", "compile", "make"}):
		return "Build_Failure"
	case containsAny(el, []string{"version", "compatibility", "requires"}):
		return "Version_Compatibility"
	case containsAny(el, []string{"permission", "denied", "unauthorized"}):
		return "Permission_Denied"
	case containsAny(el, []string{"network", "timeout", "connection"}):
		return "Network_Issue"
	case containsAny(el, []string{"missing", "not found", "no such file"}):
		return "Missing_Files"
	default:
		return "Other"
	}
}

func inferSeverity(errorText, errorType string) string {
	el := strings.ToLower(errorText)
	if containsAny(el, []string{"oom", "crashloopbackoff", "panic"}) {
		return "critical"
	}
	switch errorType {
	case "Build_Failure", "Registry_Push_Failure", "Version_Compatibility":
		return "high"
	case "Permission_Denied", "Missing_Files":
		return "medium"
	}
	if containsAny(el, []string{"warning", "skipping"}) {
		return "low"
	}
	return "medium"
}

func containsAny(s string, needles []string) bool {
	for _, n := range needles {
		if strings.Contains(s, n) {
			return true
		}
	}
	return false
}

func getString(m map[string]any, k string) string {
	if v, ok := m[k]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getMap(m map[string]any, k string) map[string]any {
	if v, ok := m[k]; ok {
		if mm, ok := v.(map[string]any); ok {
			return mm
		}
	}
	return map[string]any{}
}

func getStringFromMap(m map[string]any, k string) string {
	if v, ok := m[k]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
