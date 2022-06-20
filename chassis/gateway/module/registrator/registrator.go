package registrator

import (
	"github.com/pigeatgarlic/ideacrawler/microservice/models/event"
	"github.com/pigeatgarlic/ideacrawler/microservice/models/microservice"
)

type ServiceRegistrator interface {
	RegisterFeature(feature *microservice.Feature) error // perform by developer
	RegisterMicroservice(service *microservice.MicroService) error // perform by microservice itself

	EventLookup(target string) ([]event.Action, error) 
}