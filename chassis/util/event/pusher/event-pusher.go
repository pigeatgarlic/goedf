package eventpusher

import (
	"github.com/pigeatgarlic/ideacrawler/microservice/models/event"
	"github.com/pigeatgarlic/ideacrawler/microservice/models/microservice"
)

type EventSpeaker interface {
	Push(event *event.Event) error
	ConfigTopic([]microservice.MicroserviceListenerConfig)
}
