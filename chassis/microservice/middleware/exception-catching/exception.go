package exception

import (
	"github.com/pigeatgarlic/goedf/chassis/microservice/endpoint"
	"github.com/pigeatgarlic/goedf/chassis/microservice/instruction"
	"github.com/pigeatgarlic/goedf/chassis/microservice/instruction/throw"
	"github.com/pigeatgarlic/goedf/chassis/microservice/middleware"
	"github.com/pigeatgarlic/goedf/chassis/util/logger"
	"github.com/pigeatgarlic/goedf/models/event"
)

func InitExceptionMiddleware(logger logger.Logger) *middleware.Middleware {
	ret := middleware.InitMiddlware("Exception catching", map[string]string{
		"Author": "Pigeatgarlic",
	}, map[string]*instruction.InstructionSet{
		throw.Throw: throw.InitThowInstructionSet(logger),
	})

	ret.Handler = func(next endpoint.EndpointFunction) endpoint.EndpointFunction {
		return func(event *event.Event) error {
			if err := ret.InvokeInstruction(throw.Throw, throw.Throw, event); err != nil {
				return err
			}
			return next(event)
		}
	}
	return ret
}
