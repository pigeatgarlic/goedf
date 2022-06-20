package eventpusher

import (
	"github.com/pigeatgarlic/goedf/models/event"
	"github.com/pigeatgarlic/goedf/models/microservice"
)

type EventSpeaker interface {
	Push(event *event.Event) error
	ConfigTopic([]microservice.MicroserviceListenerConfig)
}
