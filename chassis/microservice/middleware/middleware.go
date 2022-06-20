package middleware

import (
	"github.com/pigeatgarlic/goedf/chassis/microservice/endpoint"
	"github.com/pigeatgarlic/goedf/chassis/microservice/instruction"
	"github.com/pigeatgarlic/goedf/models/event"
)

type Middleware struct {
	ID   int
	Name string
	Tags map[string]string

	instructionSets map[string]*instruction.InstructionSet

	Handler MiddlewareHandler
}

func InitMiddlware(name string,
	tags map[string]string,
	instruction map[string]*instruction.InstructionSet) *Middleware {
	return &Middleware{
		Name:            name,
		Tags:            tags,
		instructionSets: instruction,
	}
}

type MiddlewareHandler func(endpoint.EndpointFunction) endpoint.EndpointFunction

// ChainMap is a helper function for composing middlewares.
// Requests will traverse them in the order defined in dictionary.
func ChainMap(dict map[int]*Middleware) MiddlewareHandler {
	var outer MiddlewareHandler
	return func(next endpoint.EndpointFunction) endpoint.EndpointFunction {
		for i := 0; i < len(dict); i++ {
			next = dict[i].Handler(next)
		}
		return outer(next)
	}
}

func (middleware *Middleware) InvokeInstruction(instruction string,
	key string,
	event *event.Event) error {

	middleware.instructionSets[instruction].InvokeFunction(key, event)
	return nil
}
