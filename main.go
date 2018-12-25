package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	stash map[uint32]uint32
)

func init() {
	stash = map[uint32]uint32{}
}

func main() {
	// if len(os.Args) < 4 {
	// 	fmt.Fprintf(os.Stderr, "Usage: %s <broker> <group> <topics..>\n",
	// 		os.Args[0])
	// 	os.Exit(1)
	// }

	broker := "localhost"
	// broker := os.Args[1]
	group := "myGroup"
	// group := os.Args[2]
	topics := []string{"myTopic", "^aRegex.*[Tt]opic"}
	// topics := os.Args[3:]
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":               broker,
		"group.id":                        group,
		"go.application.rebalance.enable": true, // delegate Assign() responsibility to app
		"session.timeout.ms":              6000,
		"default.topic.config":            kafka.ConfigMap{"auto.offset.reset": "earliest"}})

	defer func() {
		fmt.Printf("Closing consumer\n")
		c.Close()
	}()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Consumer %v\n", c)

	err = c.SubscribeTopics(topics, nil)

	run := true

	for run == true {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := c.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case kafka.AssignedPartitions:
				parts := make([]kafka.TopicPartition,
					len(e.Partitions))
				for i, tp := range e.Partitions {
					tp.Offset = kafka.OffsetTail(5) // Set start offset to 5 messages from end of partition
					parts[i] = tp
				}
				fmt.Printf("Assign %v\n", parts)
				c.Assign(parts)
			case *kafka.Message:
				handleMessage(e)
			case kafka.PartitionEOF:
				fmt.Printf("%% Reached %v\n", e)
				run = false
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
				run = false
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}
}

func handleMessage(msg *kafka.Message) {
	delay := getDelay(msg.Headers)
	key := binary.BigEndian.Uint32(msg.Key)
	stash[key] = delay
	fmt.Printf("Message %d | %d: %v\n", len(stash), key, delay)
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
