package kafkautil

import (
	"testing"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/go-cmp/cmp"
	"github.com/spf13/viper"
)

func TestFmtDeliveryKey(t *testing.T) {
	type args struct {
		offset kafka.Offset
		delay  int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Correct Format",
			args: args{
				offset: 0,
				delay:  12345,
			},
			want: "0-12345",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FmtDeliveryKey(tt.args.offset, tt.args.delay); !cmp.Equal(got, tt.want) {
				t.Errorf("FmtDeliveryKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFmtKafkaGroup(t *testing.T) {
	type args struct {
		groupName string
		topic     string
	}
	tests := []struct {
		name     string
		args     args
		want     string
		deadline int64
	}{
		{
			name: "Correct Format",
			args: args{
				groupName: "test",
				topic:     "aTopic",
			},
			deadline: 12345,
			want:     "group.id=gohlay_test:aTopic:12345",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set("deadline", tt.deadline)
			if got := FmtKafkaGroup(tt.args.groupName, tt.args.topic); !cmp.Equal(got, tt.want) {
				t.Errorf("FmtKafkaGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}
