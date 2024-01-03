package src

import (
	"context"
	"net/http"

	"infra/common/log"
	"infra/order"

	"github.com/labstack/echo/v4"
)

type PaymentHandler interface {
	RegisterEndpoints(echo *echo.Group)
	paymentFailed(c echo.Context) error
	paymentSuccess(c echo.Context) error
}

func NewPaymentHandler(
	logger *log.Logger,
	paymentService PaymentService,
) PaymentHandler {
	return &paymentHandler{
		logger:         logger,
		paymentService: paymentService,
	}
}

type paymentHandler struct {
	logger         *log.Logger
	paymentService PaymentService
}

func (rc *paymentHandler) RegisterEndpoints(echo *echo.Group) {
	echo.POST("/failed", rc.paymentFailed)
	echo.POST("/success", rc.paymentSuccess)
	echo.GET("/orders/:orderId", rc.getOrder)
}

func (rc *paymentHandler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, "Ok")
}

func (rc *paymentHandler) paymentFailed(c echo.Context) error {
	req := &order.OrderDto{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Params")
	}

	model := &order.OrderModel{
		Id:            req.Id,
		SubStatus:     order.ORDER_PAYMENT_FAILED,
		FailureReason: req.FailureReason,
	}
	ctx := context.WithValue(context.Background(), "db", c.Get("db"))
	_, err := rc.paymentService.UpdateOrderPaymentStatus(ctx, model)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusOK, nil)
}

func (rc *paymentHandler) paymentSuccess(c echo.Context) error {
	req := &order.OrderDto{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Params")
	}

	ctx := context.WithValue(context.Background(), "db", c.Get("db"))
	model := &order.OrderModel{
		Id:        req.Id,
		SubStatus: order.ORDER_PAYMENT_PAID,
	}
	_, err := rc.paymentService.UpdateOrderPaymentStatus(ctx, model)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusOK, nil)
}

func (rc *paymentHandler) getOrder(c echo.Context) error {
	orderId := c.Param("orderId")
	if orderId == "" {
		return c.JSON(http.StatusBadRequest, "Invalid Params")
	}
	ctx := context.WithValue(context.Background(), "db", c.Get("db"))
	existingOrder, err := rc.paymentService.GetOrder(ctx, orderId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}
	if existingOrder == nil {
		return c.JSON(http.StatusNotFound, "Order is not exist")
	}
	return c.JSON(http.StatusOK, order.MapOrderFromModel(existingOrder))
}
