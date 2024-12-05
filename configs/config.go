package configs

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Load will perform all of the steps necessary to capture any user configured settings
func Load() {
	readEnv()
	readConfig()
	readKafkaConfig()
	setupLogging()
	printConfig()
}

// setupLogging setup the configured logging options
func setupLogging() {
	log.SetOutput(os.Stdout)

	if viper.GetBool("json") {
		log.SetFormatter(&log.JSONFormatter{})
	}

	logLevel := log.WarnLevel
	if viper.GetBool("verbose") {
		logLevel = log.InfoLevel
	}
	if viper.GetBool("debug") {
		logLevel = log.DebugLevel
	}
	if viper.GetBool("silent") {
		logLevel = log.PanicLevel
	}
	log.Debugf("Setting log Level %v", logLevel)
	log.SetLevel(logLevel)
}

// readConfig reads user settings from an optional config file
func readConfig() {
	viper.SetConfigName("gohlay")
	viper.AddConfigPath("/etc/gohlay/")
	viper.AddConfigPath("$HOME/.gohlay")
	viper.AddConfigPath(viper.GetString("config-dir"))
	err := viper.ReadInConfig()
	if err != nil {
		log.Infof("Optional %v", err)
	}
}

// readConfig reads user settings from an optional config file
func readEnv() {
	viper.SetEnvPrefix("gohlay")
	viper.AutomaticEnv()
}

// printConfig print all user defined settings
func printConfig() {
	for key, value := range viper.GetViper().AllSettings() {
		log.WithFields(log.Fields{
			key: value,
		}).Info("Config")
	}
}
