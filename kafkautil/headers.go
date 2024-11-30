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
	DeliveredMsg bool
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
				return meta, fmt.Errorf("calculating time remaining from GOHLAY header: %v", err)
			}
			meta.DeliveryTime = t.UnixMilli()
			meta.Gohlayed = true
		case config.HeaderOverride("GOHLAY_DELIVERED"):
			meta.Gohlayed = true
			meta.DeliveryKey = string(h.Value)
			meta.DeliveredMsg = true
		}
	}
	return meta, nil
}
