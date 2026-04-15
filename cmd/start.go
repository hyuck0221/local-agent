package cmd

import (
	"fmt"

	"github.com/hyuck0221/local-agent/internal/ollama"
	"github.com/hyuck0221/local-agent/internal/opencode"
	"github.com/hyuck0221/local-agent/internal/tui"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [model]",
	Short: "Start a local LLM and wire it into opencode",
	Long:  "Ensures Ollama is installed and running, pulls the requested model, and registers it as a provider in opencode's config.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !ollama.Installed() {
			fmt.Println("Ollama not found. Installing...")
			if err := ollama.Install(); err != nil {
				return fmt.Errorf("install ollama: %w", err)
			}
		}

		if err := ollama.StartServer(); err != nil {
			return fmt.Errorf("start ollama server: %w", err)
		}
		fmt.Println("Ollama server ready at", ollama.DefaultHost)

		var model string
		if len(args) == 1 {
			model = args[0]
		} else {
			picked, err := tui.PickModel()
			if err != nil {
				return err
			}
			if picked == "" {
				return fmt.Errorf("no model selected")
			}
			model = picked
		}

		if err := ollama.Pull(model); err != nil {
			return fmt.Errorf("pull %s: %w", model, err)
		}

		path, err := opencode.RegisterModel(model, ollama.BaseURL)
		if err != nil {
			return fmt.Errorf("update opencode config: %w", err)
		}

		fmt.Printf("\n✓ %s is running at %s\n", model, ollama.BaseURL)
		fmt.Printf("✓ Registered in %s as provider \"%s\"\n", path, opencode.ProviderID)
		fmt.Printf("\nRun `opencode` and choose model `%s/%s`.\n", opencode.ProviderID, model)
		return nil
	},
}
