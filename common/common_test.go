package common

import (
	"testing"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/spf13/viper"
)

func TestFmtMessageId(t *testing.T) {
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
			if got := FmtMessageId(tt.args.offset, tt.args.delay); got != tt.want {
				t.Errorf("FmtMessageId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFmtKafkaGroup(t *testing.T) {
	type args struct {
		reason string
		topic  string
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
				reason: "test",
				topic:  "aTopic",
			},
			deadline: 12345,
			want:     "group.id=gohlay_test:aTopic:12345",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set("deadline", tt.deadline)
			if got := FmtKafkaGroup(tt.args.reason, tt.args.topic); got != tt.want {
				t.Errorf("FmtKafkaGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
