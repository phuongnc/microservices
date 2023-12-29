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

	"github.com/google/uuid"
)

type PaymentService interface {
	PaymentConsumeEvent(c context.Context, msg *kafka.Message) error
	CreateOrder(ctx context.Context, order *order.OrderModel) (*order.OrderModel, error)
	UpdateOrder(ctx context.Context, order *order.OrderModel) (*order.OrderModel, error)
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
	existingOrder.Status = msgOrder.Status
	if msgOrder.Status == order.ORDER_PREPARATION_FAILED {
		existingOrder.Status = order.ORDER_REFUNDING
	}
	existingOrder.UpdatedAt = time.Now()
	if err := o.orderRepo.UpdateOrder(ctx, existingOrder); err != nil {
		o.logger.Error("Can not update order", err)
		return err
	}
	// incase need refund payment, publish order update event
	if existingOrder.Status == order.ORDER_REFUNDING {
		o.logger.Info("Send event to kafka")
		return o.paymentPublisher.PublishPaymentEvent(existingOrder)
	}
	return nil
}

func (o *paymentService) CreateOrder(ctx context.Context, obj *order.OrderModel) (*order.OrderModel, error) {
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
	err := o.paymentPublisher.PublishPaymentEvent(obj)
	if err != nil {
		return nil, err

	}
	return obj, nil
}

func (o *paymentService) UpdateOrder(ctx context.Context, obj *order.OrderModel) (*order.OrderModel, error) {
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
