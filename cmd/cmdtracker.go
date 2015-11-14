package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var cmdTracker = &cobra.Command{
	Use:   "cmdtracker",
	Short: "cmdtracker stores global command-line history.",
}

// Execute causes the CmdTracker to process the command-line args and run
func Execute() {
	if err := cmdTracker.Execute(); err != nil {
		// Error is already reported by cobra
		os.Exit(-1)
	}
}
