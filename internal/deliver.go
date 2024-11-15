package internal

import (
	"fmt"
	"os"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	log "github.com/sirupsen/logrus"
	"github.com/vordimous/gohlay/common"
	"github.com/vordimous/gohlay/config"
)

var (
	producer *kafka.Producer
)

// HandleDeliveries produces the gohlayed kafka messages
func HandleDeliveries() {
	p, err := kafka.NewProducer(config.GetProducer())
	if err != nil {
		log.Fatal("Failed to create producer ", err)
		os.Exit(1)
	}
	producer = p

	didD := 0
	defer func() {
		producer.Close()
		log.Info("Number of golayed messages delivered | ", didD)
	}()
	// Delivery report handler for produced messages
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Error("Delivery failed: | ", ev.TopicPartition)
				} else {
					didD++
				}
			default:
				log.Debug("Ignored produce | ", ev)
			}
		}
	}()
	log.Info("Number of golayed messages on topic | ", len(isDelivered))
	ScanTopic(deliverMsg)

	// Wait for message deliveries before shutting down
	for pending := 1; pending > 0; pending = producer.Flush(1000) {
		log.Debug("Waiting for message remaining deliveries")
	}
}

func deliverMsg(msg *kafka.Message) {
	if delay, delivered, _, gohlayed := common.ParseHeaders(msg.Headers); gohlayed && !delivered {
		deliveryKey := common.FmtKafkaKey(msg.TopicPartition.Offset, delay)
		log.Debug("Found deliverable message | ", fmt.Sprintf("%d-%s %s", msg.TopicPartition.Offset, msg.Key, deliveryKey))
		if delivered, exists := isDelivered[deliveryKey]; exists && !delivered {
			var headers = []kafka.Header{}

			// replace the deadline header with the delivered header
			for _, h := range msg.Headers {
				if h.Key == "GOHLAY" {
					headers = append(headers,
						kafka.Header{
							Key:   "GOHLAY_DELIVERED",
							Value: []byte(deliveryKey),
						})
				} else {
					headers = append(headers, h)
				}
			}

			deliveryMsg := &kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: msg.TopicPartition.Topic, Partition: msg.TopicPartition.Partition},
				Value:          msg.Value,
				Key:            msg.Key,
				Opaque:         msg.Opaque,
				Headers:        headers,
			}
			producer.Produce(deliveryMsg, nil)
			log.Info("Delivered message | ", fmt.Sprintf("%s %s", deliveryMsg.Key, deliveryKey))
		} else if delivered {
			log.Debug("Message already delivered | ", fmt.Sprintf("%d-%s %s", msg.TopicPartition.Offset, msg.Key, deliveryKey))
		}
	}
}
