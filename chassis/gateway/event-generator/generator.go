package eventgenerator

import (
	eventpusher "github.com/pigeatgarlic/ideacrawler/microservice/chassis/util/event/pusher"

	eventtranslator "github.com/pigeatgarlic/ideacrawler/microservice/chassis/gateway/event-generator/event-translator"
	registrator "github.com/pigeatgarlic/ideacrawler/microservice/chassis/gateway/module/registrator"

	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/util/logger"
	"github.com/pigeatgarlic/ideacrawler/microservice/models/request-response/request"
)

type EventGenerator struct {
	trans  *eventtranslator.EventTranslator
	pusher eventpusher.EventSpeaker

	log logger.Logger
}

func InitEventGenerator(log logger.Logger,
	pusher eventpusher.EventSpeaker,
	lut registrator.ServiceRegistrator) (*EventGenerator, error) {
	var ret EventGenerator

	trans, err := eventtranslator.InitEventTranslator(log, lut)
	if err != nil {
		return nil, err
	}

	ret.pusher = pusher
	ret.trans = trans
	ret.log = log

	return &ret, nil
}

func (gen *EventGenerator) Push(req *request.UserRequest) error {
	event, err := gen.trans.Translate(req)
	if err != nil {
		return err
	}
	return gen.pusher.Push(event)
}
