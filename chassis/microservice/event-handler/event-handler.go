package eventhandler

import (
	eventpusher "github.com/pigeatgarlic/goedf/chassis/util/event/pusher"
	eventquerier "github.com/pigeatgarlic/goedf/chassis/util/event/querier"
	"github.com/pigeatgarlic/goedf/chassis/util/logger"
	"github.com/pigeatgarlic/goedf/models/event"
	"github.com/pigeatgarlic/goedf/models/microservice"
)

type EventHandler struct {
	querier eventquerier.EventListener
	pusher  eventpusher.EventSpeaker

	queue chan (*event.Event)
}

type ExecuteEndpoint func(event *event.Event) (event_after *event.Event)

func InitEventHandler(log logger.Logger,
	querier eventquerier.EventListener,
	pusher eventpusher.EventSpeaker) (*EventHandler, error) {
	var ret EventHandler

	ret.queue = make(chan *event.Event)
	ret.querier = querier
	ret.pusher = pusher

	return &ret, nil

}

func (handler *EventHandler) Start() {
	go func() {
		for {
			event := handler.WaitProcessedEvent()
			go handler.Push(event)
		}
	}()
}

func (handler *EventHandler) Push(event *event.Event) error {
	return handler.pusher.Push(event)
}

func (handler *EventHandler) ConfigTopic(config []microservice.MicroserviceListenerConfig) {
	handler.pusher.ConfigTopic(config)
}

func (processor *EventHandler) ProcessEvent(event *event.Event, process ExecuteEndpoint) {
	result := process(event)
	processor.queue <- result
}

func (processor *EventHandler) WaitProcessedEvent() *event.Event {
	return <-processor.queue
}

func (processor *EventHandler) WaitIncomingEvent() *event.Event {
	return processor.querier.WaitIncomingEvent()
}
