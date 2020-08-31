package main

import (
	"log"
	"time"

	cors "github.com/itsjamie/gin-cors"
	"github.com/rfashwal/scs-actuator/internal"
	"github.com/rfashwal/scs-actuator/internal/config"
	httptransport "github.com/rfashwal/scs-actuator/internal/transport/http"
	"github.com/rfashwal/scs-actuator/internal/transport/mq"
	"github.com/rfashwal/scs-utilities/rabbit"
)

func main() {

	conf := config.Config().Manager
	mqManager, err := rabbit.NewRabbitMQManager(conf.RabbitURL())
	if err != nil {
		log.Fatalf("MQ server init: %s", err)
	}

	svc, err := internal.NewService(conf)
	if err != nil {
		log.Fatal("service init err", err)
	}

	go mq.TemperatureObserver(svc, mqManager, conf)

	router, err := httptransport.NewServer(svc)
	if err != nil {
		log.Fatalf("http server init: %s", err)
	}

	manager := config.EurekaManagerInit()
	manager.SendRegistrationOrFail()
	manager.ScheduleHeartBeat(conf.ServiceName(), 10)
	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE, OPTIONS",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	err = router.Run(conf.Address())

	if err != nil {
		log.Fatalf("router run: %s", err)
	}
}
