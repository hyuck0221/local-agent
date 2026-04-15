package tui

import "github.com/charmbracelet/huh"

// Recommended is a curated list of Ollama models surfaced in the picker.
// Tool-calling-capable models are preferred because opencode's agent loop
// requires them; the few entries that do NOT support tools are explicitly
// labelled "chat only" so users know they will not work with opencode's
// editing features.
var Recommended = []Model{
	{Tag: "qwen2.5-coder:7b", Label: "Qwen2.5 Coder 7B — best small coder, ~5GB (default pick)"},
	{Tag: "qwen2.5-coder:14b", Label: "Qwen2.5 Coder 14B — higher quality, ~9GB, 16GB+ RAM"},
	{Tag: "qwen3:8b", Label: "Qwen 3 8B — modern general-purpose, ~5GB"},
	{Tag: "llama3.1:8b", Label: "Llama 3.1 8B — Meta baseline with strong tools, ~5GB"},
	{Tag: "llama3.3:70b", Label: "Llama 3.3 70B — top quality, ~40GB, workstation only"},
	{Tag: "mistral-nemo:12b", Label: "Mistral Nemo 12B — fast with solid tool calling, ~7GB"},
	{Tag: "deepseek-coder-v2:16b", Label: "DeepSeek Coder V2 16B — powerful MoE coder, ~9GB"},
	{Tag: "command-r:35b", Label: "Command R 35B — RAG/tools specialist, ~20GB"},
	{Tag: "gemma4:e4b", Label: "Gemma 4 E4B — multilingual, chat only (no tool calling)"},
	{Tag: "gemma4:e2b", Label: "Gemma 4 E2B — small multilingual, chat only (no tool calling)"},
}

type Model struct {
	Tag   string
	Label string
}

// PickModel renders an interactive select prompt listing locally installed
// models first (tagged "[installed]") followed by curated recommendations.
// Returns "" if the user cancels.
func PickModel(installed []string) (string, error) {
	var choice string
	opts := make([]huh.Option[string], 0, len(installed)+len(Recommended))

	seen := make(map[string]bool, len(installed))
	for _, tag := range installed {
		opts = append(opts, huh.NewOption("[installed] "+tag, tag))
		seen[tag] = true
	}
	for _, m := range Recommended {
		if seen[m.Tag] {
			continue
		}
		opts = append(opts, huh.NewOption(m.Label, m.Tag))
	}

	err := huh.NewSelect[string]().
		Title("Pick a model to run locally").
		Description("Installed models appear first. Recommended picks support tool calling; entries marked 'chat only' will not work with opencode's agent loop.").
		Options(opts...).
		Value(&choice).
		Run()
	if err != nil {
		return "", err
	}
	return choice, nil
}
