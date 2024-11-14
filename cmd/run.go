/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vordimous/gohlay/config"
	"github.com/vordimous/gohlay/internal"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Check for gohlayed messages and deliver them.",
	Long: `Perform a check on the configured topics and deliver any
messages that are past the deadline.`,
	Run: func(cmd *cobra.Command, args []string) {
		config.SetupLogging()
		config.PrintConfig()
		internal.Run()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
