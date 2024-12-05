package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vordimous/gohlay/configs"
	"github.com/vordimous/gohlay/pkg/find"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for gohlayed messages that are past the deadline",
	Long: `Scan the kafka topics for gohlayed messages and output the result as a JSON array.
The JSON array is a list of messages identified by "<offset>-<delivery time>".
Example Output:

["8-1731612338","9-1731614359","10-1731614360","11-1731614361"]
`,
	Run: func(cmd *cobra.Command, args []string) {
		configs.Load()
		results := []CheckResult{}
		for _, f := range find.CheckForDeliveries() {
			r := CheckResult{
				Topic:        f.TopicName(),
				DeliveryKeys: f.GohlayedSlice(),
			}
			results = append(results, r)
		}
		deliveriesJson, _ := json.Marshal(results)
		fmt.Println(string(deliveriesJson))

	},
}

type CheckResult struct {
	Topic        string   `json:"topic,omitempty"`
	DeliveryKeys []string `json:"delivery_keys,omitempty"`
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
