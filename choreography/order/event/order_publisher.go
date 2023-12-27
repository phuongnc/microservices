package event

import (
	"infra/common/kafka"
	"infra/order"

	"github.com/google/uuid"
)

type OrderPublisher interface {
	PublishOrderEvent(order *order.OrderModel) error
	Destroy() error
}

func NewOrderPublisher(producer kafka.KafkaMessageProducer) OrderPublisher {
	return &orderPublisher{
		producer: producer,
	}
}

type orderPublisher struct {
	producer kafka.KafkaMessageProducer
}

func (o *orderPublisher) PublishOrderEvent(order *order.OrderModel) error {
	orderEvent, _ := o.producer.BuildMessage(uuid.New().String(), order)
	return o.producer.PublishMessage(orderEvent)
}

func (o *orderPublisher) Destroy() error {
	return o.producer.Destroy()
}
