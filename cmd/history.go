package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var historyCmd = &cobra.Command{
	Use:     "history",
	Short:   "Output the recent combined history",
	Example: "cmdtrack history --url=http://localhost:8080/",
	Run: func(cmd *cobra.Command, args []string) {
		commands, err := FetchCommands(cmdtrackURL, isVerbose(cmd))
		if err != nil {
			fmt.Println("History fetch error: %v", err.Error())
			os.Exit(-1)
		}

		for index := range commands {
			fmt.Printf("%8d  %v\n", index+1, commands[len(commands)-index-1].Command)
		}
	},
}

func init() {
	historyCmd.PersistentFlags().StringVar(&cmdtrackURL, "url", "https://cmdtrack-1127.appspot.com/", "URL for the cmdtrack server")
	historyCmd.PersistentFlags().BoolP("verbose", "v", false, "Make command verbose")
	cmdTrack.AddCommand(historyCmd)
}

func isVerbose(cmd *cobra.Command) bool {
	verbose, err := cmd.PersistentFlags().GetBool("verbose")
	if err != nil {
		panic(err)
	}
	return verbose
}
