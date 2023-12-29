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
	CreateOrder(c echo.Context) error
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
	echo.POST("", rc.CreateOrder)
}

func (rc *paymentHandler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, "Ok")
}

func (rc *paymentHandler) CreateOrder(c echo.Context) error {
	req := &order.OrderDto{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Params")
	}

	ctx := context.WithValue(context.Background(), "db", c.Get("db"))
	model := order.MapOrderToModel(req)
	_, err := rc.paymentService.CreateOrder(ctx, model)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusOK, nil)
}
