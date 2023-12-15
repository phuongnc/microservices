package protocol

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaMessageProducer struct {
	configuration *KafkaProducerConfiguration
	writer        *kafka.Writer
}

func NewKafkaMessageProducer(configuration *KafkaProducerConfiguration) *KafkaMessageProducer {
	return &KafkaMessageProducer{configuration: configuration}
}

func (k *KafkaMessageProducer) Initialize() error {
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

func (k *KafkaMessageProducer) Destroy() error {
	// TODO error handling and logging
	return k.writer.Close()
}

func (k *KafkaMessageProducer) PublishMessage(message *Message) error {
	if err := k.writer.WriteMessages(context.Background(), kafka.Message{
		Key:        []byte(message.Key),
		Value:      message.Data,
		WriterData: nil, // this can be handy with Completion function for writer
		Time:       message.Timestamp,
	}); err != nil {
		// TODO error handling and logging
		/*
			docs
			When the method returns an error, it may be of type kafka.WriteError to allow the caller
			to determine the status of each message.
		*/
		return err
	}
	return nil
}

func (k *KafkaMessageProducer) Name() string {
	return "KafkaMessageProducer"
}

func (k *KafkaMessageProducer) Topic() string {
	return k.configuration.Topic
}
