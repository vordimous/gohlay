package config

import (
	"strings"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	_, kafkaVersion := kafka.LibraryVersion()
	log.Debugf("Kafka Lib Version: %s", kafkaVersion)
}

func getBootrapServersString() (string) {
	return strings.Join(viper.GetStringSlice("bootstrap_servers"), ",")
}

// GetConsumer creates the kafka.ConfigMap based on the configured settings for a Consumer
func GetConsumer() (config *kafka.ConfigMap, topics []string) {
	config = &kafka.ConfigMap{
		"bootstrap.servers":               getBootrapServersString(),
		"go.application.rebalance.enable": true, // delegate Assign() responsibility to app
		"session.timeout.ms":              6000,
		"enable.partition.eof":            true,
		"default.topic.config":            kafka.ConfigMap{"auto.offset.reset": "earliest"},
	}
	topics = viper.GetStringSlice("topics")
	log.Debugf("Kafka consumer config: %+v", config)
	return
}

// GetProducer creates the kafka.ConfigMap based on the configured settings for a Producer
func GetProducer() (config *kafka.ConfigMap) {
	config = &kafka.ConfigMap{
		"bootstrap.servers": getBootrapServersString(),
	}
	log.Debugf("Kafka producer config: %+v", config)
	return
}
