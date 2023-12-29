package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"infra/common/kafka"
	logger "infra/common/log"

	"payment-service/event"
	src "payment-service/src"

	"infra/common/db"
	"infra/order"

	"gorm.io/gorm"
)

type runtime struct {
	appConf          *AppConfig
	logger           *logger.Logger
	db               *gorm.DB
	paymentHandler   src.PaymentHandler
	paymentPublisher event.PaymentPublisher
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
	paymentService := src.NewPaymentService(rt.logger, orderRepository, rt.paymentPublisher)
	rt.paymentHandler = src.NewPaymentHandler(rt.logger, paymentService)

	//setup kafka publisher
	kafkaConfig := &kafka.KafkaProducerConfiguration{
		BootstrapServers:  rt.appConf.KafkaConfig.BootstrapServers,
		Topic:             rt.appConf.KafkaConfig.PaymentEventTopic,
		TopicAutoCreation: true,
	}
	producer := kafka.NewKafkaMessageProducer(kafkaConfig)
	rt.paymentPublisher = event.NewPaymentPublisher(producer)
	//setup kafka consumer
	kafkaConsumerConfig := &kafka.KafkaConsumerConfiguration{
		BootstrapServers: rt.appConf.KafkaConfig.BootstrapServers,
		Topics:           []string{rt.appConf.KafkaConfig.OrderEventTopic},
		GroupID:          rt.appConf.KafkaConfig.PaymentGroup,
		AutoCommit:       false,
	}
	ctx := context.WithValue(context.Background(), "db", rt.db)
	consumer := kafka.NewKafkaMessageConsumer(ctx, kafkaConsumerConfig)
	consumer.StartConsumer(paymentService.PaymentConsumeEvent)

	return &rt
}

func (rt *runtime) Serve() {
	api := NewApi(rt.appConf, rt.db, rt.logger, rt.paymentHandler)
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
		rt.paymentPublisher.Destroy()
	}()
}
