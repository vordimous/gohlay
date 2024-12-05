package internal

import (
	"math"
	"os"
	"os/signal"
	"syscall"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	log "github.com/sirupsen/logrus"
	"github.com/vordimous/gohlay/configs"
	"github.com/vordimous/gohlay/internal/kafkautil"
)

var (
	sigchan   chan os.Signal
	maxOffset kafka.Offset
)

func init() {
	sigchan = make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	maxOffset, _ = kafka.NewOffset(math.MaxInt64)
}

type MessageHandler interface {
	TopicName() string
	GroupName() string
	HandleMessage(*kafka.Message) string
}

// ScanTopic creates a unique consumer that reads all messages on the topics
func ScanTopic(handler MessageHandler) {
	topicConfigMap := configs.Consumer()
	topic := handler.TopicName()
	partitionOffsets := map[int32]kafka.Offset{}

	if err := topicConfigMap.Set(kafkautil.FmtKafkaGroup(handler.GroupName(), topic)); err != nil {
		log.Fatalf("Failed to set the consumer groupId %v", err)
		os.Exit(1)
	}

	log.Debugf("Scanning with %+v", topicConfigMap)
	c, err := kafka.NewConsumer(&topicConfigMap)
	if err != nil {
		log.Fatal("Failed to create consumer ", err)
		os.Exit(1)
	}

	defer func() {
		log.Debugf("Closing consumer %+v", topicConfigMap)
		c.Close()
	}()

	partitions := []int32{}
	if metadata, err := c.GetMetadata(&topic, false, 100); err != nil {
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
			Topic:     &topic,
			Partition: partition,
		})
		partitionOffsets[partition] = maxOffset
	}
	if len(topicPartitions) == 0 {
		log.Infof("No partitions found for topic: %s", topic)
		return
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
				partitionOffsets[e.TopicPartition.Partition] = e.TopicPartition.Offset
				if res := handler.HandleMessage(e); res != "" {
					log.Debugf("Handling Message result: %s", res)
				}
			case kafka.PartitionEOF:
				log.Debugf("kafka.Event PartitionEOF; Reached the end of the partition: %v %v %+v", partitionOffsets[e.Partition], e.Partition, e)
				delete(partitionOffsets, e.Partition)

				// stop scanning once all partitions have reached the end
				if len(partitionOffsets) == 0 {
					run = false
				}
			case kafka.RevokedPartitions:
				log.Debugf("kafka.Event RevokedPartitions: %v", e)
				for p := range e.Partitions {
					delete(partitionOffsets, int32(p))
				}

				// stop scanning once all partitions have reached the end
				if len(partitionOffsets) == 0 {
					run = false
				}
			case kafka.Error:
				log.Errorf("kafka.Event Error: %v", e)
				run = false
			default:
				log.Debugf("kafka.Event Ignored: %v", e)
			}
		}
	}
}
