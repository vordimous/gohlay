package internal

import (
	"math"
	"os"
	"os/signal"
	"syscall"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	log "github.com/sirupsen/logrus"
	"github.com/vordimous/gohlay/common"
	"github.com/vordimous/gohlay/config"
)

var (
	sigchan     chan os.Signal
	maxOffset   kafka.Offset
)

func init() {
	sigchan = make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	maxOffset, _ = kafka.NewOffset(math.MaxInt64)
}

// ScanAll creates a unique consumer that reads all messages on the topics
func ScanAll(reason string, handleMessage func(*kafka.Message)) {
	for _, topic := range config.GetTopics() {
		scanTopic(topic, reason, handleMessage)
    }
}

// scanTopic creates a unique consumer that reads all messages on the topics
func scanTopic(topic string, reason string, handleMessage func(*kafka.Message)) {

	topicConfigMap := config.GetConsumer()
	if err := topicConfigMap.Set(common.FmtKafkaGroup(reason, topic)); err != nil {
		log.Fatalf("Failed to set the consumer groupId %v", err)
		os.Exit(1)
	}

	log.Debugf("Scanning with %+v", topicConfigMap)
	c, err := kafka.NewConsumer(topicConfigMap)
	if err != nil {
		log.Fatal("Failed to create consumer ", err)
		os.Exit(1)
	}

	defer func() {
		log.Debugf("Closing consumer %+v", topicConfigMap)
		c.Close()
	}()

	partitions := []int32{}
	if metadata, err :=c.GetMetadata(&topic, false, 100); err != nil {
		log.Warning("Failed to get Partitions, using partition 0", err)
		partitions = append(partitions, 0)
	} else {
		for _, p := range metadata.Topics[topic].Partitions {
			partitions = append(partitions, p.ID)
		}
	}
	topicPartitions := []kafka.TopicPartition{}
	for _, partition := range partitions {
		topicPartitions = append(topicPartitions, kafka.TopicPartition{
			Topic: &topic,
			Partition: partition,
		})
	}
	if err := c.Assign(topicPartitions); err != nil {
		log.Fatal("Failed subscribe ", err)
		os.Exit(1)
	}
	log.Debugf("Scanning %d Partitions: %+v", len(topicPartitions), partitions)

	run := true
	for run {
		select {
		case sig := <-sigchan:
			log.Errorf("Caught signal, terminating: %v", sig)
			run = false
		default:
			ev := c.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				if e.TopicPartition.Offset < maxOffset {
					handleMessage(e)
				} else {
					run = false
				}
			case kafka.PartitionEOF:
				maxOffset = e.Offset
				log.Debugf("%% Reached maxOffset: %v %v %+v", maxOffset, e.Partition, e)
				run = false
			case kafka.RevokedPartitions:
				log.Debugf("%% Revoked: %v", e)
				run = false
			case kafka.Error:
				log.Errorf("%% Error: %v", e)
				run = false
			default:
				log.Debugf("Ignored: %v", e)
			}
		}
	}
}

