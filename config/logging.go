package config

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// SetupLogging setup the configured logging options
func SetupLogging() {
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
	log.Debug("Setting log Level | ", logLevel)
	log.SetLevel(logLevel)
}

// PrintConfig print all user defined settings
func PrintConfig() {
	for key, value := range viper.GetViper().AllSettings() {
		log.WithFields(log.Fields{
			key: value,
		}).Info("Config")
	}
}
