package src

import (
	"infra/common/log"
	"infra/order"

	"github.com/labstack/echo/v4"
)

type OrderDomain interface {
	CreateOrder(ctx echo.Context) error
}

type orderDomain struct {
	logger    *log.Logger
	orderRepo order.OrderRepository
}

func NewOrderDomain(
	logger *log.Logger,
	orderRepo order.OrderRepository,
) OrderDomain {
	return &orderDomain{
		logger,
		orderRepo,
	}
}

func (s *orderDomain) CreateOrder(ctx echo.Context) error {
	s.logger.Info("creating order")

	return s.orderRepo.CreateOrder(ctx, &order.OrderModel{})
}
