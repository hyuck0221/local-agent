package cmd

import (
	"fmt"

	"github.com/hyuck0221/local-agent/internal/ollama"
	"github.com/hyuck0221/local-agent/internal/opencode"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show local-agent state (Ollama + opencode wiring)",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("ollama installed : %v\n", ollama.Installed())
		fmt.Printf("ollama serving   : %v (%s)\n", ollama.Reachable(), ollama.DefaultHost)

		registered, err := opencode.Registered()
		if err != nil {
			return err
		}
		if registered == "" {
			fmt.Println("opencode provider: not registered (run `local-agent start`)")
		} else {
			fmt.Printf("opencode provider: %s/%s\n", opencode.ProviderID, registered)
		}
		return nil
	},
}
