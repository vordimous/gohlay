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
	group := fmt.Sprintf("group.id=%d-%d", viper.GetInt64("deadline"), len(isDelivered))
	topicConfigMap.Set(group)
	log.Debugf("Scanning with %+v", topicConfigMap)
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
			log.Errorf("Caught signal, terminating: %v", sig)
			run = false
		default:
			ev := c.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case kafka.AssignedPartitions:
				log.Debugf("AssignedPartitions: %v", e)
				parts, err := common.GetAssignedPartitions(c, e.Partitions)
				if err != nil {
					log.Fatalf("Failed to get offset: %v", err)
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
				maxOffset = e.Offset
				log.Debugf("%% Reached maxOffset: %v %v", maxOffset, e)
				run = false
			case kafka.RevokedPartitions:
				log.Infof("%% Revoked: %v", e)
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

