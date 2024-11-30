package deliver

import (
	"testing"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func TestDeliverer_HandleMessage(t *testing.T) {
	topic := "test"
	type args struct {
		msg *kafka.Message
	}
	tests := []struct {
		name string
		d    *Deliverer
		args args
		want string
	}{
		{
			name: "Not Gohlayed",
			d: &Deliverer{
				topic: topic,
				deliveryKeyMap: map[string]bool{
					"test": true,
				},
				producer: &kafka.Producer{},
			},
			args: args{
				msg: &kafka.Message{
					TopicPartition: kafka.TopicPartition{
						Topic: &topic,
						Partition: 9,
						Offset: 9,
					},
				},
			},
			want: "message is not Gohlayed: test[9]@9",
		},
		{
			name: "Already Delivered Header",
			d: &Deliverer{
				topic: topic,
				deliveryKeyMap: map[string]bool{},
				producer: &kafka.Producer{},
			},
			args: args{
				msg: &kafka.Message{
					TopicPartition: kafka.TopicPartition{
						Topic: &topic,
						Partition: 9,
						Offset: 9,
					},
					Headers: []kafka.Header{
						{
							Key:   "GOHLAY_DELIVERED",
							Value: []byte("9-1234567890000"),
						},
					},
				},
			},
			want: "message already delivered: 9 9-1234567890000",
		},
		{
			name: "Already Delivered Map",
			d: &Deliverer{
				topic: topic,
				deliveryKeyMap: map[string]bool{
					"9-1234567890000": true,
				},
				producer: &kafka.Producer{},
			},
			args: args{
				msg: &kafka.Message{
					TopicPartition: kafka.TopicPartition{
						Topic: &topic,
						Partition: 9,
						Offset: 9,
					},
					Headers: []kafka.Header{
						{
							Key:   "GOHLAY",
							Value: []byte("Fri Feb 13 23:31:30 UTC 2009"),
						},
					},
				},
			},
			want: "message already delivered: 9 ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.HandleMessage(tt.args.msg); got != tt.want {
				t.Errorf("Deliverer.HandleMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
