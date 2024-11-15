package internal

import (
	"fmt"
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
	deliveryTime, delivered, deliveredKey, hasHeader := common.ParseHeaders(msg.Headers)
	if hasHeader {
		if !delivered && deliveryTime != 0 {
			log.Debug("Message time remaining | ", deliveryTime - deadline)
			if deliveryTime < deadline {
				deliveryKey := common.FmtKafkaKey(msg.TopicPartition.Offset, deliveryTime)
				isDelivered[deliveryKey] = false // set key to be delivered with a delivery value of false
				log.Debug("Setting message for delivery | ", fmt.Sprintf("%d-%s %s", msg.TopicPartition.Offset, msg.Key, deliveryKey))
			} else {
				log.Debug("Message not ready for delivery | ", fmt.Sprintf("%d-%s", msg.TopicPartition.Offset, msg.Key) )
			}
		} else if delivered {
			delete(isDelivered, deliveredKey)
			log.Debug("Message is already delivered | ", fmt.Sprintf("%d-%s %s", msg.TopicPartition.Offset, msg.Key, deliveredKey))
		}
	} else {
		log.Debug("Messaged is not gohlayed | ", fmt.Sprintf("%d-%s", msg.TopicPartition.Offset, msg.Key) )
	}
}
