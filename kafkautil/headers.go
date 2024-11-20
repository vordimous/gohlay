package kafkautil

import (
	"time"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	log "github.com/sirupsen/logrus"
	"github.com/vordimous/gohlay/config"
)

// ParseHeaders extracts relevant information from the message headers
func ParseHeaders(headers []kafka.Header) (deliveryTime int64, isDelivered bool, isDeliveredKey string, gohlayed bool) {
	for _, h := range headers {
		switch h.Key {
		case config.GetHeaderOverride("GOHLAY"):
			timeString := string(h.Value)
			if t, err := time.Parse(time.UnixDate, timeString); err == nil {
				deliveryTime = t.UnixMilli()
			} else {
				log.Errorf("Calculating time remaining from GOHLAY header %v", err)
			}
			gohlayed = true
		case config.GetHeaderOverride("GOHLAY_DELIVERED"):
			gohlayed = true
			isDeliveredKey = string(h.Value)
			isDelivered = true
		}
	}
	return
}
