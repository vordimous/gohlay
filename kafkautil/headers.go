package kafkautil

import (
	"fmt"
	"time"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/vordimous/gohlay/config"
)

type GohlayedMeta struct {
	Gohlayed bool
	DeliveryTime int64
	Delivered bool
	DeliveryKey string
}

// ParseHeaders extracts relevant information from the message headers
func ParseHeaders(headers []kafka.Header) (GohlayedMeta, error) {
	meta := GohlayedMeta{}
	for _, h := range headers {
		switch h.Key {
		case config.HeaderOverride("GOHLAY"):
			t, err := time.Parse(time.UnixDate, string(h.Value));
			if err != nil {
				return meta, fmt.Errorf("Calculating time remaining from GOHLAY header: %v", err)
			}
			meta.DeliveryTime = t.UnixMilli()
			meta.Gohlayed = true
		case config.HeaderOverride("GOHLAY_DELIVERED"):
			meta.Gohlayed = true
			meta.DeliveryKey = string(h.Value)
			meta.Delivered = true
		}
	}
	return meta, nil
}
