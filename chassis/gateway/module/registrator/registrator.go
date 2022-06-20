package registrator

import (
	"github.com/pigeatgarlic/goedf/models/event"
	"github.com/pigeatgarlic/goedf/models/microservice"
)

type ServiceRegistrator interface {
	RegisterFeature(feature *microservice.Feature) error           // perform by developer
	RegisterMicroservice(service *microservice.MicroService) error // perform by microservice itself

	EventLookup(target string) ([]event.Action, error)
}
