package eventquerier

import (
	"github.com/pigeatgarlic/ideacrawler/microservice/models/event"
	"github.com/pigeatgarlic/ideacrawler/microservice/models/microservice"
)

type EventListener interface {
	ConfigureTopic(topic microservice.MicroserviceListenerConfig) 
	WaitIncomingEvent() *event.Event
}
