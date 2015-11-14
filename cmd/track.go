package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var command = &Command{}

var trackCmd = &cobra.Command{
	Use:     "track",
	Short:   "Track the command provided.",
	Example: "cmdtrack track --workdir=~ -- ls",
	Run: func(cmd *cobra.Command, args []string) {
		command.command = strings.Join(args, " ")
		fmt.Println(command.command)
		fmt.Println(command.workingDir)
	},
}

func init() {
	trackCmd.PersistentFlags().StringVarP(&command.workingDir, "workdir", "d", "", "Working directory command was executed from")
	cmdTrack.AddCommand(trackCmd)
}
