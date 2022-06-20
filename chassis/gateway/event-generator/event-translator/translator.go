package eventtranslator

import (
	"math/rand"

	registrator "github.com/pigeatgarlic/goedf/chassis/gateway/module/registrator"
	"github.com/pigeatgarlic/goedf/chassis/util/logger"
	"github.com/pigeatgarlic/goedf/models/event"
	"github.com/pigeatgarlic/goedf/models/request-response/request"
)

type EventTranslator struct {
	registrator registrator.ServiceRegistrator
}

func InitEventTranslator(log logger.Logger,
	registrator registrator.ServiceRegistrator) (*EventTranslator, error) {
	var err error
	var ret EventTranslator
	ret.registrator = registrator
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func (trans *EventTranslator) Translate(req *request.UserRequest) (*event.Event, error) {
	Actions, err := trans.registrator.EventLookup(req.Target)
	if err != nil {
		return nil, err
	}
	return &event.Event{
		ID:      rand.Int(),
		Headers: req.Headers,

		Actions: Actions,
	}, nil
}
