package kafka

import (
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

/*
	Just temporarily here - some common part
*/

type Message struct {
	Topic     string
	Timestamp time.Time
	Key       string
	Data      []byte
}

func NewMessage(topic string, timestamp time.Time, key string, data []byte) *Message {
	return &Message{
		Topic:     topic,
		Timestamp: timestamp,
		Key:       key,
		Data:      data,
	}
}

type MessageProducer interface {
	Initialize() error
	Destroy() error
	PublishMessage(message *Message) error
	Topic() string
}

type MessageConsumer interface {
	Initialize() error
	Destroy() error
	StartConsumer(handler func(message *Message) error)
	Topic() string
}

type KafkaProducerConfiguration struct {
	BootstrapServers  []string
	Topic             string
	MaxAttempts       int
	Balancer          kafka.Balancer
	TopicAutoCreation bool
}

type KafkaConsumerConfiguration struct {
	BootstrapServers []string
	Topic            string
	GroupID          string
	AutoCommit       bool
}

func KafkaPrintLogger(msg string, a ...interface{}) {
	fmt.Printf(msg, a...)
	fmt.Println()
}
