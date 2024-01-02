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
	PaymentFailed(c echo.Context) error
	PaymentSuccess(c echo.Context) error
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
	echo.POST("/failed", rc.PaymentFailed)
	echo.POST("/success", rc.PaymentSuccess)
}

func (rc *paymentHandler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, "Ok")
}

func (rc *paymentHandler) PaymentFailed(c echo.Context) error {
	req := &order.OrderDto{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Params")
	}

	req.SubStatus = order.ORDER_PAYMENT_FAILED
	ctx := context.WithValue(context.Background(), "db", c.Get("db"))
	model := order.MapOrderToModel(req)
	_, err := rc.paymentService.UpdateOrderPaymentStatus(ctx, model)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusOK, nil)
}

func (rc *paymentHandler) PaymentSuccess(c echo.Context) error {
	req := &order.OrderDto{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Params")
	}

	req.SubStatus = order.ORDER_PAYMENT_PAID
	ctx := context.WithValue(context.Background(), "db", c.Get("db"))
	model := order.MapOrderToModel(req)
	_, err := rc.paymentService.UpdateOrderPaymentStatus(ctx, model)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusOK, nil)
}
