package main

import (
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/microservice"
	// "github.com/pigeatgarlic/ideacrawler/microservice/chassis/microservice/middleware"
	// "github.com/pigeatgarlic/ideacrawler/microservice/chassis/microservice/endpoint"
	// "github.com/pigeatgarlic/ideacrawler/microservice/chassis/microservice/instruction"
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/util/config"
	eslogger "github.com/pigeatgarlic/ideacrawler/microservice/chassis/util/logger/es"
)

func main() {
	logconf,_ :=config.GetESlogConfig(); 

	log,err :=	eslogger.InitLogger(logconf);
	if err != nil {
		log.Fatal(err.Error())
	}
	auth, err := microservice.InitMicroService(log, &config.KafkaConfig{})
	if err != nil {
		log.Fatal(err.Error())
	}

	auth.Start()
}
