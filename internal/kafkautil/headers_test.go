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
		name string
		args args
		want GohlayedMeta
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
				DeliveryTime: 1234567890000,
				DeliveredMsg: false,
				DeliveryKey:  "",
				Gohlayed:     true,
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
				DeliveryTime: 0,
				DeliveredMsg: true,
				DeliveryKey:  "9-1234567890",
				Gohlayed:     true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := ParseHeaders(&tt.args.headers)
			if got.DeliveryTime != tt.want.DeliveryTime {
				t.Errorf("ParseHeaders() got.DeliveryTime = %v, want %v", got.DeliveryTime, tt.want.DeliveryTime)
			}
			if got.DeliveredMsg != tt.want.DeliveredMsg {
				t.Errorf("ParseHeaders() got.Delivered = %v, want %v", got.DeliveredMsg, tt.want.DeliveredMsg)
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

func BenchmarkParseHeaders(b *testing.B) {
	// Request was generated from TestHTTPServerRequest request.
	headers := []kafka.Header{
		{
			Key:   "GOHLAY",
			Value: []byte("Fri Feb 13 23:31:30 UTC 2009"),
		},
		{
			Key:   "GOHLAY_DELIVERED",
			Value: []byte("Fri Feb 13 23:31:30 UTC 2009"),
		},
		{
			Key:   "IGNORED",
			Value: []byte("Ignore Me"),
		},
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ParseHeaders(&headers)
		if err != nil {
			b.Log(err)
		}
	}
}
