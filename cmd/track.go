package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var command = &Command{}
var cmdtrackURL string

var trackCmd = &cobra.Command{
	Use:     "track",
	Short:   "Track the command provided.",
	Example: "cmdtrack track --workdir=~ --command=ls",
	Run: func(cmd *cobra.Command, args []string) {
		if !command.IsValid() {
			fmt.Println("Command, WorkingDir, and Hostname are all required")
			os.Exit(-1)
		}

		if command.Deduplicate() {
			if isVerbose(cmd) {
				fmt.Println("Skipping tracking duplicate command")
			}
			return
		}

		err := command.Send(cmdtrackURL)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}
	},
}

func init() {
	trackCmd.PersistentFlags().StringVarP(&command.WorkingDir, "workdir", "d", "", "Working directory command was executed from")
	trackCmd.PersistentFlags().StringVarP(&command.Command, "command", "c", "", "Command that was executed")
	trackCmd.PersistentFlags().StringVarP(&command.Hostname, "hostname", "n", "", "Hostname the command was executed on")
	trackCmd.PersistentFlags().Int64Var(&command.Timestamp, "timestamp", 0, "Timestamp to use for this command - sets to current time if 0")
	trackCmd.PersistentFlags().StringVar(&cmdtrackURL, "url", "https://cmdtrack-1127.appspot.com/", "URL for the cmdtrack server")
	trackCmd.PersistentFlags().BoolP("verbose", "v", false, "Make command verbose")
	cmdTrack.AddCommand(trackCmd)
}
