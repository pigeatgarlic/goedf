package redispubsub

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/util/config"
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/util/logger"
)

type redisPubSub struct {
	channel string

	redis   *redis.Client
	context context.Context

	pubsub *redis.PubSub
	queue  chan (*[]byte)

	logger logger.Logger
}

func InitRedisPubSub(conf *config.PubsubConfig, logger logger.Logger) (*redisPubSub, error) {
	ctx := context.Background()
	redisClient := redis.NewClient(&redis.Options{
		Addr: conf.Server,
	})

	res, err := redisClient.Ping(ctx).Result()

	if err != nil || res != "PONG" {
		return nil, err
	}

	pubsub := redisClient.PSubscribe(ctx, conf.Channel)
	queue := make(chan *[]byte)
	go func() {
		for {
			msg, err := pubsub.Receive(ctx)
			if err != nil {
				logger.Fatal(err.Error())
			}
			switch msg := msg.(type) {
			case *redis.Message:
				logger.Debug(fmt.Sprintf("Got redis message %s", msg.Payload))
				data := []byte(msg.Payload)
				queue <- &data
			default:
				logger.Debug("Got unknown redis message")
			}

		}
	}()

	return &redisPubSub{
		redis:   redisClient,
		context: ctx,
		pubsub:  pubsub,
		queue:   queue,
		channel: conf.Channel,
	}, nil
}

func (e *redisPubSub) Publish(message []byte) error {
	return e.redis.Publish(e.context, e.channel, message).Err()
}

func (e *redisPubSub) Subscribe() []byte {
	return *<-e.queue
}
