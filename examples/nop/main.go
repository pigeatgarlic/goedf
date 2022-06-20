package main

import (
	"github.com/pigeatgarlic/goedf/chassis/microservice"
	"github.com/pigeatgarlic/goedf/chassis/util/config"
	eslogger "github.com/pigeatgarlic/goedf/chassis/util/logger/es"
)

func main() {
	logconf, _ := config.GetESlogConfig()

	log, err := eslogger.InitLogger(logconf)
	if err != nil {
		log.Fatal(err.Error())
	}
	auth, err := microservice.InitMicroService(log, &config.KafkaConfig{})
	if err != nil {
		log.Fatal(err.Error())
	}


	auth.Start()
}
