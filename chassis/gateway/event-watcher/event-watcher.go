package eventwatcher

import (
	"encoding/json"
	"math/rand"
	"time"

	eventtranslator "github.com/pigeatgarlic/ideacrawler/microservice/chassis/gateway/event-watcher/event-translator"
	eventquerier "github.com/pigeatgarlic/ideacrawler/microservice/chassis/util/event/querier"

	registrator "github.com/pigeatgarlic/ideacrawler/microservice/chassis/gateway/module/registrator"
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/util/logger"
	"github.com/pigeatgarlic/ideacrawler/microservice/models/event"
	"github.com/pigeatgarlic/ideacrawler/microservice/models/microservice"
	"github.com/pigeatgarlic/ideacrawler/microservice/models/request-response/response"
)

type EventWatcher struct {
	querier     eventquerier.EventListener
	translator  *eventtranslator.EventTranslator
	registrator registrator.ServiceRegistrator

	middlewares []GatewayMiddleware
	channel     chan (*response.UserResponse)
}

type GatewayMiddleware func(*EventWatcher, *event.Event)

func InitEventWatcher(log logger.Logger,
	conf eventquerier.EventListener,
	registrator registrator.ServiceRegistrator) (*EventWatcher, error) {
	var err error
	var ret EventWatcher
	var querier eventquerier.EventListener
	trans := eventtranslator.InitEventTranslator()
	if err != nil {
		return nil, err
	}

	ret.middlewares = append(ret.middlewares, controlMiddleware)

	ret.registrator = registrator
	ret.channel = make(chan *response.UserResponse)
	ret.querier = querier
	ret.translator = trans
	go func() {
		for {
			event := ret.querier.WaitIncomingEvent()
			for i := 0; i < len(ret.middlewares); i++ {
				ret.middlewares[i](&ret, event)
			}

			resp, err := ret.translator.EventToResponse(event)
			if err != nil {
				break
			}
			ret.channel <- resp
		}
	}()

	return &EventWatcher{
		querier:    querier,
		translator: trans,
	}, nil
}

func controlMiddleware(watcher *EventWatcher, event *event.Event) {
	if event.Headers["Control"] == "MicroserviceRegister" {
		var err error
		var tag map[string]string
		err = json.Unmarshal([]byte(event.CurrentAction().Result.Data["Tags"]), &tag)
		if err != nil {
			return
		}

		var endpoints []microservice.Endpoint
		err = json.Unmarshal([]byte(event.CurrentAction().Result.Data["Endpoints"]), &endpoints)
		if err != nil {
			return
		}

		svc := microservice.MicroService{
			ID:        rand.Int(),
			Name:      event.CurrentAction().Result.Data["Name"],
			Tags:      tag,
			Endpoints: endpoints,
		}

		svc.Tags["RegisterAt"] = time.Now().Format(time.RFC3339)
		err = watcher.registrator.RegisterMicroservice(&svc)
		if err != nil {
			return
		}
	}
}

func (watcher *EventWatcher) Wait() *response.UserResponse {
	return <-watcher.channel
}
