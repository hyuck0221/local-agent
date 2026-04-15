package cmd

import (
	"fmt"

	"github.com/hyuck0221/local-agent/internal/ollama"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Ollama if it is not already present",
	RunE: func(cmd *cobra.Command, args []string) error {
		if ollama.Installed() {
			fmt.Println("Ollama already installed.")
			return nil
		}
		fmt.Println("Installing Ollama...")
		if err := ollama.Install(); err != nil {
			return err
		}
		fmt.Println("Ollama installed. Run `local-agent start` next.")
		return nil
	},
}
