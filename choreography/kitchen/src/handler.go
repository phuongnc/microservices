package src

import (
	"context"
	"net/http"

	"infra/common/log"
	"infra/order"

	"github.com/labstack/echo/v4"
)

type KitchenHandler interface {
	RegisterEndpoints(echo *echo.Group)
	CreateOrder(c echo.Context) error
}

func NewKitchenHandler(
	logger *log.Logger,
	kitchenService KitchenService,
) KitchenHandler {
	return &kitchenHandler{
		logger:         logger,
		kitchenService: kitchenService,
	}
}

type kitchenHandler struct {
	logger         *log.Logger
	kitchenService KitchenService
}

func (rc *kitchenHandler) RegisterEndpoints(echo *echo.Group) {
	echo.POST("", rc.CreateOrder)
}

func (rc *kitchenHandler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, "Ok")
}

func (rc *kitchenHandler) CreateOrder(c echo.Context) error {
	req := &order.OrderDto{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Params")
	}

	ctx := context.WithValue(context.Background(), "db", c.Get("db"))
	model := order.MapOrderToModel(req)
	_, err := rc.kitchenService.CreateOrder(ctx, model)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusOK, nil)
}
