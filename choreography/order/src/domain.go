package src

import (
	"context"
	"encoding/json"
	"errors"
	"infra/common/kafka"
	"infra/common/log"
	"infra/order"
	"order-service/event"
	"time"

	"github.com/google/uuid"
)

type OrderService interface {
	OrderConsumeEvent(c context.Context, msg *kafka.Message) error
	CreateOrder(ctx context.Context, order *order.OrderModel) (*order.OrderModel, error)
	UpdateOrder(ctx context.Context, order *order.OrderModel) (*order.OrderModel, error)
}

func NewOrderService(logger *log.Logger, orderRepo order.OrderRepository, orderPublisher event.OrderPublisher) OrderService {
	return &orderService{
		logger:         logger,
		orderRepo:      orderRepo,
		orderPublisher: orderPublisher,
	}
}

type orderService struct {
	logger         *log.Logger
	orderRepo      order.OrderRepository
	orderPublisher event.OrderPublisher
}

func (o *orderService) OrderConsumeEvent(ctx context.Context, msg *kafka.Message) error {
	var msgOrder order.OrderModel
	if err := json.Unmarshal(msg.Data, &msgOrder); err != nil {
		o.logger.Error("Can not parse order from kafka ", err)
		return err
	}
	// check existing order
	existingOrder, err := o.orderRepo.Query(ctx).ById(msgOrder.Id).Result()
	if err != nil {
		o.logger.Error("Can not get order by Id ", err)
		return err
	}
	if existingOrder == nil {
		o.logger.Error("Order is not exist", err)
		return nil
	}

	existingOrder.SubStatus = msgOrder.SubStatus
	if existingOrder.SubStatus == order.ORDER_PAYMENT_PAID {
		existingOrder.Status = order.ORDER_PROCESSING
	} else if existingOrder.SubStatus == order.ORDER_PAYMENT_FAILED {
		existingOrder.Status = order.ORDER_FAILED
	} else if existingOrder.SubStatus == order.ORDER_KTCHENT_PREPARATION_FAILED {
		existingOrder.Status = order.ORDER_REFUNDING
	} else if existingOrder.SubStatus == order.ORDER_REFUNDED {
		existingOrder.Status = order.ORDER_FAILED
	}

	existingOrder.FailureReason = msgOrder.FailureReason
	existingOrder.UpdatedAt = time.Now()
	if err := o.orderRepo.UpdateOrder(ctx, existingOrder); err != nil {
		o.logger.Error("Can not update order", err)
		return err
	}
	// incase need refund payment, publish order update event
	if existingOrder.Status == order.ORDER_REFUNDING {
		o.logger.Info("Send event to kafka")
		return o.orderPublisher.PublishOrderEvent(existingOrder)
	}
	return nil
}

func (o *orderService) CreateOrder(ctx context.Context, obj *order.OrderModel) (*order.OrderModel, error) {
	o.logger.Info("Create new order")
	obj.Id = uuid.New().String()
	obj.CreatedAt = time.Now()
	obj.UpdatedAt = time.Now()
	obj.Status = order.ORDER_CREATED
	if err := o.orderRepo.CreateOrder(ctx, obj); err != nil {
		o.logger.Error("Can not create new order ", err)
		return nil, err
	}
	o.logger.Info("Send event to kafka")
	err := o.orderPublisher.PublishOrderEvent(obj)
	if err != nil {
		return nil, err

	}
	return obj, nil
}

func (o *orderService) UpdateOrder(ctx context.Context, obj *order.OrderModel) (*order.OrderModel, error) {
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
	existingOrder.Amount = obj.Amount
	existingOrder.Detail = obj.Detail
	existingOrder.Status = obj.Status
	existingOrder.UserId = obj.UserId
	existingOrder.UpdatedAt = time.Now()
	if err := o.orderRepo.UpdateOrder(ctx, existingOrder); err != nil {
		o.logger.Error("Can not update order", err)
		return nil, err
	}
	return existingOrder, nil
}
