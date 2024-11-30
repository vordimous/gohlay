package kafkautil

import (
	"testing"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func TestParseHeaders(t *testing.T) {
	type args struct {
		headers []kafka.Header
	}
	tests := []struct {
		name               string
		args               args
		want  GohlayedMeta
	}{
		{
			name: "Gohlayed message",
			args: args{
				headers: []kafka.Header{
					{
						Key:   "GOHLAY",
						Value: []byte("Fri Feb 13 23:31:30 UTC 2009"),
					},
				},
			},
			want: GohlayedMeta{
				DeliveryTime:   1234567890000,
				Delivered:    false,
				DeliveryKey: "",
				Gohlayed:       true,
			},
		},
		{
			name: "Gohlayed message is Delivered",
			args: args{
				headers: []kafka.Header{
					{
						Key:   "GOHLAY_DELIVERED",
						Value: []byte("9-1234567890"),
					},
				},
			},
			want: GohlayedMeta{
				DeliveryTime:   0,
				Delivered:    true,
				DeliveryKey: "9-1234567890",
				Gohlayed:       true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := ParseHeaders(tt.args.headers)
			if got.DeliveryTime != tt.want.DeliveryTime {
				t.Errorf("ParseHeaders() got.DeliveryTime = %v, want %v", got.DeliveryTime, tt.want.DeliveryTime)
			}
			if got.Delivered != tt.want.Delivered {
				t.Errorf("ParseHeaders() got.Delivered = %v, want %v", got.Delivered, tt.want.Delivered)
			}
			if got.DeliveryKey != tt.want.DeliveryKey {
				t.Errorf("ParseHeaders() got.DeliveryKey = %v, want %v", got.DeliveryKey, tt.want.DeliveryKey)
			}
			if got.Gohlayed != tt.want.Gohlayed {
				t.Errorf("ParseHeaders() got.Gohlayed = %v, want %v", got.Gohlayed, tt.want.Gohlayed)
			}
		})
	}
}
