package main

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var (
	isDelivered    map[string]bool
	sigchan        chan os.Signal
	broker         string
	topics         []string
	topicConfigMap *kafka.ConfigMap
	selfProducer   *kafka.Producer
	maxOffset      kafka.Offset
	timeNow        uint64
)

func init() {
	_, kafkaVersion := kafka.LibraryVersion()
	fmt.Fprintf(os.Stderr, "Kafka Lib Version: %v\n", kafkaVersion)
	isDelivered = map[string]bool{}
	sigchan = make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	timeNow = uint64(time.Now().UnixNano() / int64(time.Millisecond))
	broker = "localhost:9092"
	topics = []string{"myTopic"}
	maxOffset, _ = kafka.NewOffset(math.MaxInt64)
}

func main() {
	topicConfigMap = &kafka.ConfigMap{
		"bootstrap.servers":               broker,
		"go.application.rebalance.enable": true, // delegate Assign() responsibility to app
		"session.timeout.ms":              6000,
		"enable.partition.eof":            true,
		"default.topic.config":            kafka.ConfigMap{"auto.offset.reset": "earliest"},
	}
	scanTopic(indexMsg)
	doDeliver()
}

func scanTopic(handleMessage func(*kafka.Message)) {
	group := fmt.Sprintf("group.id=%d-%d", timeNow, len(isDelivered))
	fmt.Printf("scanTopic group %s\n", group)
	topicConfigMap.Set(group)
	c, err := kafka.NewConsumer(topicConfigMap)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}
	if err := c.SubscribeTopics(topics, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Failed subscribe: %s\n", err)
		os.Exit(1)
	}

	defer func() {
		fmt.Printf("Closing consumer\n")
		c.Close()
	}()

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
				parts, err := getPartitions(c, e.Partitions)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to get offset: %s\n", err)
					os.Exit(1)
				}
				c.Assign(parts)
			case *kafka.Message:
				if e.TopicPartition.Offset < maxOffset {
					handleMessage(e)
				} else {
					run = false
				}
			case kafka.PartitionEOF:
				fmt.Printf("%% Reached %v\n", e)
				maxOffset = e.Offset
				run = false
			case kafka.RevokedPartitions:
				fmt.Printf("%% Revoked %v\n", e)
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

func indexMsg(msg *kafka.Message) {
	fmt.Printf("indexMsg %v\n", msg)
	delay, finKey, hasHeader := getDelay(msg.Headers)
	if hasHeader {
		if finKey == "" && delay != 0 {
			if delay < timeNow {
				isDelivered[getKey(msg.TopicPartition.Offset, delay)] = false // set key to be delivered with a delivery value of false
				fmt.Println("isDelivered: ", getKey(msg.TopicPartition.Offset, delay))
			} else {
				fmt.Println("not time yet: ", msg.TopicPartition.Offset)
			}
		} else if finKey != "" {
			delete(isDelivered, finKey)
		}
	} else {
		fmt.Println("no gohlay: ", msg.TopicPartition.Offset)
	}
}

func getKey(offset kafka.Offset, delay uint64) string {
	fmt.Printf("getKey %v %d\n", offset, delay)
	return fmt.Sprintf("%v-%d", offset, delay)
}

func getDelay(headers []kafka.Header) (delay uint64, key string, exists bool) {
	fmt.Printf("getDelay %v\n", headers)
	for _, h := range headers {
		if h.Key == "GOHLAY" {
			timeString := string(h.Value)
			if deliveryTime, err := time.Parse(time.UnixDate, timeString); err == nil {
				delay = uint64(deliveryTime.UnixNano() / int64(time.Millisecond))
			} else {
				fmt.Fprintf(os.Stderr, "Reading GOHLAY header: %s\n", err)
			}
			exists = true
		}
		if h.Key == "GOHLAY_FIN" {
			key = string(h.Value)
			exists = true
		}
	}
	return
}

func getPartitions(c *kafka.Consumer, partitions []kafka.TopicPartition) ([]kafka.TopicPartition, error) {
	fmt.Printf("%v %v\n", c, partitions)
	parts := make([]kafka.TopicPartition, len(partitions))
	var err error
	if false {
		limit := time.Now().Add(time.Duration(-5)*time.Minute).UnixNano() / int64(time.Millisecond)
		for i, tp := range partitions {
			offset, _ := kafka.NewOffset(limit)
			tp.Offset = offset
			fmt.Printf("offset query time: %v\n", tp.Offset)
			parts[i] = tp
		}
		parts, err = c.OffsetsForTimes(parts, 10000)
	} else {
		for i, tp := range partitions {
			offset, _ := kafka.NewOffset(0)
			tp.Offset = offset
			fmt.Printf("offset query value: %v\n", tp.Offset)
			parts[i] = tp
		}
	}
	fmt.Printf("Assign partition(s) %v\n", parts)
	return parts, err
}

func doDeliver() {
	fmt.Printf("doDeliver\n")
	if p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": broker}); err != nil {
		panic(err)
	} else {
		selfProducer = p
	}

	fmt.Printf("selfProducer %v\n", selfProducer)
	didD := 0
	defer func() {
		selfProducer.Close()
		fmt.Println("# delivered: ", didD)
	}()
	// Delivery report handler for produced messages
	go func() {
		for e := range selfProducer.Events() {
			fmt.Printf("e %v\n", e)
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
					didD++
				}
			default:
				fmt.Printf("Ignored produce %v\n", ev)
			}
		}
	}()
	fmt.Println("# to deliver: ", len(isDelivered))
	scanTopic(sendMsg)

	// Wait for message deliveries before shutting down
	for pending := 1; pending > 0; pending = selfProducer.Flush(1000) {
		fmt.Printf("waiting for message deliveries, %d remaining\n", pending)
	}
	fmt.Println("fin")
}

func sendMsg(msg *kafka.Message) {
	fmt.Printf("sendMsg %v\n", msg)
	if delay, fin, exists := getDelay(msg.Headers); exists && fin == "" {
		key := getKey(msg.TopicPartition.Offset, delay)
		delivered, exists := isDelivered[key]
		fmt.Printf("exists, delivered, key %v %v %v\n", exists, delivered, key)
		if delivered, exists := isDelivered[key]; exists && !delivered {
			fmt.Printf("exists, delivered %v %v\n", exists, delivered)
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
			selfProducer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: msg.TopicPartition.Topic, Partition: msg.TopicPartition.Partition},
				Value:          msg.Value,
				Key:            msg.Key,
				Headers:        headers,
			}, nil)
		}
	}
}
