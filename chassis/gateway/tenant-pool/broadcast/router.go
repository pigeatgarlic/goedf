package broadcast

import (
	redispubsub "github.com/pigeatgarlic/ideacrawler/microservice/chassis/gateway/tenant-pool/broadcast/redis"
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/util/config"
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/util/logger"
	"github.com/pigeatgarlic/ideacrawler/microservice/models/request-response/response"
)

type PubSubProvider interface {
	Publish(message []byte) error
	Subscribe() []byte
}

type Broadcast struct {
	provider PubSubProvider
}

func InitBroadcaster(conf *config.PubsubConfig, logger logger.Logger) (*Broadcast, error) {
	redis, err := redispubsub.InitRedisPubSub(conf, logger)
	if err != nil {
		return nil, err
	}

	return &Broadcast{
		provider: redis,
	}, nil
}

// TODO add retry
func (broadcast *Broadcast) Publish(resp *response.UserResponse) error {
	data, err := resp.ToProtobytes()
	if err != nil {
		return err
	}
	err = broadcast.provider.Publish(data)
	if err != nil {
		return err
	}
	return nil
}

func (broadcast *Broadcast) Subscribe() (*response.UserResponse, error) {
	data := broadcast.provider.Subscribe()
	return response.FromProtobytes(&data)
}
