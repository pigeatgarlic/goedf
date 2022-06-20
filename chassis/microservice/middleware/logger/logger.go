package loggermiddleware

import (
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/microservice/endpoint"
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/microservice/instruction"
	loggerinstrcution "github.com/pigeatgarlic/ideacrawler/microservice/chassis/microservice/instruction/logger"
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/microservice/middleware"
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/util/logger"
	"github.com/pigeatgarlic/ideacrawler/microservice/models/event"
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
