package control

import (
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/microservice/endpoint"
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/microservice/instruction"
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/microservice/middleware"
	controlins "github.com/pigeatgarlic/ideacrawler/microservice/chassis/microservice/instruction/control"
	"github.com/pigeatgarlic/ideacrawler/microservice/models/event"
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



