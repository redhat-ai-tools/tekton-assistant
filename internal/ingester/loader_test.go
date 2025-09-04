package ingester

import (
	"os"
	"path/filepath"
	"testing"

	ttypes "tekton-assistant/internal/types"
)

func TestInferErrorType(t *testing.T) {
	cases := []struct{ in, want string }{
		{"Process killed due to OOM and exit code 137", "OOM"},
		{"Failed to push image to registry", "Registry_Push_Failure"},
		{"Build failed: make returned non-zero", "Build_Failure"},
		{"Requires version >= 1.2 due to compatibility", "Version_Compatibility"},
		{"Permission denied on secret", "Permission_Denied"},
		{"Connection timeout on network", "Network_Issue"},
		{"File not found: no such file", "Missing_Files"},
		{"Some random error", "Other"},
	}
	for _, c := range cases {
		got := inferErrorType(c.in)
		if got != c.want {
			t.Fatalf("inferErrorType(%q)=%q want %q", c.in, got, c.want)
		}
	}
}

func TestInferSeverity(t *testing.T) {
	if got := inferSeverity("panic in container", "Other"); got != "critical" {
		t.Fatalf("severity panic=%q want critical", got)
	}
	if got := inferSeverity("", "Build_Failure"); got != "high" {
		t.Fatalf("severity build=%q want high", got)
	}
	if got := inferSeverity("", "Permission_Denied"); got != "medium" {
		t.Fatalf("severity perm=%q want medium", got)
	}
	if got := inferSeverity("warning: skipping step", "Other"); got != "low" {
		t.Fatalf("severity warning=%q want low", got)
	}
}

func TestBuildCombinedText(t *testing.T) {
	e := ttypes.KnowledgeEntry{
		Error:    "Build failed",
		Context:  "during step compile",
		Solution: "Ensure dependencies installed",
		Metadata: map[string]any{
			"task_name":      "build-task",
			"failed_step":    "compile",
			"repository":     "org/repo",
			"related_errors": []any{"E1", "E2"},
		},
	}
	got := buildCombinedText(e)
	// Basic line presence checks
	wantSubs := []string{
		"ERROR: Build failed",
		"CONTEXT: during step compile",
		"SOLUTION: Ensure dependencies installed",
		"TEKTON_TASK: build-task",
		"FAILED_STEP: compile",
		"REPOSITORY: org/repo",
		"RELATED_ERRORS: E1, E2",
	}
	for _, sub := range wantSubs {
		if !containsAny(got, []string{sub}) {
			t.Fatalf("combined_text missing %q. got=\n%s", sub, got)
		}
	}
}

func TestLoadKB_FromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "kb.json")
	content := `[
      {
        "error": "Registry push failed with 401",
        "context": "pushing image",
        "solution": "login to registry",
        "reference": "docs",
        "source": "tests",
        "metadata": {"task_name": "push", "repository": "org/repo"}
      }
    ]`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	l := NewLoader()
	entries, err := l.LoadKB(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("want 1 entry got %d", len(entries))
	}
	e := entries[0]
	if e.ErrorType != "Registry_Push_Failure" {
		t.Fatalf("wrong error type: %s", e.ErrorType)
	}
	if e.CombinedText == "" {
		t.Fatalf("combined text should not be empty")
	}
	if len(e.ID) == 0 {
		t.Fatalf("id should be set")
	}
}
