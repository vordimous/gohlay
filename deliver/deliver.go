package deliver

import (
	"fmt"
	"os"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	log "github.com/sirupsen/logrus"
	"github.com/vordimous/gohlay/config"
	"github.com/vordimous/gohlay/internal"
	"github.com/vordimous/gohlay/kafkautil"
)

// HandleDeliveries produces the gohlayed kafka messages
func HandleDeliveries(topic string, deliveryKeyMap map[string]bool) {
	topicConfigMap := config.Producer()
	p, err := kafka.NewProducer(&topicConfigMap)
	if err != nil {
		log.Fatal("Failed to create producer ", err)
		os.Exit(1)
	}
	d := &Deliverer{
		topic: topic,
		producer: p,
		deliveryKeyMap: deliveryKeyMap,
	}
	d.doDeliveries()
}

type Deliverer struct {
	topic string
	deliveryKeyMap map[string]bool
	producer *kafka.Producer
}

// TopicName is the name of the kafka topic
func (d *Deliverer) TopicName() string {
	return d.topic
}

// GroupName is a human readable name of the purpose for the message handler
func (d *Deliverer) GroupName() string {
	return "delivering"
}

// HandleMessage will deliver any gohlayed message that isn't already delivered
func (d *Deliverer) HandleMessage(msg *kafka.Message) string {
	gohlayedMeta, err := kafkautil.ParseHeaders(msg.Headers)
	if err != nil {
		return fmt.Sprintf("could't parse headers: %v", err)
	}

	if !gohlayedMeta.Gohlayed {
		return fmt.Sprintf("message is not Gohlayed: %+v", msg.TopicPartition)

	}
	deliveryKey := kafkautil.FmtDeliveryKey(msg.TopicPartition.Offset, gohlayedMeta.DeliveryTime)
	if gohlayedMeta.Delivered || d.deliveryKeyMap[deliveryKey] {
		return fmt.Sprintf("message already delivered: %+v %s", msg.TopicPartition.Offset, gohlayedMeta.DeliveryKey)

	}

	log.Debugf("Found deliverable message: %d %d-%s %s", msg.TopicPartition.Partition, msg.TopicPartition.Offset, msg.Key, deliveryKey)
	var headers = []kafka.Header{}

	// replace the deadline header with the delivered header
	for _, h := range msg.Headers {
		if h.Key == config.HeaderOverride("GOHLAY") {
			headers = append(headers,
				kafka.Header{
					Key:   config.HeaderOverride("GOHLAY_DELIVERED"),
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

	if err := d.producer.Produce(deliveryMsg, nil); err != nil {
		log.Fatalf("Failed to produce message: %v", err)
		return fmt.Sprintf("Failed to produce message: %v", err)
	}
	log.Infof("Delivered message: %s %s", deliveryMsg.Key, deliveryKey)
	return ""
}

func (d *Deliverer) doDeliveries() {
	didD := 0
	defer func() {
		d.producer.Close()
		log.Infof("Number of gohlayed messages delivered: %v", didD)
	}()

	// Delivery report handler for produced messages
	go func() {
		for e := range d.producer.Events() {
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
	log.Infof("Number of gohlayed messages on topic: %v", len(d.deliveryKeyMap))

	internal.ScanTopic(d)

	// Wait for message deliveries before shutting down
	for pending := 1; pending > 0; pending = d.producer.Flush(1000) {
		log.Debug("Waiting for remaining deliveries")
	}
}
