package tui

import "github.com/charmbracelet/huh"

// Recommended is a curated list of coding-friendly Ollama models shown in the
// picker when the user does not pass an explicit model argument.
var Recommended = []Model{
	{Tag: "qwen2.5-coder:7b", Label: "Qwen2.5 Coder 7B — strong small coder, ~5GB"},
	{Tag: "qwen2.5-coder:14b", Label: "Qwen2.5 Coder 14B — best quality if you have 16GB+ RAM"},
	{Tag: "llama3.2:3b", Label: "Llama 3.2 3B — fastest, ~2GB, light laptops"},
	{Tag: "deepseek-coder-v2:16b", Label: "DeepSeek Coder V2 16B — high quality, needs ~10GB"},
	{Tag: "gemma3:4b", Label: "Gemma 3 4B — balanced general-purpose"},
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
		Options(opts...).
		Value(&choice).
		Run()
	if err != nil {
		return "", err
	}
	return choice, nil
}
