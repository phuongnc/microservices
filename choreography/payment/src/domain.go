package src

import (
	"context"
	"encoding/json"
	"errors"
	"infra/common/kafka"
	"infra/common/log"
	"infra/order"
	"payment-service/event"
	"time"
)

type PaymentService interface {
	PaymentConsumeEvent(c context.Context, msg *kafka.Message) error
	UpdateOrderPaymentStatus(ctx context.Context, order *order.OrderModel) (*order.OrderModel, error)
	GetOrder(ctx context.Context, orderId string) (*order.OrderModel, error)
}

func NewPaymentService(logger *log.Logger, orderRepo order.OrderRepository, paymentPublisher event.PaymentPublisher) PaymentService {
	return &paymentService{
		logger:           logger,
		orderRepo:        orderRepo,
		paymentPublisher: paymentPublisher,
	}
}

type paymentService struct {
	logger           *log.Logger
	orderRepo        order.OrderRepository
	paymentPublisher event.PaymentPublisher
}

func (o *paymentService) PaymentConsumeEvent(ctx context.Context, msg *kafka.Message) error {
	var msgOrder order.OrderModel
	if err := json.Unmarshal(msg.Data, &msgOrder); err != nil {
		o.logger.Error("Can not parse order from kafka ", err)
		return err
	}

	if msgOrder.Status == order.ORDER_CREATED {
		if err := o.orderRepo.CreateOrder(ctx, &msgOrder); err != nil {
			o.logger.Error("Can not save order ", err)
			return err
		}
	} else if msgOrder.Status == order.ORDER_REFUNDING {
		//update order to refunded
		msgOrder.SubStatus = order.ORDER_REFUNDED
		_, err := o.UpdateOrderPaymentStatus(ctx, &msgOrder)
		if err != nil {
			o.logger.Error("Can not update order payment ", err)
			return err
		}
	}
	return nil
}

func (o *paymentService) UpdateOrderPaymentStatus(ctx context.Context, obj *order.OrderModel) (*order.OrderModel, error) {
	// update order
	existingOrder, err := o.orderRepo.Query(ctx).ById(obj.Id).Result()
	if err != nil {
		o.logger.Error("Can not get order by Id ", err)
		return nil, err
	}
	if existingOrder == nil {
		o.logger.Error("Order is not exist", err)
		return nil, errors.New("Invalid order")
	}

	if existingOrder.Status != order.ORDER_CREATED || existingOrder.Status != order.ORDER_REFUNDING {
		return nil, errors.New("Order has been processed")
	}

	existingOrder.Status = obj.Status
	existingOrder.SubStatus = obj.SubStatus
	existingOrder.FailureReason = obj.FailureReason
	existingOrder.UpdatedAt = time.Now()
	if err := o.orderRepo.UpdateOrder(ctx, existingOrder); err != nil {
		o.logger.Error("Can not update order", err)
		return nil, err
	}
	// publish message to payment event
	err = o.paymentPublisher.PublishPaymentEvent(existingOrder)
	if err != nil {
		o.logger.Error("Can not publish payment event ", err)
		return nil, err
	}
	return existingOrder, nil
}

func (o *paymentService) GetOrder(ctx context.Context, orderId string) (*order.OrderModel, error) {
	return o.orderRepo.Query(ctx).ById(orderId).Result()
}
