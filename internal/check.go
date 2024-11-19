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
	ScanAll(indexMsg)
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
			log.Debugf("Message time remaining: %v", deliveryTime - deadline)
			if deliveryTime < deadline {
				messageId := common.FmtMessageId(msg.TopicPartition.Offset, deliveryTime)
				isDelivered[messageId] = false // set key to be delivered with a delivery value of false
				log.Debugf("Setting message for delivery: %d %d-%s %s", msg.TopicPartition.Partition, msg.TopicPartition.Offset, msg.Key, messageId)
			} else {
				log.Debugf("Message not ready for delivery: %d-%s", msg.TopicPartition.Offset, msg.Key)
			}
		} else if delivered {
			delete(isDelivered, deliveredKey)
			log.Debugf("Message is already delivered: %d-%s %s", msg.TopicPartition.Offset, msg.Key, deliveredKey)
		}
	} else {
		log.Debugf("Messaged is not gohlayed: %d-%s", msg.TopicPartition.Offset, msg.Key)
	}
}
