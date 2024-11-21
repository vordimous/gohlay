package find

import (
	"maps"
	"slices"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vordimous/gohlay/config"
	"github.com/vordimous/gohlay/internal"
	"github.com/vordimous/gohlay/kafkautil"
)

// CheckForDeliveries will scan the topic and build a map of messages to be delivered
func CheckForDeliveries() (found []*Finder) {
	for _, topic := range config.GetTopics() {
		f := &Finder{
			topic: topic,
		}
		f.findGohlayed()
		found = append(found, f)
	}
	return
}

type Finder struct {
	topic string
	gohlayedMessages map[string]bool
}

// GohlayedSlice returns the map of gohlayed message keys
func (f *Finder) GohlayedMap() map[string]bool {
	return f.gohlayedMessages
}

// GohlayedSlice creates an array of strings from the gohlayed map keys
func (f *Finder) GohlayedSlice() []string {
	if d := slices.Collect(maps.Keys(f.GohlayedMap())); d != nil {
		return d
	}
	return []string{}
}

// GroupName is a human readable name of the purpose for the message handler
func (f *Finder) TopicName() string {
	return f.topic
}

// GroupName is a human readable name of the purpose for the message handler
func (f *Finder) GroupName() string {
	return "finding"
}

// HandleMessage index any gohlayed message with a delivery time passed the deadline
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

// GroupName is a human readable name of the purpose for the message handler
func (f *Finder) findGohlayed() {
	f.gohlayedMessages = map[string]bool{}
	internal.ScanTopic(f)
}
