package config

import (
	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	_, kafkaVersion := kafka.LibraryVersion()
	log.Debug("Kafka Lib Version | ", kafkaVersion)
}

// GetConsumer creates the kafka.ConfigMap based on the configured settings for a Consumer
func GetConsumer() (config *kafka.ConfigMap, topics []string) {
	config = &kafka.ConfigMap{
		"bootstrap.servers":               viper.GetString("bootstrap_servers"),
		"go.application.rebalance.enable": true, // delegate Assign() responsibility to app
		"session.timeout.ms":              6000,
		"enable.partition.eof":            true,
		"default.topic.config":            kafka.ConfigMap{"auto.offset.reset": "earliest"},
	}
	topics = viper.GetStringSlice("topics")

	return
}

// GetProducer creates the kafka.ConfigMap based on the configured settings for a Producer
func GetProducer() (config *kafka.ConfigMap) {
	config = &kafka.ConfigMap{
		"bootstrap.servers": viper.GetString("bootstrap_servers"),
	}

	return
}
