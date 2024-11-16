/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vordimous/gohlay/config"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gohlay",
	Short: "Replay messages on a Kafka topic after a deadline",
	Long: `Gohlay is a delayed delivery tool for producing messages onto
Kafka topics on a schedule set by a Kafka message header.`,
	Run: func(cmd *cobra.Command, args []string) {
		config.Load()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// dev options
var cfgFileDir string
var Silent bool
var Verbose bool
var Debug bool
var JsonOut bool

// Gohlay options
var Deadline int64

// kafka options
var bootstrap_servers []string
var Topics []string

func init() {

	// Dev
	rootCmd.PersistentFlags().StringVar(&cfgFileDir, "config_dir", ".", "config file directory.")
	viper.BindPFlag("config_dir", rootCmd.PersistentFlags().Lookup("config_dir"))
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Display more verbose output in console output.")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	rootCmd.PersistentFlags().BoolVar(&Debug, "debug", false, "Display debugging output in the console.")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	rootCmd.PersistentFlags().BoolVar(&Silent, "silent", false, "Don't display output in the console.")
	viper.BindPFlag("silent", rootCmd.PersistentFlags().Lookup("silent"))
	rootCmd.PersistentFlags().BoolVar(&JsonOut, "json", false, "Display output in the console as JSON.")
	viper.BindPFlag("json", rootCmd.PersistentFlags().Lookup("json"))

	// Gohlay
	rootCmd.PersistentFlags().Int64VarP(&Deadline, "deadline", "d", time.Now().UnixMilli(), "Sets the delivery deadline. (Format: Unix Timestamp)")
	viper.BindPFlag("deadline", rootCmd.PersistentFlags().Lookup("deadline"))

	// Kafka
	rootCmd.PersistentFlags().StringArrayVarP(&bootstrap_servers, "bootstrap_servers", "b", []string{"localhost:9092"}, "Sets the \"bootstrap.servers\" parameter in the kafka.ConfigMap")
	viper.BindPFlag("bootstrap_servers", rootCmd.PersistentFlags().Lookup("bootstrap_servers"))
	rootCmd.PersistentFlags().StringArrayVarP(&Topics, "topics", "t", []string{"gohlay"}, "Sets the kafka topics to use")
	viper.BindPFlag("topics", rootCmd.PersistentFlags().Lookup("topics"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}


