package kafkautil

import (
	"fmt"
	"time"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/vordimous/gohlay/configs"
)

type GohlayedMeta struct {
	Gohlayed     bool
	DeliveryTime int64
	DeliveredMsg bool
	DeliveryKey  string
}

// ParseHeaders extracts relevant information from the message headers
func ParseHeaders(headers *[]kafka.Header) (meta GohlayedMeta, err error) {
	for _, h := range *headers {
		switch h.Key {
		case configs.HeaderOverride("GOHLAY"):
			t, err := time.Parse(time.UnixDate, string(h.Value))
			if err != nil {
				return meta, fmt.Errorf("calculating time remaining from GOHLAY header: %v", err)
			}
			meta.DeliveryTime = t.UnixMilli()
			meta.Gohlayed = true
		case configs.HeaderOverride("GOHLAY_DELIVERED"):
			meta.Gohlayed = true
			meta.DeliveryKey = string(h.Value)
			meta.DeliveredMsg = true
		}
	}
	return
}
