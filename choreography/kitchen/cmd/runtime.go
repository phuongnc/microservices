package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"infra/common/kafka"
	logger "infra/common/log"

	"kitchen-service/event"
	src "kitchen-service/src"

	"infra/common/db"
	"infra/order"

	"gorm.io/gorm"
)

type runtime struct {
	appConf          *AppConfig
	logger           *logger.Logger
	db               *gorm.DB
	kitchenHandler   src.KitchenHandler
	kitchenPublisher event.KitchenPublisher
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

	//setup kafka publisher
	kafkaConfig := &kafka.KafkaProducerConfiguration{
		BootstrapServers:  rt.appConf.KafkaConfig.BootstrapServers,
		Topic:             rt.appConf.KafkaConfig.KitchenEventTopic,
		TopicAutoCreation: true,
	}
	producer := kafka.NewKafkaMessageProducer(kafkaConfig)
	rt.kitchenPublisher = event.NewKitchenPublisher(producer)

	// init service handler
	orderRepository := order.NewOrderRepo()
	kitchenService := src.NewKitchenService(rt.logger, orderRepository, rt.kitchenPublisher)
	rt.kitchenHandler = src.NewKitchenHandler(rt.logger, kitchenService)

	//setup kafka consumer
	kafkaConsumerConfig := &kafka.KafkaConsumerConfiguration{
		BootstrapServers: rt.appConf.KafkaConfig.BootstrapServers,
		Topics:           []string{rt.appConf.KafkaConfig.PaymentEventTopic},
		GroupID:          rt.appConf.KafkaConfig.KitchenGroup,
		AutoCommit:       false,
	}
	ctx := context.WithValue(context.Background(), "db", rt.db)
	consumer := kafka.NewKafkaMessageConsumer(ctx, kafkaConsumerConfig)
	consumer.StartConsumer(kitchenService.PaymentConsumeEvent)

	return &rt
}

func (rt *runtime) Serve() {
	api := NewApi(rt.appConf, rt.db, rt.logger, rt.kitchenHandler)
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
		rt.kitchenPublisher.Destroy()
	}()
}
