package kafkautil

import (
	"fmt"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/spf13/viper"
)

// FmtMessageId formats a gohlayed message id
func FmtMessageId(offset kafka.Offset, delay int64) string {
	return fmt.Sprintf("%v-%d", offset, delay)
}

// FmtKafkaGroup formats a group id
func FmtKafkaGroup(reason string, topic string) string {
	return fmt.Sprintf("group.id=gohlay_%s:%s:%d", reason, topic, viper.GetInt64("deadline"))
}
