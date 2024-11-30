package find

import (
	"fmt"
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
	topic            string
	gohlayedMessages map[string]bool
}

// GohlayedMap returns the map of gohlayed message keys
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

// TopicName is the name of the topic to use when finding
func (f *Finder) TopicName() string {
	return f.topic
}

// GroupName is a human readable name of the purpose for the message handler
func (f *Finder) GroupName() string {
	return "finding"
}

// HandleMessage index any gohlayed message with a delivery time passed the deadline
func (f *Finder) HandleMessage(msg *kafka.Message) string {
	gohlayedMeta, err := kafkautil.ParseHeaders(msg.Headers)
	if err != nil {
		return fmt.Sprintf("could't parse headers: %v", err)
	}
	
	if !gohlayedMeta.Gohlayed {
		return fmt.Sprintf("message is not Gohlayed: %v %d %d", msg.TopicPartition.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset)

	}
	deliveryKey := kafkautil.FmtDeliveryKey(msg.TopicPartition.Offset, gohlayedMeta.DeliveryTime)
	if gohlayedMeta.Delivered || gohlayedMeta.DeliveryKey != "" || f.gohlayedMessages[gohlayedMeta.DeliveryKey] {
		f.gohlayedMessages[gohlayedMeta.DeliveryKey] = true
		return fmt.Sprintf("message already delivered: %d-%s %s", msg.TopicPartition.Offset, msg.Key, gohlayedMeta.DeliveryKey)
	}

	deadline := viper.GetInt64("deadline")
	deliveryTime := gohlayedMeta.DeliveryTime
	if deliveryTime > deadline {
		return fmt.Sprintf("message not ready for delivery: %d-%s %d", msg.TopicPartition.Offset, msg.Key, deliveryTime-deadline)
	}

	f.gohlayedMessages[deliveryKey] = false // set key to be delivered with a delivery value of false
	log.Infof("Setting message for delivery: %d %d-%s %s", msg.TopicPartition.Partition, msg.TopicPartition.Offset, msg.Key, deliveryKey)
	return ""
}

func (f *Finder) findGohlayed() {
	f.gohlayedMessages = map[string]bool{}
	internal.ScanTopic(f)
}
