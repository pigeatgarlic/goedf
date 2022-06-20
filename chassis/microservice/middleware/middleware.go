package middleware

import (
	"fmt"

	"github.com/pigeatgarlic/goedf/chassis/microservice/endpoint"
	"github.com/pigeatgarlic/goedf/chassis/microservice/instruction"
	"github.com/pigeatgarlic/goedf/models/event"
)

type Middleware struct {
	ID   int
	Name string
	Tags map[string]string

	handler MiddlewareHandler
}

func InitMiddlware(name string,
	tags map[string]string,
	instruction instruction.Instruction) *Middleware {
	middleware := Middleware{
		Name:            name,
		Tags:            tags,
	}
	middleware.handler = middleware.describeHandler(instruction);
	return &middleware
}

type MiddlewareHandler func(endpoint.EndpointFunction) endpoint.EndpointFunction

func (middle *Middleware) describeHandler(ins instruction.Instruction) MiddlewareHandler{
	return func(ef endpoint.EndpointFunction) endpoint.EndpointFunction {
		return func(event *event.Event) error {
			err := ins( &event.PreviousAction().Result,
						&event.CurrentAction().Result,
						event.ID,
						event.Headers)

			current_action := event.CurrentAction()
			current_action.SignedAuthority = append(current_action.SignedAuthority, 
				fmt.Sprintf("Middleware %s", middle.Name))
			if err != nil {
				return err
			}
			return ef(event);
		}
	}
}

// ChainMap is a helper function for composing middlewares.
// Requests will traverse them in the order defined in dictionary.
func ChainMap(dict map[int]*Middleware) MiddlewareHandler {
	var outer MiddlewareHandler
	return func(next endpoint.EndpointFunction) endpoint.EndpointFunction {
		for i := 0; i < len(dict); i++ {
			next = dict[i].handler(next)
		}
		return outer(next)
	}
}

