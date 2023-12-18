package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	logger "infra/log"

	src "order-service/src"

	"gorm.io/gorm"
)

type runtime struct {
	appConf      *src.AppConfig
	logger       *logger.Logger
	db           *gorm.DB
	orderHandler *src.OrderService
}

func NewRuntime() *runtime {
	rt := runtime{}
	var err error

	if rt.appConf, err = src.BuildConfiguration(); err != nil {
		fmt.Sprintf("can't load application configuration: %v", err)
	}

	// rt.db, err = db.NewSQL(rt.appConf.DatabaseConfig)
	// if err != nil {
	// 	fmt.Sprintf("cannot connect to db: %v", err)
	// }

	rt.logger = logger.New()

	orderDomain := src.NewOrderDomain(rt.logger)
	rt.orderHandler = src.NewOrderHandler(rt.logger, orderDomain)

	// service, err := NewService(rt.logger)
	// if err != nil {
	// 	rt.logger.Error("creating new service instance", err)
	// }
	// rt.service = service

	return &rt
}

func (rt *runtime) Serve() {
	// api := NewApi(rt.appConf, rt.db, rt.logger, rt.service)
	// registerSignalsHandler(api)
	// api.Run()

	// run kafka
	//kafka.TestKafka()

}

func registerSignalsHandler(api *Api) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Printf("Received termination signal: [%s], stopping app", sig)
		api.Stop()
	}()
}
