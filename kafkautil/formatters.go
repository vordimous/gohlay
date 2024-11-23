package kafkautil

import (
	"fmt"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/spf13/viper"
)

// FmtDeliveryKey formats a gohlayed message id
func FmtDeliveryKey(offset kafka.Offset, delay int64) string {
	return fmt.Sprintf("%v-%d", offset, delay)
}

// FmtKafkaGroup formats a group id
func FmtKafkaGroup(groupName string, topic string) string {
	return fmt.Sprintf("group.id=gohlay_%s:%s:%d", groupName, topic, viper.GetInt64("deadline"))
}
