package src

import (
	"context"
	"net/http"

	"infra/common/log"
	"infra/order"

	"github.com/labstack/echo/v4"
)

type OrderHandler interface {
	RegisterEndpoints(echo *echo.Group)
	CreateOrder(c echo.Context) error
}

func NewOrderHandler(
	logger *log.Logger,
	orderService OrderService,
) OrderHandler {
	return &orderHandler{
		logger:       logger,
		orderService: orderService,
	}
}

type orderHandler struct {
	logger       *log.Logger
	orderService OrderService
}

func (rc *orderHandler) RegisterEndpoints(echo *echo.Group) {
	echo.POST("", rc.CreateOrder)
}

func (rc *orderHandler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, "Ok")
}

func (rc *orderHandler) CreateOrder(c echo.Context) error {
	req := &order.OrderDto{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Params")
	}

	ctx := context.WithValue(context.Background(), "db", c.Get("db"))
	model := order.MapOrderToModel(req)
	_, err := rc.orderService.CreateOrder(ctx, model)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusOK, nil)
}
