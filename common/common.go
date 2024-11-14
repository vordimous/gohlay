package common

import (
	"fmt"
	"time"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	log "github.com/sirupsen/logrus"
)

// FmtKafkaKey generates a unique identifier
func FmtKafkaKey(offset kafka.Offset, delay int64) string {
	return fmt.Sprintf("%v-%d", offset, delay)
}

// ParseHeaders extracts relevant information from the message headers
func ParseHeaders(headers []kafka.Header) (deliveryTime int64, isDelivered bool, isDeliveredKey string, gohlayed bool) {
	for _, h := range headers {
		if h.Key == "GOHLAY" {
			timeString := string(h.Value)
			if t, err := time.Parse(time.UnixDate, timeString); err == nil {
				deliveryTime = t.Unix()
			} else {
				log.Error("Calculating time remaining from GOHLAY header | ", err)
			}
			gohlayed = true
		}
		if h.Key == "GOHLAY_FIN" {
			gohlayed = true
			isDeliveredKey = string(h.Value)
			isDelivered = true
		}
	}
	return
}

// GetAssignedPartitions finds the assigned partitions
func GetAssignedPartitions(c *kafka.Consumer, partitions []kafka.TopicPartition) ([]kafka.TopicPartition, error) {
	log.Info("GetAssignedPartitions | ", c, partitions)
	parts := make([]kafka.TopicPartition, len(partitions))
	var err error
	if false {
		limit := time.Now().Add(time.Duration(-5)*time.Minute).UnixNano() / int64(time.Millisecond)
		for i, tp := range partitions {
			offset, _ := kafka.NewOffset(limit)
			tp.Offset = offset
			log.Debug("GetAssignedPartitions offset query time | ", tp.Offset)
			parts[i] = tp
		}
		parts, err = c.OffsetsForTimes(parts, 10000)
	} else {
		for i, tp := range partitions {
			offset, _ := kafka.NewOffset(0)
			tp.Offset = offset
			log.Info("GetAssignedPartitions offset query value | ", tp.Offset)
			parts[i] = tp
		}
	}
	log.Info("Assigned partition(s) | ", parts)
	return parts, err
}
