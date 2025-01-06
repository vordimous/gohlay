package kafkautil

import (
	"testing"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/go-cmp/cmp"
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
			if got, err := ParseHeaders(&tt.args.headers); err != nil || !cmp.Equal(got, tt.want) {
				t.Errorf("ParseHeaders() = %v, want %v. err: %v", got, tt.want, err)
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
