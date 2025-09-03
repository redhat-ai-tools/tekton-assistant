package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func printResult(cfg *Config, result map[string]any) error {
	if cfg.JSONOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}
	status := fmt.Sprint(result["status"])
	action := fmt.Sprint(result["action"])
	emoji := map[string]string{"success": "✅", "error": "❌", "warning": "⚠️"}[status]
	fmt.Printf("%s %s %s\n", emoji, capitalize(action), capitalize(status))
	fmt.Println("==================================================")
	if msg, ok := result["message"].(string); ok && msg != "" {
		fmt.Printf("Message: %s\n", msg)
	}

	// Pretty-print action-specific payloads
	switch action {
	case "stats":
		if stats, ok := result["stats"].(map[string]any); ok {
			keys := make([]string, 0, len(stats))
			for k := range stats {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				fmt.Printf("%s: %v\n", k, stats[k])
			}
		}
	case "validate":
		if v, ok := result["validation"].(map[string]any); ok {
			keys := make([]string, 0, len(v))
			for k := range v {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				val := v[k]
				if b, ok := val.(bool); ok {
					icon := "❌"
					if b {
						icon = "✅"
					}
					fmt.Printf("%s %s: %v\n", icon, k, b)
				} else {
					fmt.Printf("%s: %v\n", k, val)
				}
			}
		}
	}
	return nil
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = toUpper(r[0])
	return string(r)
}

func toUpper(r rune) rune {
	if r >= 'a' && r <= 'z' {
		return r - ('a' - 'A')
	}
	return r
}

func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

func trunc(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

func resolveKBPath(p string) string {
	candidates := []string{
		p,
		filepath.Join("data", "kb.json"),
		filepath.Join("..", "data", "kb.json"),
	}
	for _, cand := range candidates {
		if fileExists(cand) {
			return cand
		}
	}
	return p
}
