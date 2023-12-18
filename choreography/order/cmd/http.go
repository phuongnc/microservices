package cmd

import (
	"context"
	"fmt"

	"infra/log"
	"infra/middleware"
	"order-service/src"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Handles struct {
	echo *echo.Echo
}

type Api struct {
	handles *Handles
	config  *src.AppConfig
	logger  *log.Logger
}

func NewApi(config *src.AppConfig, db *gorm.DB, logger *log.Logger, orderHandler src.OrderHandler) *Api {
	a := &Api{
		config:  config,
		handles: &Handles{},
		logger:  logger,
	}

	a.handles.echo = createHttpServer(db, orderHandler, a.logger)
	return a
}

func (a *Api) Run() error {
	a.logger.Info("starting")
	return a.handles.echo.Start(fmt.Sprintf(":%s", a.config.HttpServerConfig.Port))
}

func (a *Api) Stop() {
	_ = a.handles.echo.Shutdown(context.Background())
}

func createHttpServer(db *gorm.DB, orderHandler src.OrderHandler, log *log.Logger) *echo.Echo {
	e := echo.New()
	v1 := e.Group("/v1")
	e.Use(middleware.HttpDb(db))
	e.Use(middleware.Cors())

	// e.HTTPErrorHandler = func(err error, c echo.Context) {
	// 	if c.Response().Committed {
	// 		return
	// 	}

	// 	var apiError *rest.Error
	// 	if errors.As(err, &apiError) {
	// 		if err := c.JSON(apiError.HttpCode, apiError); err != nil {
	// 			log.Error("echo-error-handler", err)
	// 		}
	// 		return
	// 	}

	// 	e.DefaultHTTPErrorHandler(err, c)
	// }
	orderHandler.RegisterEndpoints(v1)
	return e
}
