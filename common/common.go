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
func FmtKafkaGroup(topic string, deadline int64) string {
	return fmt.Sprintf("%s:%d", topic, viper.GetInt64("deadline"))
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

// GetAssignedPartitions finds the assigned partitions
func GetAssignedPartitions(c *kafka.Consumer, partitions []kafka.TopicPartition) ([]kafka.TopicPartition, error) {
	parts := make([]kafka.TopicPartition, len(partitions))
	var err error
	if false {
		limit := time.Now().Add(time.Duration(-5)*time.Minute).UnixNano() / int64(time.Millisecond)
		for i, tp := range partitions {
			offset, _ := kafka.NewOffset(limit)
			tp.Offset = offset
			log.Debugf("GetAssignedPartitions offset query time %v", tp.Offset)
			parts[i] = tp
		}
		parts, err = c.OffsetsForTimes(parts, 10000)
	} else {
		for i, tp := range partitions {
			offset, _ := kafka.NewOffset(0)
			tp.Offset = offset
			log.Infof("GetAssignedPartitions offset query value %v", tp.Offset)
			parts[i] = tp
		}
	}
	log.Infof("Assigned partition(s) %v", parts)
	return parts, err
}
