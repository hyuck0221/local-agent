package main

import (
	"fmt"
	"os"

	"github.com/hyuck0221/local-agent/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
