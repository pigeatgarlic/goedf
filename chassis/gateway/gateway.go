package gateway

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	eventgenerator "github.com/pigeatgarlic/goedf/chassis/gateway/event-generator"
	eventwatcher "github.com/pigeatgarlic/goedf/chassis/gateway/event-watcher"
	helperendpoint "github.com/pigeatgarlic/goedf/chassis/gateway/helper"
	authenticator "github.com/pigeatgarlic/goedf/chassis/gateway/module/auth"
	registrator "github.com/pigeatgarlic/goedf/chassis/gateway/module/registrator"
	"github.com/pigeatgarlic/goedf/chassis/gateway/module/registrator/etcd"
	"github.com/pigeatgarlic/goedf/chassis/gateway/module/registrator/redis"
	tenantpool "github.com/pigeatgarlic/goedf/chassis/gateway/tenant-pool"
	tenantwatcher "github.com/pigeatgarlic/goedf/chassis/gateway/tenant-watcher"
	"github.com/pigeatgarlic/goedf/chassis/util/config"

	eventpusher "github.com/pigeatgarlic/goedf/chassis/util/event/pusher"
	eventquerier "github.com/pigeatgarlic/goedf/chassis/util/event/querier"

	"github.com/pigeatgarlic/goedf/chassis/util/event/provider/kafka"

	"github.com/pigeatgarlic/goedf/chassis/util/logger"
)

type Gateway struct {
	auth *authenticator.Authenticator

	watcher        *tenantwatcher.TenantWatcher
	helper         *helperendpoint.HelperWatcher
	pool           *tenantpool.TenantPool
	eventWatcher   *eventwatcher.EventWatcher
	eventGenerator *eventgenerator.EventGenerator
}

func NewGateway(
	sec *config.SecurityConfig,
	conf *config.GatewayConfig,
	cache_config config.CacheInterface,
	mq config.MQInterface,
	pubsub *config.PubsubConfig,
	log logger.Logger) error {
	var err error

	var gw Gateway
	var reg registrator.ServiceRegistrator
	var pusher eventpusher.EventSpeaker
	var querier eventquerier.EventListener

	switch cache_config.GetProvider() {
	case "redis":
		reg, err = redis.InitRedisCache(cache_config, log)
	case "etcd":
		reg, err = etcd.InitEtcdCache(cache_config, log)
	}

	switch mq.GetProvider() {
	case "kafka":
		querier, err = kafka.InitKafkaQuerier(log, mq)
		pusher, err = kafka.InitKafkaPusher(log, mq)
	case "redis":
	case "rabbitmq":
	case "nats":
	}

	// init authenticator module
	gw.auth = authenticator.InitAuthenticator(sec)

	// init sub-modules
	gw.eventWatcher, err = eventwatcher.InitEventWatcher(log, querier, reg)
	if err != nil {
		return err
	}
	gw.helper, err = helperendpoint.InitHelperEndpoint(log, sec, conf, reg)
	if err != nil {
		return err
	}
	gw.watcher, err = tenantwatcher.InitTenantWatcher(log, sec, conf)
	if err != nil {
		return err
	}
	gw.pool, err = tenantpool.InitTenantPool(pubsub, log)
	if err != nil {
		return err
	}
	gw.eventGenerator, err = eventgenerator.InitEventGenerator(log, pusher, reg)
	if err != nil {
		return err
	}

	// init deliver thread
	go func() {
		for {
			gw.pool.SendResponse(gw.eventWatcher.Wait())
		}
	}()

	go func() {
		for {
			gw.pool.NewTenant(gw.watcher.WaitTenant())
		}
	}()

	go func() {
		for {
			gw.pool.KillTenant(gw.watcher.WaitClose())
		}
	}()

	go func() {
		for {
			gw.eventGenerator.Push(gw.watcher.WaitRequest())
		}
	}()

	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- errors.New(fmt.Sprintf("%s", <-c))
	}()
	return <-errc
}
