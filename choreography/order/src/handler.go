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
	createOrder(c echo.Context) error
	getOrder(c echo.Context) error
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
	echo.POST("", rc.createOrder)
	echo.GET("/:orderId", rc.getOrder)
}

func (rc *orderHandler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, "Ok")
}

func (rc *orderHandler) createOrder(c echo.Context) error {
	req := &order.OrderDto{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Params")
	}

	ctx := context.WithValue(context.Background(), "db", c.Get("db"))
	model := order.MapOrderToModel(req)
	newOrder, err := rc.orderService.CreateOrder(ctx, model)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusOK, order.MapOrderFromModel(newOrder))
}

func (rc *orderHandler) getOrder(c echo.Context) error {
	orderId := c.Param("orderId")
	if orderId == "" {
		return c.JSON(http.StatusBadRequest, "Invalid Params")
	}
	ctx := context.WithValue(context.Background(), "db", c.Get("db"))
	existingOrder, err := rc.orderService.GetOrder(ctx, orderId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}
	if existingOrder == nil {
		return c.JSON(http.StatusNotFound, "Order is not exist")
	}
	return c.JSON(http.StatusOK, order.MapOrderFromModel(existingOrder))
}
