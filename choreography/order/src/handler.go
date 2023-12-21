package src

import (
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
	orderRepo order.OrderRepository,
) OrderHandler {
	return &orderHandler{
		logger:    logger,
		orderRepo: orderRepo,
	}
}

type orderHandler struct {
	logger    *log.Logger
	orderRepo order.OrderRepository
}

func (rc *orderHandler) RegisterEndpoints(echo *echo.Group) {
	echo.POST("/create-order", rc.CreateOrder)
}

func (rc *orderHandler) HealthCheck(c echo.Context) error {

	return c.JSON(http.StatusOK, "Ok")
}

func (rc *orderHandler) CreateOrder(c echo.Context) error {
	req := &order.OrderDto{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Params")
	}
	model := order.MapOrderToModel(req)
	if err := rc.orderRepo.CreateOrder(c, model); err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusOK, nil)
}
