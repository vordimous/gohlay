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
		wantDeliveryTime   int64
		wantIsDelivered    bool
		wantIsDeliveredKey string
		wantGohlayed       bool
	}{
		{
			name: "Gohlayed message",
			args: args{
				headers: []kafka.Header{
					{
						Key:   "GOHLAY",
						Value: []byte("Tue Nov 03 15:28:57 UTC 1989"),
					},
				},
			},
			wantDeliveryTime:   626110137000,
			wantIsDelivered:    false,
			wantIsDeliveredKey: "",
			wantGohlayed:       true,
		},
		{
			name: "Gohlayed message is Delivered",
			args: args{
				headers: []kafka.Header{
					{
						Key:   "GOHLAY_DELIVERED",
						Value: []byte("0-12345"),
					},
				},
			},
			wantDeliveryTime:   0,
			wantIsDelivered:    true,
			wantIsDeliveredKey: "0-12345",
			wantGohlayed:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDeliveryTime, gotIsDelivered, gotIsDeliveredKey, gotGohlayed := ParseHeaders(tt.args.headers)
			if gotDeliveryTime != tt.wantDeliveryTime {
				t.Errorf("ParseHeaders() gotDeliveryTime = %v, want %v", gotDeliveryTime, tt.wantDeliveryTime)
			}
			if gotIsDelivered != tt.wantIsDelivered {
				t.Errorf("ParseHeaders() gotIsDelivered = %v, want %v", gotIsDelivered, tt.wantIsDelivered)
			}
			if gotIsDeliveredKey != tt.wantIsDeliveredKey {
				t.Errorf("ParseHeaders() gotIsDeliveredKey = %v, want %v", gotIsDeliveredKey, tt.wantIsDeliveredKey)
			}
			if gotGohlayed != tt.wantGohlayed {
				t.Errorf("ParseHeaders() gotGohlayed = %v, want %v", gotGohlayed, tt.wantGohlayed)
			}
		})
	}
}
