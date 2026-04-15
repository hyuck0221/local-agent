package tui

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

type Model struct {
	Tag   string
	Label string // short description, no tag prefix
}

// Recommended is a curated list of Ollama models surfaced in the picker.
// Every entry advertises the `tools` capability in `ollama show`, which is
// what opencode's agent loop needs. Gemma 3 / DeepSeek-R1 / Phi3 are
// intentionally excluded because they cannot emit structured tool calls.
//
// Gemma 4 *does* advertise tool support but in practice hallucinates tools
// outside the provided schema (e.g. `google:search`). It stays in the list
// for multimodal use cases but is labelled accordingly.
var Recommended = []Model{
	{Tag: "qwen2.5-coder:7b", Label: "best small coder, ~5GB (default)"},
	{Tag: "qwen2.5-coder:14b", Label: "higher quality coder, ~9GB, 16GB+ RAM"},
	{Tag: "qwen3:8b", Label: "modern general-purpose, ~5GB"},
	{Tag: "llama3.1:8b", Label: "Meta baseline with strong tools, ~5GB"},
	{Tag: "llama3.3:70b", Label: "top quality, ~40GB, workstation only"},
	{Tag: "mistral-nemo:12b", Label: "fast with solid tool calling, ~7GB"},
	{Tag: "deepseek-coder-v2:16b", Label: "powerful MoE coder, ~9GB"},
	{Tag: "command-r:35b", Label: "RAG/tools specialist, ~20GB"},
	{Tag: "gemma4:e4b", Label: "multimodal, ~8B — tools unstable in opencode"},
	{Tag: "gemma4:e2b", Label: "small multimodal, ~5B — tools unstable in opencode"},
}

// PickModel renders a unified select prompt where each row shows whether
// the model is already installed locally. Order: recommended-installed,
// recommended-not-installed, then any locally installed models that are not
// in the recommended list.
func PickModel(installed []string) (string, error) {
	instSet := make(map[string]bool, len(installed))
	for _, t := range installed {
		instSet[t] = true
	}
	recSet := make(map[string]bool, len(Recommended))
	for _, m := range Recommended {
		recSet[m.Tag] = true
	}

	var recInstalled, recMissing []Model
	for _, m := range Recommended {
		if instSet[m.Tag] {
			recInstalled = append(recInstalled, m)
		} else {
			recMissing = append(recMissing, m)
		}
	}
	var extra []string
	for _, t := range installed {
		if !recSet[t] {
			extra = append(extra, t)
		}
	}

	// Column width for the tag, so labels align.
	tagWidth := 0
	widen := func(s string) {
		if len(s) > tagWidth {
			tagWidth = len(s)
		}
	}
	for _, m := range Recommended {
		widen(m.Tag)
	}
	for _, t := range extra {
		widen(t)
	}

	row := func(mark, tag, label string) string {
		return fmt.Sprintf("%s  %-*s  %s", mark, tagWidth, tag, label)
	}

	opts := make([]huh.Option[string], 0, len(Recommended)+len(extra))
	for _, m := range recInstalled {
		opts = append(opts, huh.NewOption(row("●", m.Tag, m.Label), m.Tag))
	}
	for _, m := range recMissing {
		opts = append(opts, huh.NewOption(row("○", m.Tag, m.Label), m.Tag))
	}
	for _, t := range extra {
		opts = append(opts, huh.NewOption(row("●", t, "locally installed"), t))
	}

	var choice string
	err := huh.NewSelect[string]().
		Title("Pick a model to run locally").
		Description("●  already installed      ○  will be downloaded      (opencode needs tool-calling)").
		Options(opts...).
		Height(min(len(opts)+4, 18)).
		Value(&choice).
		Run()
	if err != nil {
		return "", err
	}
	return choice, nil
}
