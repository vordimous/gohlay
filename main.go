package main

import (
	"encoding/binary"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	defer func() {
		err = c.Close()
		if err != nil {
			panic(err)
		}
	}()

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{"myTopic", "^aRegex.*[Tt]opic"}, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s | %d: %v\n", msg.TopicPartition, binary.BigEndian.Uint32(msg.Key), getDelay(msg.Headers))
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}

func getDelay(headers []kafka.Header) uint32 {
	d := uint32(0)
	for _, h := range headers {
		if h.Key == "delay" {
			d = binary.BigEndian.Uint32(h.Value)
		}
	}
	return d
}
