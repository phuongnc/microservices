package event

import (
	"infra/common/kafka"
	"infra/order"

	"github.com/google/uuid"
)

type KitchenPublisher interface {
	PublishKitchenEvent(order *order.OrderModel) error
	Destroy() error
}

func NewKitchenPublisher(producer kafka.KafkaMessageProducer) KitchenPublisher {
	return &kitchenPublisher{
		producer: producer,
	}
}

type kitchenPublisher struct {
	producer kafka.KafkaMessageProducer
}

func (o *kitchenPublisher) PublishKitchenEvent(order *order.OrderModel) error {
	paymentEvent, _ := o.producer.BuildMessage(uuid.New().String(), order)
	return o.producer.PublishMessage(paymentEvent)
}

func (o *kitchenPublisher) Destroy() error {
	return o.producer.Destroy()
}
