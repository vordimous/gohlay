package common

import (
	"fmt"
	"time"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	log "github.com/sirupsen/logrus"
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

// ParseHeaders extracts relevant information from the message headers
func ParseHeaders(headers []kafka.Header) (deliveryTime int64, isDelivered bool, isDeliveredKey string, gohlayed bool) {
	for _, h := range headers {
		if h.Key == "GOHLAY" {
			timeString := string(h.Value)
			if t, err := time.Parse(time.UnixDate, timeString); err == nil {
				deliveryTime = t.UnixMilli()
			} else {
				log.Errorf("Calculating time remaining from GOHLAY header %v", err)
			}
			gohlayed = true
		}
		if h.Key == "GOHLAY_DELIVERED" {
			gohlayed = true
			isDeliveredKey = string(h.Value)
			isDelivered = true
		}
	}
	return
}
