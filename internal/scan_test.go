package internal

import (
	"testing"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/spf13/viper"
	"github.com/vordimous/gohlay/configs"
)

type ScanTester struct {}

// TopicName is the name of the kafka topic
func (s *ScanTester) TopicName() string {
	return "gohlay"
}

// GroupName is a human readable name of the purpose for the message handler
func (s *ScanTester) GroupName() string {
	return "testScan"
}

// HandleMessage will deliver any gohlayed message that isn't already delivered
func (s *ScanTester) HandleMessage(msg *kafka.Message) string {
	return ""
}

func BenchmarkScanTopic(b *testing.B) {
	type args struct {
		handler MessageHandler
	}
	benchmarks := []struct {
		name string
		args args
	}{
		{
			name: "No read",
			args: args{
				handler: &ScanTester{},
			},
		},
	}
	viper.Set("bootstrap-servers", "localhost:9092")
	viper.Set("topics", "gohlay")
	configs.Load()
	for _, bm := range benchmarks {
		b.ReportAllocs()
		b.ResetTimer()
		b.Run(bm.name, func(b *testing.B) {
			ScanTopic(bm.args.handler)
		})
	}
}
