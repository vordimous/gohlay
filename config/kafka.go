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

func GetHeaderOverride(header string) string {
	overrideMap := map[string]string{}
	for _, kv := range strings.Split(viper.GetString("override_headers"), ",") {
		p := strings.Split(kv, "=")
		if len(p) == 2 {
			overrideMap[p[0]] = p[1]
		}
	}
	if overrideMap[header] != "" {
		return overrideMap[header]
	}

	return header
}

func getBootrapServersString() string {
	return strings.Join(viper.GetStringSlice("bootstrap_servers"), ",")
}

// GetConsumer creates the kafka.ConfigMap based on the configured settings for a Consumer
func GetTopics() (topics []string) {
	topics = viper.GetStringSlice("topics")
	if len(topics) == 0 {
		log.Error("No topics defined")
	}
	return
}

// GetConsumer creates the kafka.ConfigMap based on the configured settings for a Consumer
func GetConsumer() (config *kafka.ConfigMap) {
	config = &kafka.ConfigMap{
		"bootstrap.servers":               getBootrapServersString(),
		"go.application.rebalance.enable": true, // delegate Assign() responsibility to app
		"session.timeout.ms":              6000,
		"enable.partition.eof":            true,
		"default.topic.config":            kafka.ConfigMap{"auto.offset.reset": "earliest"},
	}
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
