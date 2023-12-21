package protocol

import (
	"context"
	"fmt"
	"io"

	"github.com/segmentio/kafka-go"
)

type KafkaMessageConsumer struct {
	configuration *KafkaConsumerConfiguration
	reader        *kafka.Reader
}

func NewKafkaMessageConsumer(configuration *KafkaConsumerConfiguration) *KafkaMessageConsumer {
	return &KafkaMessageConsumer{configuration: configuration}
}

func (k *KafkaMessageConsumer) Initialize() error {
	c := k.configuration
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: c.BootstrapServers,
		GroupID: c.GroupID,
		Topic:   c.Topic,
		//QueueCapacity:     0,
		MaxBytes: 10e6, // 10MB
		//MaxWait:           0,
		CommitInterval: 1, // batches commits in 1 second interval
		//SessionTimeout:    0,
		//JoinGroupBackoff:  0,
		//RetentionTime:     0,
		//StartOffset:       0,
		//ReadBackoffMin:    0,
		//ReadBackoffMax:    0,
		Logger:      kafka.LoggerFunc(KafkaPrintLogger),
		ErrorLogger: kafka.LoggerFunc(KafkaPrintLogger),
	})
	k.reader = reader
	return nil
}

func (k *KafkaMessageConsumer) Destroy() error {
	// TODO error handling and logging
	return k.reader.Close()
}

func (k *KafkaMessageConsumer) StartConsumer(handler func(message *Message) error) {
	if k.configuration.AutoCommit {
		k.startConsumerAutoCommit(handler)
	} else {
		k.startConsumerManualCommit(handler)
	}
}

func (k *KafkaMessageConsumer) startConsumerAutoCommit(handler func(message *Message) error) {
	for {
		kMsg, err := k.reader.ReadMessage(context.Background())
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Kafka consumer is closed")
			}
			// TODO handle errors ? - check specific errors
			fmt.Printf("Error in kafka consumer: %s", err)
			break
		}
		msg := &Message{
			Topic:     kMsg.Topic,
			Timestamp: kMsg.Time,
			Key:       string(kMsg.Key),
			Data:      kMsg.Value,
		}
		if pErr := handler(msg); pErr != nil {
			// TODO handle errors ? - check specific errors
			// when error occurs I'm not doing the commit! as I need message do be redelivered and processed!
			continue
		}
	}
}

func (k *KafkaMessageConsumer) startConsumerManualCommit(handler func(message *Message) error) {
	for {
		kMsg, err := k.reader.FetchMessage(context.Background())
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Kafka consumer is closed")
			}
			// TODO handle errors ? - check specific errors
			fmt.Printf("Error in kafka consumer: %s", err)
			break
		}
		msg := &Message{
			Topic:     kMsg.Topic,
			Timestamp: kMsg.Time,
			Key:       string(kMsg.Key),
			Data:      kMsg.Value,
		}
		if pErr := handler(msg); pErr != nil {
			// TODO handle errors ? - check specific errors
			// when error occurs I'm not doing the commit! as I need message do be redelivered and processed!
			continue
		}
		/*
			Because kafka consumer groups track a single offset per partition, the highest message offset
			passed to CommitMessages will cause all previous messages to be committed!!!
			This needs to be verified
		*/
		k.reader.CommitMessages(context.Background(), kMsg)
	}
}

func (k *KafkaMessageConsumer) Name() string {
	return "KafkaMessageConsumer"
}

func (k *KafkaMessageConsumer) Topic() string {
	return k.configuration.Topic
}
