package internal

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	log "github.com/sirupsen/logrus"
	"github.com/vordimous/gohlay/config"
)

var (
	isDelivered map[string]bool
	sigchan     chan os.Signal
	producer    *kafka.Producer
	maxOffset   kafka.Offset
	timeNow     uint64
)

func init() {
	isDelivered = map[string]bool{}
	sigchan = make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	timeNow = uint64(time.Now().UnixNano() / int64(time.Millisecond))
	maxOffset, _ = kafka.NewOffset(math.MaxInt64)
}

// Run will scan a configured topic and deliver messages
func Run() {
	scanTopic(indexMsg)
	doDeliver()
}

func scanTopic(handleMessage func(*kafka.Message)) {
	topicConfigMap, topics := config.GetConsumer()
	group := fmt.Sprintf("group.id=%d-%d | ", timeNow, len(isDelivered))
	topicConfigMap.Set(group)
	log.Info("scanTopic | ", topicConfigMap)
	log.Info("scanTopic | ", topics)
	c, err := kafka.NewConsumer(topicConfigMap)
	if err != nil {
		log.Error("Failed to create consumer ", err)
		os.Exit(1)
	}
	if err := c.SubscribeTopics(topics, nil); err != nil {
		log.Error("Failed subscribe ", err)
		os.Exit(1)
	}

	defer func() {
		log.Info("Closing consumer")
		c.Close()
	}()

	run := true
	for run {
		select {
		case sig := <-sigchan:
			log.Error("Caught signal, terminating | ", sig)
			run = false
		default:
			ev := c.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case kafka.AssignedPartitions:
				log.Info("scanTopic AssignedPartitions | ", e)
				parts, err := getPartitions(c, e.Partitions)
				if err != nil {
					log.Error("Failed to get offset| ", err)
					os.Exit(1)
				}
				c.Assign(parts)
			case *kafka.Message:
				log.Info("scanTopic Message | ", e)
				if e.TopicPartition.Offset < maxOffset {
					handleMessage(e)
				} else {
					run = false
				}
			case kafka.PartitionEOF:
				log.Info("%% Reached | ", e)
				maxOffset = e.Offset
				log.Info("%% maxOffset | ", maxOffset)
				run = false
			case kafka.RevokedPartitions:
				log.Info("%% Revoked | ", e)
				run = false
			case kafka.Error:
				log.Error("%% Error | ", e)
				run = false
			default:
				log.Info("Ignored | ", e)
			}
		}
	}
}

func indexMsg(msg *kafka.Message) {
	log.Info("indexMsg | ", msg)
	delay, finKey, hasHeader := getDelay(msg.Headers)
	if hasHeader {
		if finKey == "" && delay != 0 {
			if delay < timeNow {
				isDelivered[getKey(msg.TopicPartition.Offset, delay)] = false // set key to be delivered with a delivery value of false
				log.Info("isDelivered | ", getKey(msg.TopicPartition.Offset, delay))
			} else {
				log.Info("not time yet | ", msg.TopicPartition.Offset)
			}
		} else if finKey != "" {
			delete(isDelivered, finKey)
		}
	} else {
		log.Info("no gohlay | ", msg.TopicPartition.Offset)
	}
}

func getKey(offset kafka.Offset, delay uint64) string {
	log.Info("getKey offset | ", offset)
	log.Info("getKey delay | ", delay)
	return fmt.Sprintf("%v-%d", offset, delay)
}

func getDelay(headers []kafka.Header) (delay uint64, key string, exists bool) {
	log.Info("getDelay | ", headers)
	for _, h := range headers {
		if h.Key == "GOHLAY" {
			timeString := string(h.Value)
			if deliveryTime, err := time.Parse(time.UnixDate, timeString); err == nil {
				delay = uint64(deliveryTime.UnixNano() / int64(time.Millisecond))
			} else {
				log.Error("Reading GOHLAY header ", err)
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
	log.Info("getPartitions | ", c, partitions)
	parts := make([]kafka.TopicPartition, len(partitions))
	var err error
	if false {
		limit := time.Now().Add(time.Duration(-5)*time.Minute).UnixNano() / int64(time.Millisecond)
		for i, tp := range partitions {
			offset, _ := kafka.NewOffset(limit)
			tp.Offset = offset
			log.Info("offset query time | ", tp.Offset)
			parts[i] = tp
		}
		parts, err = c.OffsetsForTimes(parts, 10000)
	} else {
		for i, tp := range partitions {
			offset, _ := kafka.NewOffset(0)
			tp.Offset = offset
			log.Info("offset query value | ", tp.Offset)
			parts[i] = tp
		}
	}
	log.Info("Assign partition(s) | ", parts)
	return parts, err
}

func doDeliver() {
	log.Info("doDeliver")
	p, err := kafka.NewProducer(config.GetProducer())
	if err != nil {
		log.Error("Failed to create producer ", err)
		os.Exit(1)
	}
	producer = p

	log.Info("producer | ", producer)
	didD := 0
	defer func() {
		producer.Close()
		log.Info("# delivered | ", didD)
	}()
	// Delivery report handler for produced messages
	go func() {
		for e := range producer.Events() {
			log.Info("producer.Event | ", e)
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Info("Delivery failed: | ", ev.TopicPartition)
				} else {
					log.Info("Delivered message to | ", ev.TopicPartition)
					didD++
				}
			default:
				log.Info("Ignored produce | ", ev)
			}
		}
	}()
	log.Info("# to deliver | ", len(isDelivered))
	scanTopic(sendMsg)

	// Wait for message deliveries before shutting down
	for pending := 1; pending > 0; pending = producer.Flush(1000) {
		log.Info("waiting for message remaining deliveries | ", pending)
	}
	log.Info("fin")
}

func sendMsg(msg *kafka.Message) {
	log.Info("sendMsg | ", msg)
	if delay, fin, exists := getDelay(msg.Headers); exists && fin == "" {
		key := getKey(msg.TopicPartition.Offset, delay)
		delivered, exists := isDelivered[key]
		log.Info("exists, delivered, key | ", exists, delivered, key)
		if delivered, exists := isDelivered[key]; exists && !delivered {
			log.Info("exists, delivered | ", exists, delivered)
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
			producer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: msg.TopicPartition.Topic, Partition: msg.TopicPartition.Partition},
				Value:          msg.Value,
				Key:            msg.Key,
				Headers:        headers,
			}, nil)
		}
	}
}
