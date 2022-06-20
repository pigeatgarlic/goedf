package redis

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pigeatgarlic/goedf/chassis/util/config"
	"github.com/pigeatgarlic/goedf/chassis/util/logger"
	"github.com/pigeatgarlic/goedf/models/event"
	"github.com/pigeatgarlic/goedf/models/microservice"
)

type redisCache struct {
	channel string

	redis   *redis.Client
	context context.Context

	logger logger.Logger
}

const (
	featureKey      = "Feature"
	microserviceKey = "Microservices"
)

func InitRedisCache(conf config.CacheInterface,
	logger logger.Logger) (*redisCache, error) {
	ctx := context.Background()
	redisClient := redis.NewClient(&redis.Options{
		Addr: conf.GetServer(),
	})

	res, err := redisClient.Ping(ctx).Result()
	if err != nil || res != "PONG" {
		return nil, err
	}

	return &redisCache{
		redis:   redisClient,
		context: ctx,
		channel: conf.GetChannel(),
		logger:  logger,
	}, nil
}

func (cache *redisCache) getAllEndpoint() ([]microservice.Endpoint, error) {
	var er error
	var ret []microservice.Endpoint
	cmd, er := cache.redis.Get(cache.context, microserviceKey).Result()
	if er != nil {
		return nil, er
	}

	var services []microservice.MicroService
	json.Unmarshal([]byte(cmd), &services)

	for _, svc := range services {
		ret = append(ret, svc.Endpoints...)
	}

	return ret, nil
}

func (cache *redisCache) Get(target string) (microservice.ActionSerires, error) {
	var ret microservice.ActionSerires
	cmd, er := cache.redis.Get(cache.context, featureKey).Result()
	if er != nil {
		return nil, er
	}

	var features []microservice.Feature
	json.Unmarshal([]byte(cmd), &features)

	endpoints, er := cache.getAllEndpoint()
	if er != nil {
		return nil, er
	}

	for _, feature := range features {
		if feature.Name == "target" {
			for _, endpoint := range endpoints {
				pick := false
				for _, work := range feature.EndpointIDs {
					if work == endpoint.ID {
						pick = true
					}
				}
				if pick {
					er = ret.Add(endpoint)
					if er != nil {
						return nil, er
					}
				}
			}
		}
	}
	return ret, nil
}

func (cache *redisCache) RegisterFeature(feature *microservice.Feature) error {
	var err error
	endpoint_snapshoot, err := cache.redis.Get(cache.context, microserviceKey).Result()
	feature_snapshoot, err := cache.redis.Get(cache.context, featureKey).Result()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if feature_snapshoot != "" {
				cache.redis.Set(cache.context, featureKey, feature_snapshoot, time.Second).Result()
			}
			if endpoint_snapshoot != "" {
				cache.redis.Set(cache.context, microserviceKey, endpoint_snapshoot, time.Second).Result()
			}
		}
	}()

	cmd, er := cache.redis.Get(cache.context, featureKey).Result()
	if er != nil {
		return er
	}
	var services []microservice.Feature
	json.Unmarshal([]byte(cmd), &services)
	services = append(services, *feature)
	_, er = cache.redis.Set(cache.context, featureKey, services, time.Second).Result()
	if er != nil {
		return er
	}

	// validate every possible feature
	cmd, er = cache.redis.Get(cache.context, featureKey).Result()
	if er != nil {
		return er
	}
	var features []microservice.Feature
	json.Unmarshal([]byte(cmd), &features)
	for _, feature := range features {
		_, err = cache.Get(feature.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cache *redisCache) RegisterMicroservice(service *microservice.MicroService) error {
	var err error
	endpoint_snapshoot, err := cache.redis.Get(cache.context, microserviceKey).Result()
	feature_snapshoot, err := cache.redis.Get(cache.context, featureKey).Result()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if feature_snapshoot != "" {
				cache.redis.Set(cache.context, featureKey, feature_snapshoot, time.Second).Result()
			}
			if endpoint_snapshoot != "" {
				cache.redis.Set(cache.context, microserviceKey, endpoint_snapshoot, time.Second).Result()
			}
		}
	}()

	cmd, er := cache.redis.Get(cache.context, microserviceKey).Result()
	if er != nil {
		return er
	}
	var services []microservice.MicroService
	json.Unmarshal([]byte(cmd), &services)
	services = append(services, *service)
	_, er = cache.redis.Set(cache.context, microserviceKey, services, time.Second).Result()
	if er != nil {
		return er
	}

	// validate every possible feature
	cmd, er = cache.redis.Get(cache.context, featureKey).Result()
	if er != nil {
		return er
	}
	var features []microservice.Feature
	json.Unmarshal([]byte(cmd), &features)
	for _, feature := range features {
		_, err = cache.Get(feature.Name)
		if err != nil {
			return err
		}
	}
	return nil

}

func (lut *redisCache) EventLookup(target string) ([]event.Action, error) {
	var err error
	var ret []event.Action
	result := make(microservice.ActionSerires)
	if err != nil {
		return nil, err
	}

	count := 0
	var prev_action *event.Action
	for i := 0; i < len(result); i++ {
		prev := count
		count = rand.Int()
		last := event.Action{
			ID:   count,
			Prev: prev,
			Next: prev_action.ID,

			Service:  result[i].MicroserviceID,
			Endpoint: result[i].EndpointID,
		}
		if i == 0 {
			last.Done = true
		}
		ret = append(ret, last)
		prev_action = &last
	}

	return ret, nil
}
