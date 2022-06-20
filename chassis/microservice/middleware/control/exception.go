package control

import (
	"github.com/pigeatgarlic/goedf/chassis/microservice/endpoint"
	"github.com/pigeatgarlic/goedf/chassis/microservice/instruction"
	controlins "github.com/pigeatgarlic/goedf/chassis/microservice/instruction/control"
	"github.com/pigeatgarlic/goedf/chassis/microservice/middleware"
	"github.com/pigeatgarlic/goedf/models/event"
)

const (
	Name = "control"
)


func InitControlMiddleware(control *instruction.InstructionSet) *middleware.Middleware {
	ret := middleware.InitMiddlware("Control",map[string]string{
		"Author": "Pigeatgarlic",
	},map[string]*instruction.InstructionSet{
		"Control": control,
	})

	ret.Handler = func(next endpoint.EndpointFunction) endpoint.EndpointFunction {
		return func(event *event.Event) error {
			ret.InvokeInstruction("Control",controlins.Filter, event)
			return next(event)
		}
	}
	return ret
}



