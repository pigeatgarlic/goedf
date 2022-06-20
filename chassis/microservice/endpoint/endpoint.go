package endpoint

import (
	"fmt"

	"github.com/pigeatgarlic/goedf/chassis/microservice/instruction"
	"github.com/pigeatgarlic/goedf/models/event"
)

type Endpoint struct {
	ID   int
	Name string
	Tags map[string]string

	finalHandler 	EndpointFunction
}


type EndpointFunction func(event *event.Event) error

func InitEndpoint(name string,
				  tag map[string]string,
				  handler instruction.Instruction) *Endpoint {
	endpoint := Endpoint{
		Name:            name,
		Tags:            tag,
	}
	endpoint.finalHandler = endpoint.describeEndpointFunction(handler);
	return &endpoint;
}


func (endpoint *Endpoint) describeEndpointFunction(ins instruction.Instruction) EndpointFunction{
	return func(event *event.Event) error {
		err := ins( &event.PreviousAction().Result,
					&event.CurrentAction().Result,
					 event.ID,
					 event.Headers)

		current_action := event.CurrentAction()
		current_action.SignedAuthority = append(current_action.SignedAuthority, fmt.Sprintf("Endpoint %s", endpoint.Name))

		return err
	}
}


func (endpoint *Endpoint) GetEndpointHandler() EndpointFunction {
	return endpoint.finalHandler
}
