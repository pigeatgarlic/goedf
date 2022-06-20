package eventquerier

import (
	"github.com/pigeatgarlic/goedf/models/event"
	"github.com/pigeatgarlic/goedf/models/microservice"
)

type EventListener interface {
	ConfigureTopic(topic microservice.MicroserviceListenerConfig)
	WaitIncomingEvent() *event.Event
}
