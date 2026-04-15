package cmd

import (
	"fmt"

	"github.com/hyuck0221/local-agent/internal/ollama"
	"github.com/hyuck0221/local-agent/internal/tui"
	"github.com/spf13/cobra"
)

var modelsCmd = &cobra.Command{
	Use:   "models",
	Short: "List locally installed models and recommended picks",
	RunE: func(cmd *cobra.Command, args []string) error {
		installed, err := ollama.List()
		if err != nil {
			fmt.Println("could not query ollama:", err)
		} else if len(installed) == 0 {
			fmt.Println("No models installed locally.")
		} else {
			fmt.Println("Installed:")
			for _, m := range installed {
				fmt.Println("  -", m)
			}
		}
		fmt.Println("\nRecommended:")
		for _, m := range tui.Recommended {
			fmt.Printf("  - %s\n      %s\n", m.Tag, m.Label)
		}
		return nil
	},
}
