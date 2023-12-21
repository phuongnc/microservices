package src

import (
	"net/http"

	"infra/common/log"

	"github.com/labstack/echo/v4"
)

type OrderHandler interface {
	RegisterEndpoints(echo *echo.Group)
	CreateOrder(c echo.Context) error
}

func NewOrderHandler(
	logger *log.Logger,
	service OrderDomain,
) OrderHandler {
	return &orderHandler{
		logger:  logger,
		service: service,
	}
}

type orderHandler struct {
	logger  *log.Logger
	service OrderDomain
}

func (rc *orderHandler) RegisterEndpoints(echo *echo.Group) {
	echo.POST("/create-order", rc.CreateOrder)
}

func (rc *orderHandler) HealthCheck(c echo.Context) error {

	return c.JSON(http.StatusOK, "Ok")
}

func (rc *orderHandler) CreateOrder(c echo.Context) error {

	return c.JSON(http.StatusOK, nil)
}
