package main

import (
	"github.com/pigeatgarlic/goedf/chassis/gateway"
	"github.com/pigeatgarlic/goedf/chassis/util/config"
	eslogger "github.com/pigeatgarlic/goedf/chassis/util/logger/es"
)

func main() {

	sec, _ := config.GetSecurityConfig()
	gw, _ := config.GetGatewayConfig()
	cache, _ := config.GetCacheConfig()
	kafka, _ := config.GetMQConfig()
	pubsub, _ := config.GetPubsubConfig()
	logconf, _ := config.GetESlogConfig()

	log, _ := eslogger.InitLogger(logconf)
	gateway.NewGateway(sec, gw, cache, kafka, pubsub, log)
}
