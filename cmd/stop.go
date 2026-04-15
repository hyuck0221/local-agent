package cmd

import (
	"fmt"

	"github.com/hyuck0221/local-agent/internal/ollama"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the local Ollama server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !ollama.Serving() {
			fmt.Println("Ollama server is not running.")
			return nil
		}
		if err := ollama.Stop(); err != nil {
			return err
		}
		fmt.Println("Ollama server stopped.")
		return nil
	},
}
