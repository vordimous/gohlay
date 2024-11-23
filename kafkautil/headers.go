package kafkautil

import (
	"time"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	log "github.com/sirupsen/logrus"
	"github.com/vordimous/gohlay/config"
)

type GohlayedMeta struct {
	Gohlayed bool
	DeliveryTime int64
	Delivered bool
	DeliveryKey string
}

// ParseHeaders extracts relevant information from the message headers
func ParseHeaders(headers []kafka.Header) (GohlayedMeta) {
	meta := GohlayedMeta{}
	for _, h := range headers {
		switch h.Key {
		case config.GetHeaderOverride("GOHLAY"):
			if t, err := time.Parse(time.UnixDate, string(h.Value)); err == nil {
				meta.DeliveryTime = t.UnixMilli()
			} else {
				log.Errorf("Calculating time remaining from GOHLAY header %v", err)
			}
			meta.Gohlayed = true
		case config.GetHeaderOverride("GOHLAY_DELIVERED"):
			meta.Gohlayed = true
			meta.DeliveryKey = string(h.Value)
			meta.Delivered = true
		}
	}
	return meta
}
