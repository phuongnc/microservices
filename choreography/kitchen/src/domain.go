package src

import (
	"context"
	"encoding/json"
	"errors"
	"infra/common/kafka"
	"infra/common/log"
	"infra/order"
	"kitchen-service/event"
	"time"
)

type KitchenService interface {
	PaymentConsumeEvent(c context.Context, msg *kafka.Message) error
	UpdateOrderKitchenStatus(ctx context.Context, obj *order.OrderModel) (*order.OrderModel, error)
	GetOrder(ctx context.Context, orderId string) (*order.OrderModel, error)
}

func NewKitchenService(logger *log.Logger, orderRepo order.OrderRepository, kitchenPublisher event.KitchenPublisher) KitchenService {
	return &kitchenService{
		logger:           logger,
		orderRepo:        orderRepo,
		kitchenPublisher: kitchenPublisher,
	}
}

type kitchenService struct {
	logger           *log.Logger
	orderRepo        order.OrderRepository
	kitchenPublisher event.KitchenPublisher
}

func (o *kitchenService) PaymentConsumeEvent(ctx context.Context, msg *kafka.Message) error {
	var msgOrder order.OrderModel
	if err := json.Unmarshal(msg.Data, &msgOrder); err != nil {
		o.logger.Error("Can not parse order from kafka ", err)
		return err
	}
	if msgOrder.SubStatus == order.ORDER_PAYMENT_PAID {
		if err := o.orderRepo.CreateOrder(ctx, &msgOrder); err != nil {
			o.logger.Error("Can not save order ", err)
			return err
		}
	}
	return nil
}

func (o *kitchenService) UpdateOrderKitchenStatus(ctx context.Context, obj *order.OrderModel) (*order.OrderModel, error) {
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
	existingOrder.SubStatus = obj.SubStatus
	existingOrder.FailureReason = obj.FailureReason
	existingOrder.UpdatedAt = time.Now()
	if err := o.orderRepo.UpdateOrder(ctx, existingOrder); err != nil {
		o.logger.Error("Can not update order", err)
		return nil, err
	}
	// publish message to payment event
	err = o.kitchenPublisher.PublishKitchenEvent(existingOrder)
	if err != nil {
		o.logger.Error("Can not publish kitchen event ", err)
		return nil, err
	}
	return existingOrder, nil
}

func (o *kitchenService) GetOrder(ctx context.Context, orderId string) (*order.OrderModel, error) {
	return o.orderRepo.Query(ctx).ById(orderId).Result()
}
