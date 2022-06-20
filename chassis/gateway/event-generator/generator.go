package eventgenerator

import (
	eventpusher "github.com/pigeatgarlic/goedf/chassis/util/event/pusher"

	eventtranslator "github.com/pigeatgarlic/goedf/chassis/gateway/event-generator/event-translator"
	registrator "github.com/pigeatgarlic/goedf/chassis/gateway/module/registrator"

	"github.com/pigeatgarlic/goedf/chassis/util/logger"
	"github.com/pigeatgarlic/goedf/models/request-response/request"
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
