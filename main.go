package main

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	fmt.Println("sex")
	a, b := kafka.LibraryVersion()
	fmt.Printf("Kafka Lib Version: %#v | %#v\n", a, b)
}
