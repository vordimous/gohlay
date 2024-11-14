package internal

import (
	"os"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	log "github.com/sirupsen/logrus"
	"github.com/vordimous/gohlay/common"
	"github.com/vordimous/gohlay/config"
)
var (
	producer    *kafka.Producer
)

// HandleDeliveries produces the gohlayed kafka messages
func HandleDeliveries() {
	log.Info("doDeliver")
	p, err := kafka.NewProducer(config.GetProducer())
	if err != nil {
		log.Error("Failed to create producer ", err)
		os.Exit(1)
	}
	producer = p

	didD := 0
	defer func() {
		producer.Close()
		log.Info("# delivered | ", didD)
	}()
	// Delivery report handler for produced messages
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Debug("Delivery failed: | ", ev.TopicPartition)
				} else {
					log.Debug("Delivered message to | ", ev.TopicPartition)
					didD++
				}
			default:
				log.Debug("Ignored produce | ", ev)
			}
		}
	}()
	log.Info("# to deliver | ", len(isDelivered))
	ScanTopic(deliverMsg)

	// Wait for message deliveries before shutting down
	for pending := 1; pending > 0; pending = producer.Flush(1000) {
		log.Info("waiting for message remaining deliveries | ", pending)
	}
	log.Info("fin")
}

func deliverMsg(msg *kafka.Message) {
	if delay, fin, _, gohlayed := common.ParseHeaders(msg.Headers); gohlayed && !fin {
		key := common.FmtKafkaKey(msg.TopicPartition.Offset, delay)
		log.Info("Found gohlayed message | ", msg)
		if delivered, exists := isDelivered[key]; exists && !delivered {
			var headers = []kafka.Header{}
			for _, h := range msg.Headers {
				if h.Key == "GOHLAY" {
					headers = append(headers,
						kafka.Header{
							Key:   "GOHLAY_FIN",
							Value: []byte(key),
						})
				} else {
					headers = append(headers, h)
				}
			}
			deliveryMsg := &kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: msg.TopicPartition.Topic, Partition: msg.TopicPartition.Partition},
				Value:          msg.Value,
				Key:            msg.Key,
				Headers:        headers,
			}
			producer.Produce(deliveryMsg, nil)
			log.Info("Delivered gohlayed message | ", deliveryMsg)
		}
	}
}
