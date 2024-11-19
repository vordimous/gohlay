package config

import (
	"reflect"
	"testing"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/spf13/viper"
)

func Test_getBootrapServersString(t *testing.T) {
	tests := []struct {
		name              string
		want              string
		bootstrap_servers []string
	}{
		{
			name:              "Test single server",
			bootstrap_servers: []string{"localhost:9092"},
			want:              "localhost:9092",
		},
		{
			name:              "Test multiple servers",
			bootstrap_servers: []string{"localhost:9092", "localhost:9093"},
			want:              "localhost:9092,localhost:9093",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set("bootstrap_servers", tt.bootstrap_servers)
			if got := getBootrapServersString(); got != tt.want {
				t.Errorf("getBootrapServersString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetConsumer(t *testing.T) {
	tests := []struct {
		name              string
		bootstrap_servers []string
		wantConfig        *kafka.ConfigMap
	}{
		{
			name:              "Test single server and topic",
			bootstrap_servers: []string{"localhost:9092"},
			wantConfig: &kafka.ConfigMap{
				"bootstrap.servers":               "localhost:9092",
				"go.application.rebalance.enable": true,
				"session.timeout.ms":              6000,
				"enable.partition.eof":            true,
				"default.topic.config":            kafka.ConfigMap{"auto.offset.reset": "earliest"},
			},
		},
		{
			name:              "Test multiple servers and topics",
			bootstrap_servers: []string{"localhost:9092", "localhost:9093"},
			wantConfig: &kafka.ConfigMap{
				"bootstrap.servers":               "localhost:9092,localhost:9093",
				"go.application.rebalance.enable": true,
				"session.timeout.ms":              6000,
				"enable.partition.eof":            true,
				"default.topic.config":            kafka.ConfigMap{"auto.offset.reset": "earliest"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set("bootstrap_servers", tt.bootstrap_servers)
			gotConfig := GetConsumer()
			if !reflect.DeepEqual(gotConfig, tt.wantConfig) {
				t.Errorf("GetConsumer() gotConfig = %v, want %v", gotConfig, tt.wantConfig)
			}
		})
	}
}

func TestGetProducer(t *testing.T) {
	tests := []struct {
		name              string
		bootstrap_servers []string
		wantConfig        *kafka.ConfigMap
	}{
		{
			name:              "Test single server and topic",
			bootstrap_servers: []string{"localhost:9092"},
			wantConfig: &kafka.ConfigMap{
				"bootstrap.servers": "localhost:9092",
			},
		},
		{
			name:              "Test multiple servers and topics",
			bootstrap_servers: []string{"localhost:9092", "localhost:9093"},
			wantConfig: &kafka.ConfigMap{
				"bootstrap.servers": "localhost:9092,localhost:9093",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set("bootstrap_servers", tt.bootstrap_servers)
			if gotConfig := GetProducer(); !reflect.DeepEqual(gotConfig, tt.wantConfig) {
				t.Errorf("GetProducer() = %v, want %v", gotConfig, tt.wantConfig)
			}
		})
	}
}
