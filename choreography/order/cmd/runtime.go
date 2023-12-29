package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"infra/common/kafka"
	logger "infra/common/log"

	"order-service/event"
	src "order-service/src"

	"infra/common/db"
	"infra/order"

	"gorm.io/gorm"
)

type runtime struct {
	appConf        *AppConfig
	logger         *logger.Logger
	db             *gorm.DB
	orderHandler   src.OrderHandler
	orderPublisher event.OrderPublisher
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
	rt.migrateDB()

	orderRepository := order.NewOrderRepo()
	orderService := src.NewOrderService(rt.logger, orderRepository, rt.orderPublisher)
	rt.orderHandler = src.NewOrderHandler(rt.logger, orderService)

	//setup kafka publisher
	kafkaConfig := &kafka.KafkaProducerConfiguration{
		BootstrapServers:  rt.appConf.KafkaConfig.BootstrapServers,
		Topic:             rt.appConf.KafkaConfig.OrderEventTopic,
		TopicAutoCreation: true,
	}
	producer := kafka.NewKafkaMessageProducer(kafkaConfig)
	rt.orderPublisher = event.NewOrderPublisher(producer)
	//setup kafka consumer
	kafkaConsumerConfig := &kafka.KafkaConsumerConfiguration{
		BootstrapServers: rt.appConf.KafkaConfig.BootstrapServers,
		Topics:           []string{rt.appConf.KafkaConfig.PaymentEventTopic, rt.appConf.KafkaConfig.kitchenEventTopic},
		GroupID:          rt.appConf.KafkaConfig.OrderGroup,
		AutoCommit:       false,
	}
	ctx := context.WithValue(context.Background(), "db", rt.db)
	consumer := kafka.NewKafkaMessageConsumer(ctx, kafkaConsumerConfig)
	consumer.StartConsumer(orderService.OrderConsumeEvent)

	return &rt
}

func (rt *runtime) Serve() {
	api := NewApi(rt.appConf, rt.db, rt.logger, rt.orderHandler)
	rt.registerSignalsHandler(api)
	api.Run()
	rt.logger.Info("call here")

}

func (rt *runtime) migrateDB() {
	rt.db.Table("order").AutoMigrate(&order.OrderEntity{})
}

func (rt *runtime) registerSignalsHandler(api *Api) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Printf("Received termination signal: [%s], stopping app", sig)
		api.Stop()
		rt.orderPublisher.Destroy()
	}()
}
