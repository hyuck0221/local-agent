package opencode

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/hyuck0221/local-agent/internal/platform"
)

// ProviderID is the key under `provider` that this tool owns in opencode.json.
// Only this key is ever created or replaced; unrelated providers stay intact.
const ProviderID = "local-agent"

// RegisterModel merges a local-agent provider entry into opencode.json,
// preserving any other providers and top-level keys the user may have set.
// baseURL should be the OpenAI-compatible root (e.g. http://localhost:11434/v1).
func RegisterModel(model, baseURL string) (string, error) {
	path, err := platform.OpencodeConfigPath()
	if err != nil {
		return "", err
	}

	raw, err := os.ReadFile(path)
	var root map[string]any
	switch {
	case err == nil:
		if err := json.Unmarshal(raw, &root); err != nil {
			return "", fmt.Errorf("parse %s: %w", path, err)
		}
	case errors.Is(err, fs.ErrNotExist):
		root = map[string]any{}
	default:
		return "", fmt.Errorf("read %s: %w", path, err)
	}
	if root == nil {
		root = map[string]any{}
	}

	providers, _ := root["provider"].(map[string]any)
	if providers == nil {
		providers = map[string]any{}
	}

	providers[ProviderID] = map[string]any{
		"npm":  "@ai-sdk/openai-compatible",
		"name": "Local Agent (Ollama)",
		"options": map[string]any{
			"baseURL": baseURL,
			"apiKey":  "ollama",
		},
		"models": map[string]any{
			model: map[string]any{"name": model},
		},
	}
	root["provider"] = providers

	if err := writeAtomic(path, root); err != nil {
		return "", err
	}
	return path, nil
}

// Registered returns the currently registered local-agent model name, or "" if none.
func Registered() (string, error) {
	path, err := platform.OpencodeConfigPath()
	if err != nil {
		return "", err
	}
	raw, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	var root map[string]any
	if err := json.Unmarshal(raw, &root); err != nil {
		return "", err
	}
	providers, _ := root["provider"].(map[string]any)
	entry, _ := providers[ProviderID].(map[string]any)
	models, _ := entry["models"].(map[string]any)
	for name := range models {
		return name, nil
	}
	return "", nil
}

func writeAtomic(path string, data map[string]any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".opencode-*.json")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(out); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return err
	}
	return os.Rename(tmpName, path)
}
