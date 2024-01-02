package cmd

import (
	"context"
	"fmt"

	"infra/common/log"
	"infra/common/middleware"
	"kitchen-service/src"

	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Api struct {
	echo   *echo.Echo
	config *AppConfig
	logger *log.Logger
}

func NewApi(config *AppConfig, gormDB *gorm.DB, logger *log.Logger, kitchenHandler src.KitchenHandler) *Api {
	a := &Api{
		config: config,
		logger: logger,
	}
	a.echo = createHttpServer(gormDB, kitchenHandler, a.logger)
	return a
}

func (a *Api) Run() error {
	a.logger.Info("starting order service")
	return a.echo.Start(fmt.Sprintf(":%s", a.config.HttpServerConfig.Port))
}

func (a *Api) Stop() {
	a.logger.Info("Shutdown order service")
	_ = a.echo.Shutdown(context.Background())
}

func createHttpServer(gormDB *gorm.DB, kitchenHandler src.KitchenHandler, log *log.Logger) *echo.Echo {
	e := echo.New()
	e.GET("/health", healthCheck)
	path := e.Group("/kitchens")

	e.Use(middleware.HttpDb(gormDB))
	e.Use(middleware.Cors())
	kitchenHandler.RegisterEndpoints(path)
	return e
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, "Running")
}
