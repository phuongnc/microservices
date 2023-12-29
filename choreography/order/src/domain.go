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
	"github.com/labstack/echo/v4"
)

type OrderService interface {
	OrderConsumeEvent(c context.Context, msg *kafka.Message) error
	CreateOrder(ctx context.Context, order *order.OrderModel) (*order.OrderModel, error)
	UpdateOrder(ctx echo.Context, order *order.OrderModel) (*order.OrderModel, error)
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
	var order order.OrderModel
	if err := json.Unmarshal(msg.Data, &order); err != nil {
		o.logger.Error("Can not parse order from kafka ", err)
		return err
	}

	// zexistingOrder, err := o.orderRepo.Query(ctx).ById(obj.Id).Result()
	// if err != nil {
	// 	o.logger.Error("Can not get order by Id ", err)
	// 	return nil, err
	// }
	// if existingOrder == nil {
	// 	o.logger.Error("Order is not exist", err)
	// 	return nil, errors.New("Invalid order")
	// }
	// existingOrder.Amount = obj.Amount
	// existingOrder.Detail = obj.Detail
	// existingOrder.Status = obj.Status
	// existingOrder.UserId = obj.UserId
	// existingOrder.UpdatedAt = time.Now()
	// if err := o.orderRepo.UpdateOrder(ctx, existingOrder); err != nil {
	// 	o.logger.Error("Can not update order", err)
	// 	return nil, err
	// }

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

func (o *orderService) UpdateOrder(ctx echo.Context, obj *order.OrderModel) (*order.OrderModel, error) {
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
	if err := o.orderRepo.UpdateOrder(ctx.Request().Context(), existingOrder); err != nil {
		o.logger.Error("Can not update order", err)
		return nil, err
	}
	//send order-update event
	return existingOrder, nil
}
