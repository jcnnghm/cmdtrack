package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var command = &Command{}

var trackCmd = &cobra.Command{
	Use:   "track",
	Short: "Track the command provided.",
	Run: func(cmd *cobra.Command, args []string) {
		command.command = strings.Join(args, " ")
		fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
		fmt.Println(command.command)
		fmt.Println(command.workingDir)
	},
}

func init() {
	trackCmd.PersistentFlags().StringVarP(&command.workingDir, "workdir", "d", "", "Working directory command was executed from")
	cmdTracker.AddCommand(trackCmd)
}
