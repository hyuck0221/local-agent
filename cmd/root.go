package cmd

import "github.com/spf13/cobra"

var version = "dev"

var rootCmd = &cobra.Command{
	Use:           "local-agent",
	Short:         "Run a local LLM and wire it into opencode with one command",
	Long:          "local-agent bootstraps Ollama, picks a model, and registers it as an OpenAI-compatible provider in opencode's config.",
	Version:       version,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(modelsCmd)
}
