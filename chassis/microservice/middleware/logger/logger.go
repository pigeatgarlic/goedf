package loggermiddleware

import (
	"github.com/pigeatgarlic/goedf/chassis/microservice/endpoint"
	"github.com/pigeatgarlic/goedf/chassis/microservice/instruction"
	loggerinstrcution "github.com/pigeatgarlic/goedf/chassis/microservice/instruction/logger"
	"github.com/pigeatgarlic/goedf/chassis/microservice/middleware"
	"github.com/pigeatgarlic/goedf/chassis/util/logger"
	"github.com/pigeatgarlic/goedf/models/event"
)

func InitLoggerMiddleware(log logger.Logger) *middleware.Middleware {
	ret := middleware.InitMiddlware("Logger", map[string]string{
		"Author": "Pigeatgarlic",
	}, map[string]*instruction.InstructionSet{
		"Logger": loggerinstrcution.InitLoggerInstruction(log),
	})

	ret.Handler = func(next endpoint.EndpointFunction) endpoint.EndpointFunction {
		return func(event *event.Event) error {
			ret.InvokeInstruction("Logger", loggerinstrcution.Infor, event)
			return next(event)
		}
	}
	return ret
}
