package kafka

import (
	"context"
	"fmt"
	"io"

	"github.com/segmentio/kafka-go"
)

type ConsumerHandler func(ctx context.Context, message *Message) error

type KafkaMessageConsumer interface {
	StartConsumer(handler ConsumerHandler)
	Destroy() error
}

type kafkaMessageConsumer struct {
	configuration *KafkaConsumerConfiguration
	context       context.Context
	reader        *kafka.Reader
	stopSync      chan bool
}

func NewKafkaMessageConsumer(ctx context.Context, configuration *KafkaConsumerConfiguration) KafkaMessageConsumer {
	consumer := &kafkaMessageConsumer{configuration: configuration, context: ctx}
	consumer.initialize()
	return consumer
}

func (k *kafkaMessageConsumer) initialize() {
	c := k.configuration
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        c.BootstrapServers,
		GroupID:        c.GroupID,
		GroupTopics:    c.Topics,
		MaxBytes:       10e6,
		CommitInterval: 1,
		//Logger:         kafka.LoggerFunc(KafkaPrintLogger),
		ErrorLogger: kafka.LoggerFunc(KafkaPrintLogger),
	})
	k.reader = reader
}

func (k *kafkaMessageConsumer) Destroy() error {
	//stop sync
	if k.stopSync != nil {
		fmt.Println("stop to receive messages")
		k.stopSync <- true
		close(k.stopSync)
		k.stopSync = nil
	}
	//close reader
	return k.reader.Close()
}

func (k *kafkaMessageConsumer) StartConsumer(handler ConsumerHandler) {
	if k.configuration.AutoCommit {
		k.startConsumerAutoCommit(handler)
	} else {
		k.startConsumerManualCommit(handler)
	}
}

func (k *kafkaMessageConsumer) startConsumerAutoCommit(handler ConsumerHandler) {
	for {
		kMsg, err := k.reader.ReadMessage(k.context)
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
		if pErr := handler(k.context, msg); pErr != nil {
			// TODO handle errors ? - check specific errors
			// when error occurs I'm not doing the commit! as I need message do be redelivered and processed!
			continue
		}
	}
}

func (k *kafkaMessageConsumer) startConsumerManualCommit(handlerFunc ConsumerHandler) {
	if k.stopSync == nil {
		k.stopSync = make(chan bool, 1)
	}
	go func(handlerFunc ConsumerHandler, stopSync <-chan bool) {
		for {
			select {
			case <-stopSync:
				return
			default:
				kMsg, err := k.reader.FetchMessage(k.context)
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
				if pErr := handlerFunc(k.context, msg); pErr != nil {
					// TODO handle errors ? - check specific errors
					// when error occurs I'm not doing the commit! as I need message do be redelivered and processed!
					continue
				}
				k.reader.CommitMessages(k.context, kMsg)

			}
		}
	}(handlerFunc, k.stopSync)
}

// func (k *kafkaMessageConsumer) Name() string {
// 	return "KafkaMessageConsumer"
// }

// func (k *kafkaMessageConsumer) Topic() string {
// 	return k.configuration.Topic
// }
