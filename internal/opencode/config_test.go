package opencode

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestMergePreservesExistingProviders verifies the merge logic does not
// clobber other providers or top-level keys the user may have added.
func TestMergePreservesExistingProviders(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(dir, ".config"))

	path := filepath.Join(dir, ".config", "opencode", "opencode.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	existing := map[string]any{
		"theme": "dracula",
		"provider": map[string]any{
			"anthropic": map[string]any{"api_key": "secret"},
		},
	}
	raw, _ := json.MarshalIndent(existing, "", "  ")
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := RegisterModel("qwen2.5-coder:7b", "http://localhost:11434/v1"); err != nil {
		t.Fatalf("RegisterModel: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var parsed map[string]any
	if err := json.Unmarshal(got, &parsed); err != nil {
		t.Fatal(err)
	}
	if parsed["theme"] != "dracula" {
		t.Errorf("top-level key lost: got %v", parsed["theme"])
	}
	providers := parsed["provider"].(map[string]any)
	if _, ok := providers["anthropic"]; !ok {
		t.Error("existing anthropic provider was removed")
	}
	if _, ok := providers[ProviderID]; !ok {
		t.Error("local-agent provider was not added")
	}
}

// TestIdempotent verifies repeated calls update the same entry rather than accumulating.
func TestIdempotent(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	if _, err := RegisterModel("a", "http://localhost:11434/v1"); err != nil {
		t.Fatal(err)
	}
	if _, err := RegisterModel("b", "http://localhost:11434/v1"); err != nil {
		t.Fatal(err)
	}
	name, err := Registered()
	if err != nil {
		t.Fatal(err)
	}
	if name != "b" {
		t.Errorf("expected b, got %q", name)
	}
}
