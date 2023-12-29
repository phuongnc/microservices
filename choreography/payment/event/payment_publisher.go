package event

import (
	"infra/common/kafka"
	"infra/order"

	"github.com/google/uuid"
)

type PaymentPublisher interface {
	PublishPaymentEvent(order *order.OrderModel) error
	Destroy() error
}

func NewPaymentPublisher(producer kafka.KafkaMessageProducer) PaymentPublisher {
	return &paymentPublisher{
		producer: producer,
	}
}

type paymentPublisher struct {
	producer kafka.KafkaMessageProducer
}

func (o *paymentPublisher) PublishPaymentEvent(order *order.OrderModel) error {
	paymentEvent, _ := o.producer.BuildMessage(uuid.New().String(), order)
	return o.producer.PublishMessage(paymentEvent)
}

func (o *paymentPublisher) Destroy() error {
	return o.producer.Destroy()
}
