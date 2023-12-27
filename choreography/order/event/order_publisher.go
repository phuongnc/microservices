package event

import (
	"encoding/json"
	"fmt"
	"infra/common/kafka/protocol"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer interface {
}

type kafkaProducer struct {
	bootstrapServer []string
	topic           string
	producer        *protocol.KafkaMessageProducer
}

func (k *kafkaProducer) SetupKafkaProducer() {
	producerConfig := &protocol.KafkaProducerConfiguration{
		BootstrapServers:  k.bootstrapServer,
		Topic:             k.topic,
		MaxAttempts:       10,
		Balancer:          &kafka.Hash{},
		TopicAutoCreation: true,
	}
	k.producer = protocol.NewKafkaMessageProducer(producerConfig)
}

func (k *kafkaProducer) BuildMessage(key string, body interface{}) *protocol.Message {
	fmt.Println("Message: ", body)
	data, err := json.Marshal(body)

	if err != nil {
		panic("Marshalling error")
	}

	return &protocol.Message{
		Topic:     k.topic,
		Timestamp: time.Now(),
		Key:       key,
		Data:      data,
	}
}

func (k *kafkaProducer) PublishMessage(msg *protocol.Message) error {
	return k.producer.PublishMessage(msg)
}
