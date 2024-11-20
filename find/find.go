package find

import (
	"maps"
	"slices"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vordimous/gohlay/internal"
	"github.com/vordimous/gohlay/kafkautil"
)

type Finder struct {
	reason string
	gohlayedMessages map[string]bool
}

// CheckForDeliveries will scan the topic and build a map of messages to be delivered
func CheckForDeliveries() (*Finder) {
	f := new(Finder)
	f.reason = "indexing"
	internal.ScanAll(f)
	return f
}


// GetGohlayedSlice returns the map of gohlayed message keys
func (f *Finder) GetGohlayed() map[string]bool {
	return f.gohlayedMessages
}

// GetGohlayedSlice creates an array of strings from the gohlayed map keys
func (f *Finder) GetGohlayedSlice() []string {
	if d := slices.Collect(maps.Keys(f.GetGohlayed())); d != nil {
		return d
	}
	return []string{}
}

func (f *Finder) GetReason() string {
	return f.reason
}

func (f *Finder) HandleMessage(msg *kafka.Message) {
	deadline := viper.GetInt64("deadline")
	deliveryTime, delivered, deliveredKey, hasHeader := kafkautil.ParseHeaders(msg.Headers)
	if hasHeader {
		if !delivered && deliveryTime != 0 {
			log.Debugf("Message time remaining: %v", deliveryTime-deadline)
			if deliveryTime < deadline {
				messageId := kafkautil.FmtMessageId(msg.TopicPartition.Offset, deliveryTime)
				f.gohlayedMessages[messageId] = false // set key to be delivered with a delivery value of false
				log.Debugf("Setting message for delivery: %d %d-%s %s", msg.TopicPartition.Partition, msg.TopicPartition.Offset, msg.Key, messageId)
			} else {
				log.Debugf("Message not ready for delivery: %d-%s", msg.TopicPartition.Offset, msg.Key)
			}
		} else if delivered {
			delete(f.gohlayedMessages, deliveredKey)
			log.Debugf("Message is already delivered: %d-%s %s", msg.TopicPartition.Offset, msg.Key, deliveredKey)
		}
	} else {
		log.Debugf("Messaged is not gohlayed: %d-%s", msg.TopicPartition.Offset, msg.Key)
	}
}
