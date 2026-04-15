package tui

import "github.com/charmbracelet/huh"

// Recommended is a curated list of Ollama models that support OpenAI-style
// tool calling, which opencode requires for its agent loop. Models that cannot
// emit structured tool calls (Gemma, DeepSeek-R1, Phi3, …) are intentionally
// excluded — opencode fails immediately against them.
var Recommended = []Model{
	{Tag: "qwen2.5-coder:7b", Label: "Qwen2.5 Coder 7B — best small coder, ~5GB (default pick)"},
	{Tag: "qwen2.5-coder:14b", Label: "Qwen2.5 Coder 14B — higher quality, ~9GB, 16GB+ RAM"},
	{Tag: "qwen3:8b", Label: "Qwen 3 8B — modern general-purpose, ~5GB"},
	{Tag: "llama3.1:8b", Label: "Llama 3.1 8B — Meta baseline with strong tools, ~5GB"},
	{Tag: "llama3.3:70b", Label: "Llama 3.3 70B — top quality, ~40GB, workstation only"},
	{Tag: "mistral-nemo:12b", Label: "Mistral Nemo 12B — fast with solid tool calling, ~7GB"},
	{Tag: "deepseek-coder-v2:16b", Label: "DeepSeek Coder V2 16B — powerful MoE coder, ~9GB"},
	{Tag: "command-r:35b", Label: "Command R 35B — RAG/tools specialist, ~20GB"},
}

type Model struct {
	Tag   string
	Label string
}

// PickModel renders an interactive select prompt and returns the chosen tag.
// Returns "" if the user cancels.
func PickModel() (string, error) {
	var choice string
	opts := make([]huh.Option[string], 0, len(Recommended))
	for _, m := range Recommended {
		opts = append(opts, huh.NewOption(m.Label, m.Tag))
	}
	err := huh.NewSelect[string]().
		Title("Pick a model to run locally").
		Description("All listed models support tool calling, which opencode requires.").
		Options(opts...).
		Value(&choice).
		Run()
	if err != nil {
		return "", err
	}
	return choice, nil
}
