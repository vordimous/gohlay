package internal

import (
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
		log.Infof("Number of gohlayed messages delivered: %v", didD)
	}()
	// Delivery report handler for produced messages
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Errorf("Delivery failed: %+v", ev.TopicPartition)
				} else {
					didD++
				}
			default:
				log.Debugf("Ignored produce: %+v", ev)
			}
		}
	}()
	log.Infof("Number of gohlayed messages on topic: %v", len(isDelivered))

	ScanAll("deliver", deliverMsg)

	// Wait for message deliveries before shutting down
	for pending := 1; pending > 0; pending = producer.Flush(1000) {
		log.Debug("Waiting for remaining deliveries")
	}
}

func deliverMsg(msg *kafka.Message) {
	if delay, delivered, _, gohlayed := common.ParseHeaders(msg.Headers); gohlayed && !delivered {
		messageId := common.FmtMessageId(msg.TopicPartition.Offset, delay)
		log.Debugf("Found deliverable message: %d %d-%s %s", msg.TopicPartition.Partition, msg.TopicPartition.Offset, msg.Key, messageId)
		if delivered, exists := isDelivered[messageId]; exists && !delivered {
			var headers = []kafka.Header{}

			// replace the deadline header with the delivered header
			for _, h := range msg.Headers {
				if h.Key == config.GetHeaderOverride("GOHLAY") {
					headers = append(headers,
						kafka.Header{
							Key:   config.GetHeaderOverride("GOHLAY_DELIVERED"),
							Value: []byte(messageId),
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
			log.Infof("Delivered message: %s %s", deliveryMsg.Key, messageId)
		} else if delivered {
			log.Debugf("Message already delivered: %d-%s %s", msg.TopicPartition.Offset, msg.Key, messageId)
		}
	}
}
