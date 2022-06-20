package etcd

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/util/config"
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/util/logger"
	"github.com/pigeatgarlic/ideacrawler/microservice/models/event"
	"github.com/pigeatgarlic/ideacrawler/microservice/models/microservice"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type etcdCache struct {
	client *clientv3.Client
}

const (
	featureKey      = "Feature"
	microserviceKey = "Microservices"
)

func InitEtcdCache(conf config.CacheInterface,
	logger logger.Logger) (*etcdCache, error) {

	etcdConfig := clientv3.Config{
		Endpoints:   []string{conf.GetServer()},
		DialTimeout: 2 * time.Second,
		Username:    conf.GetUser(),
		Password:    conf.GetPassword(),
	}

	etcdClient, err := clientv3.New(etcdConfig)

	if err != nil {
		return nil, err
	}

	return &etcdCache{
		client: etcdClient,
	}, nil
}

func (cache *etcdCache) getAllEndpoint() ([]microservice.Endpoint, error) {
	var ret []microservice.Endpoint
	return ret, nil
}

func (cache *etcdCache) Get(target string) (microservice.ActionSerires, error) {
	var ret microservice.ActionSerires
	return ret, nil
}

func (cache *etcdCache) SetKeyValue(key string, value string) (err error) {
	ctx := context.Background()
	_, err = cache.client.Put(ctx, key, value)
	return
}

func (cache *etcdCache) GetKeyValue(key string) (value string, err error) {
	ctx := context.Background()
	resp, err := cache.client.Get(ctx, key)
	if err != nil {
		return "", err
	}

	for _, i := range resp.Kvs {
		value = string(i.Value)
	}
	return
}

func (cache *etcdCache) RegisterFeature(feature *microservice.Feature) error {
	var err error

	endpoint_snapshoot, err := cache.GetKeyValue(microserviceKey)
	feature_snapshoot, err := cache.GetKeyValue(featureKey)

	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if feature_snapshoot != "" {
				cache.SetKeyValue(featureKey, feature_snapshoot)
			}
			if endpoint_snapshoot != "" {
				cache.SetKeyValue(microserviceKey, endpoint_snapshoot)
			}
		}
	}()

	cmd, er := cache.GetKeyValue(featureKey)
	if er != nil {
		return er
	}
	var services []microservice.Feature
	json.Unmarshal([]byte(cmd), &services)
	services = append(services, *feature)
	data, err := json.Marshal(services)
	er = cache.SetKeyValue(featureKey, string(data))
	if er != nil {
		return er
	}

	// validate every possible feature
	cmd, er = cache.GetKeyValue(featureKey)
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

func (cache *etcdCache) RegisterMicroservice(service *microservice.MicroService) error {
	var err error

	endpoint_snapshoot, err := cache.GetKeyValue(microserviceKey)
	feature_snapshoot, err := cache.GetKeyValue(featureKey)

	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if feature_snapshoot != "" {
				cache.SetKeyValue(featureKey, feature_snapshoot)
			}
			if endpoint_snapshoot != "" {
				cache.SetKeyValue(microserviceKey, endpoint_snapshoot)
			}
		}
	}()

	cmd, er := cache.GetKeyValue(microserviceKey)
	if er != nil {
		return er
	}
	var services []microservice.MicroService
	json.Unmarshal([]byte(cmd), &services)
	services = append(services, *service)
	data, err := json.Marshal(services)
	er = cache.SetKeyValue(featureKey, string(data))
	if er != nil {
		return er
	}

	// validate every possible feature
	cmd, er = cache.GetKeyValue(featureKey)
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

func (lut *etcdCache) EventLookup(target string) ([]event.Action, error) {
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
