package main

import (
	"fmt"
	"os"

	"github.com/ankitpokhrel/jira-cli/internal/cmd"
)

// main is the entry point for the jira-cli application.
// It delegates execution to the root command defined in the cmd package.
func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
