package internal

import (
	"maps"
	"slices"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vordimous/gohlay/common"
)

var (
	isDelivered map[string]bool
)

func init() {
	isDelivered = map[string]bool{}
}

// CheckForDeliveries will scan the topic and build a map of messages to be delivered
func CheckForDeliveries() {
	ScanTopic(indexMsg)
}

// GetDeliveries creates an array of strings from the deliveries map keys
func GetDeliveries() ([]string) {
	if d := slices.Collect(maps.Keys(isDelivered)); d != nil {
		return d
	}
	return []string{}
}

func indexMsg(msg *kafka.Message) {
	deadline := viper.GetInt64("deadline")
	deliveryTime, fin, finKey, hasHeader := common.ParseHeaders(msg.Headers)
	if hasHeader {
		if !fin && deliveryTime != 0 {
			log.Debug("Message time remaining | ", deliveryTime - deadline)
			if deliveryTime < deadline {
				isDelivered[common.FmtKafkaKey(msg.TopicPartition.Offset, deliveryTime)] = false // set key to be delivered with a delivery value of false
				log.Debug("Setting message for delivery | ", common.FmtKafkaKey(msg.TopicPartition.Offset, deliveryTime))
			} else {
				log.Debug("Message not ready for delivery | ", msg.TopicPartition.Offset)
			}
		} else if fin {
			delete(isDelivered, finKey)
			log.Debug("Message not ready for delivery | ", finKey)
		}
	} else {
		log.Debug("Messaged is not gohlayed | ", msg.TopicPartition.Offset)
	}
}
