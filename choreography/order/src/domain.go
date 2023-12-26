package src

import (
	"errors"
	"infra/common/log"
	"infra/order"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type OrderService interface {
	CreateOrder(ctx echo.Context, order *order.OrderModel) (*order.OrderModel, error)
	UpdateOrder(ctx echo.Context, order *order.OrderModel) (*order.OrderModel, error)
}

func NewOrderService(logger *log.Logger, orderRepo order.OrderRepository) OrderService {
	return &orderService{
		logger:    logger,
		orderRepo: orderRepo,
	}
}

type orderService struct {
	logger    *log.Logger
	orderRepo order.OrderRepository
}

func (o *orderService) CreateOrder(ctx echo.Context, order *order.OrderModel) (*order.OrderModel, error) {
	o.logger.Info("Create new order")
	order.Id = uuid.New().String()
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	if err := o.orderRepo.CreateOrder(ctx, order); err != nil {
		o.logger.Error("Can not create new order ", err)
		return nil, err
	}
	o.logger.Info("Send event to kafka")
	//send order-update event
	return order, nil
}

func (o *orderService) UpdateOrder(ctx echo.Context, order *order.OrderModel) (*order.OrderModel, error) {
	// update order
	existingOrder, err := o.orderRepo.Query(ctx).ById(order.Id).Result()
	if err != nil {
		o.logger.Error("Can not get order by Id ", err)
		return nil, err
	}
	if existingOrder == nil {
		o.logger.Error("Order is not exist", err)
		return nil, errors.New("Invalid order")
	}
	existingOrder.Amount = order.Amount
	existingOrder.Detail = order.Detail
	existingOrder.Status = order.Status
	existingOrder.UserId = order.UserId
	existingOrder.UpdatedAt = time.Now()
	if err := o.orderRepo.UpdateOrder(ctx, order); err != nil {
		o.logger.Error("Can not update order", err)
		return nil, err
	}
	//send order-update event
	return order, nil
}
