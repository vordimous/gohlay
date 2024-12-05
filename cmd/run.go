package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vordimous/gohlay/configs"
	"github.com/vordimous/gohlay/pkg/deliver"
	"github.com/vordimous/gohlay/pkg/find"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Check for gohlayed messages and deliver them.",
	Long: `Perform a check on the configured topics and deliver any
messages that are past the deadline.`,
	Run: func(cmd *cobra.Command, args []string) {
		configs.Load()
		for _, f := range find.CheckForDeliveries() {
			deliver.HandleDeliveries(f.TopicName(), f.GohlayedMap())
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
