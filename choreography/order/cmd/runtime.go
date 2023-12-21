package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	logger "infra/common/log"

	src "order-service/src"

	"infra/common/db"
	"infra/order"

	"gorm.io/gorm"
)

type runtime struct {
	appConf      *AppConfig
	logger       *logger.Logger
	db           *gorm.DB
	orderHandler src.OrderHandler
}

func NewRuntime() *runtime {
	rt := runtime{}
	var err error
	rt.logger = logger.New()

	if rt.appConf, err = BuildConfiguration(); err != nil {
		rt.logger.Error("Can not build config ", err)
	}

	rt.db, err = db.NewSQL(rt.appConf.DatabaseConfig)
	if err != nil {
		rt.logger.Error("Can not connect to database ", err)
	}

	orderRepository := order.NewOrderRepo()
	rt.orderHandler = src.NewOrderHandler(rt.logger, orderRepository)

	return &rt
}

func (rt *runtime) Serve() {
	api := NewApi(rt.appConf, rt.db, rt.logger, rt.orderHandler)
	registerSignalsHandler(api)
	api.Run()

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
