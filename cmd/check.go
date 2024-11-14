/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vordimous/gohlay/config"
	"github.com/vordimous/gohlay/internal"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for golayed messages that are past the deadline",
	Long: `Scan the kafka topics for golayed messages and output the result as a JSON array.
The JSON array is a list of messages identified by "<offset>-<delivery time>".
Example Output:

["8-1731612338","9-1731614359","10-1731614360","11-1731614361"]
`,
	Run: func(cmd *cobra.Command, args []string) {
		config.SetupLogging()
		config.PrintConfig()
		internal.CheckForDeliveries()
		deliveriesJson, _ := json.Marshal(internal.GetDeliveries())
		fmt.Println(string(deliveriesJson))

	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
