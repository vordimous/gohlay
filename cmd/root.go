package cmd

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
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
	if err := viper.BindPFlag("config_dir", rootCmd.PersistentFlags().Lookup("config_dir")); err != nil { log.Error(err) }
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Display more verbose output in console output.")
	if err := viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose")); err != nil { log.Error(err) }
	rootCmd.PersistentFlags().BoolVar(&Debug, "debug", false, "Display debugging output in the console.")
	if err := viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug")); err != nil { log.Error(err) }
	rootCmd.PersistentFlags().BoolVar(&Silent, "silent", false, "Don't display output in the console.")
	if err := viper.BindPFlag("silent", rootCmd.PersistentFlags().Lookup("silent")); err != nil { log.Error(err) }
	rootCmd.PersistentFlags().BoolVar(&JsonOut, "json", false, "Display output in the console as JSON.")
	if err := viper.BindPFlag("json", rootCmd.PersistentFlags().Lookup("json")); err != nil { log.Error(err) }

	// Gohlay
	rootCmd.PersistentFlags().Int64VarP(&Deadline, "deadline", "d", time.Now().UnixMilli(), "Sets the delivery deadline. (Format: Unix Timestamp)")
	if err := viper.BindPFlag("deadline", rootCmd.PersistentFlags().Lookup("deadline")); err != nil { log.Error(err) }

	// Kafka
	rootCmd.PersistentFlags().StringArrayVarP(&bootstrap_servers, "bootstrap_servers", "b", []string{"localhost:9092"}, "Sets the \"bootstrap.servers\" parameter in the kafka.ConfigMap")
	if err := viper.BindPFlag("bootstrap_servers", rootCmd.PersistentFlags().Lookup("bootstrap_servers")); err != nil { log.Error(err) }
	rootCmd.PersistentFlags().StringArrayVarP(&Topics, "topics", "t", []string{"gohlay"}, "Sets the kafka topics to use")
	if err := viper.BindPFlag("topics", rootCmd.PersistentFlags().Lookup("topics")); err != nil { log.Error(err) }
}


