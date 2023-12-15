package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	config "bitbucket.org/inseinc/report-service/config"
	logger "bitbucket.org/inseinc/report-service/log"
	"gorm.io/gorm"
)

type runtime struct {
	appConf *config.AppConfig
	logger  *logger.Logger
	db      *gorm.DB
	service *Service
}

func NewRuntime() *runtime {
	rt := runtime{}
	var err error

	if rt.appConf, err = config.BuildConfiguration(); err != nil {
		fmt.Sprintf("can't load application configuration: %v", err)
	}

	// rt.db, err = db.NewSQL(rt.appConf.DatabaseConfig)
	// if err != nil {
	// 	fmt.Sprintf("cannot connect to db: %v", err)
	// }

	rt.logger = logger.New()

	service, err := NewService(rt.logger)
	if err != nil {
		rt.logger.Error("creating new service instance", err)
	}
	rt.service = service

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
