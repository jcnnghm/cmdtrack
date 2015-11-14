package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var cmdTrack = &cobra.Command{
	Use:   "cmdtrack",
	Short: "cmdtrack stores global command-line history.",
}

// Execute causes the CmdTracker to process the command-line args and run
func Execute() {
	if err := cmdTrack.Execute(); err != nil {
		// Error is already reported by cobra
		os.Exit(-1)
	}
}
