package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaMessageProducer interface {
	Initialize() error
	Destroy() error
	PublishMessage(message *Message) error
	BuildMessage(key string, body interface{}) (*Message, error)
}

type kafkaMessageProducer struct {
	configuration *KafkaProducerConfiguration
	writer        *kafka.Writer
}

func NewKafkaMessageProducer(configuration *KafkaProducerConfiguration) KafkaMessageProducer {
	return &kafkaMessageProducer{configuration: configuration}
}

func (k *kafkaMessageProducer) Initialize() error {
	c := k.configuration
	w := &kafka.Writer{
		Addr:                   kafka.TCP(c.BootstrapServers...),
		Topic:                  c.Topic,
		Balancer:               c.Balancer,
		BatchSize:              1, // important!
		MaxAttempts:            c.MaxAttempts,
		ReadTimeout:            time.Second * 5,
		WriteTimeout:           time.Second * 5,
		Async:                  false,               // important! Use this only if you don't care about guarantees of whether the messages were written to kafka.
		AllowAutoTopicCreation: c.TopicAutoCreation, // important - need to handler errors
		Logger:                 kafka.LoggerFunc(KafkaPrintLogger),
		ErrorLogger:            kafka.LoggerFunc(KafkaPrintLogger),
	}
	k.writer = w
	return nil
}

func (k *kafkaMessageProducer) Destroy() error {
	return k.writer.Close()
}

func (k *kafkaMessageProducer) PublishMessage(message *Message) error {
	if err := k.writer.WriteMessages(context.Background(), kafka.Message{
		Key:        []byte(message.Key),
		Value:      message.Data,
		WriterData: nil, // this can be handy with Completion function for writer
		Time:       message.Timestamp,
	}); err != nil {
		return err
	}
	return nil
}

func (k *kafkaMessageProducer) BuildMessage(key string, body interface{}) (*Message, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return &Message{
		Topic:     k.configuration.Topic,
		Timestamp: time.Now(),
		Key:       key,
		Data:      data,
	}, nil
}
