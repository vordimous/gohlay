/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vordimous/gohlay/config"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gohlay",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		config.SetupLogging()
		config.PrintConfig()
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
var Silent bool
var Verbose bool
var Debug bool
var JsonOut bool

// kafka options
var BootstrapServers string
var Topics []string

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gohlay.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Display more verbose output in console output.")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	rootCmd.PersistentFlags().BoolVar(&Debug, "debug", false, "Display debugging output in the console.")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	rootCmd.PersistentFlags().BoolVar(&Silent, "silent", false, "Don't display output in the console.")
	viper.BindPFlag("silent", rootCmd.PersistentFlags().Lookup("silent"))
	rootCmd.PersistentFlags().BoolVar(&JsonOut, "json", false, "Display output in the console as JSON.")
	viper.BindPFlag("json", rootCmd.PersistentFlags().Lookup("json"))

	// kafka
	rootCmd.PersistentFlags().StringVarP(&BootstrapServers, "bootstrap_servers", "b", "localhost:9092", "Sets the \"bootstrap.servers\" parameter in the kafka.ConfigMap")
	viper.BindPFlag("bootstrap_servers", rootCmd.PersistentFlags().Lookup("bootstrap_servers"))
	rootCmd.PersistentFlags().StringArrayVarP(&Topics, "topics", "t", []string{"gohlay"}, "Sets the kafka topics to use")
	viper.BindPFlag("topics", rootCmd.PersistentFlags().Lookup("topics"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}


