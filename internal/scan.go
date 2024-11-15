package internal

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"syscall"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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

// ScanTopic creates a unique consumer that reads all messages on the topics
func ScanTopic(handleMessage func(*kafka.Message)) {
	topicConfigMap, topics := config.GetConsumer()
	group := fmt.Sprintf("group.id=%d-%d | ", viper.GetInt64("deadline"), len(isDelivered))
	topicConfigMap.Set(group)
	log.Debug("scan with group.id | ", group)
	log.Debug(topicConfigMap)
	c, err := kafka.NewConsumer(topicConfigMap)
	if err != nil {
		log.Fatal("Failed to create consumer ", err)
		os.Exit(1)
	}
	if err := c.SubscribeTopics(topics, nil); err != nil {
		log.Fatal("Failed subscribe ", err)
		os.Exit(1)
	}

	defer func() {
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
				log.Debug("AssignedPartitions | ", e)
				parts, err := common.GetAssignedPartitions(c, e.Partitions)
				if err != nil {
					log.Fatal("Failed to get offset | ", err)
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
				log.Debug("%% Reached | ", e)
				maxOffset = e.Offset
				log.Debug("%% maxOffset | ", maxOffset)
				run = false
			case kafka.RevokedPartitions:
				log.Debug("%% Revoked | ", e)
				run = false
			case kafka.Error:
				log.Error("%% Error | ", e)
				run = false
			default:
				log.Debug("Ignored | ", e)
			}
		}
	}
}

