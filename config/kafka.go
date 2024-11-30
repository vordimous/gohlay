package config

import (
	"strings"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	overrideMap          map[string]string
	bootrapServersString string
	topics               []string
	kafkaSettingsMap     kafka.ConfigMap
)

func init() {
	_, kafkaVersion := kafka.LibraryVersion()
	log.Debugf("Kafka Lib Version: %s", kafkaVersion)
	overrideMap = map[string]string{}
	bootrapServersString = ""
	topics = []string{}
	kafkaSettingsMap = kafka.ConfigMap{}
	Consumer()
}

func readKafkaConfig() {

	topics = viper.GetStringSlice("topics")

	for _, kv := range viper.GetStringSlice("override_headers") {
		p := strings.Split(kv, "=")
		if len(p) == 2 {
			overrideMap[p[0]] = p[1]
		}
	}

	for _, settingKV := range viper.GetStringSlice("kafka_properties") {
		p := strings.Split(settingKV, "=")
		if len(p) == 2 {
			kafkaSettingsMap[p[0]] = p[1]
		}
	}

	bootrapServersString = strings.Join(viper.GetStringSlice("bootstrap_servers"), ",")
	kafkaSettingsMap["bootstrap.servers"] = bootrapServersString
}

// HeaderOverride checks for a configured name override for a default header name
func HeaderOverride(header string) string {
	if overrideMap[header] != "" {
		return overrideMap[header]
	}
	return header
}

// Topics returns the list of configured topics
func Topics() []string {
	return topics
}

// Consumer creates the kafka.ConfigMap based on the configured settings for a Consumer
func Consumer() (config kafka.ConfigMap) {
	config = kafka.ConfigMap{}
	for k, v := range kafkaSettingsMap {
		config[k] = v
	}
	config["go.application.rebalance.enable"] = true // delegate Assign() responsibility to app
	config["session.timeout.ms"] = 6000
	config["enable.partition.eof"] = true
	config["default.topic.config"] = kafka.ConfigMap{"auto.offset.reset": "earliest"}
	log.Debugf("Kafka consumer config: %+v", config)
	return
}

// Producer creates the kafka.ConfigMap based on the configured settings for a Producer
func Producer() (config kafka.ConfigMap) {
	config = kafka.ConfigMap{}
	for k, v := range kafkaSettingsMap {
		config[k] = v
	}
	log.Debugf("Kafka producer config: %+v", config)
	return
}
