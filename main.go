package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	stash map[uint32]uint32
)

func init() {
	stash = map[uint32]uint32{}
}

func main() {
	broker := "localhost"
	group := "myGroup"
	topics := []string{"myTopic", "^aRegex.*[Tt]opic"}
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
				limit := time.Now().Add(time.Duration(-24)*time.Hour).UnixNano() / int64(time.Millisecond)
				parts := make([]kafka.TopicPartition, len(e.Partitions))
				for i, tp := range e.Partitions {
					offset, _ := kafka.NewOffset(limit)
					tp.Offset = offset
					fmt.Printf("offset query: %v\n", tp.Offset)
					// tp.Offset = kafka.OffsetTail(5) // Set start offset to 5 messages from end of partition
					parts[i] = tp
				}
				fmt.Printf("Assign %v\n", parts)
				fmt.Printf("time limit: %d\n", limit)
				parts, err = c.OffsetsForTimes(parts, 10000)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to get offset: %s\n", err)
					os.Exit(1)
				}
				for _, tp := range parts {
					fmt.Printf("offset: %v\n", tp.Offset)
				}
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
