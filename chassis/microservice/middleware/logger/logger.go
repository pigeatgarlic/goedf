package loggermiddleware

import (
	ins "github.com/pigeatgarlic/goedf/chassis/microservice/instruction/logger"
	"github.com/pigeatgarlic/goedf/chassis/microservice/middleware"
	"github.com/pigeatgarlic/goedf/chassis/util/logger"
)

func InitLoggerMiddleware(log logger.Logger) *middleware.Middleware {
	return middleware.InitMiddlware("Logger", map[string]string{
		"Author": "Pigeatgarlic",
	}, ins.InitLoggerInstruction(log),
	)
}
