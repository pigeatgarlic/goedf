package exception

import (
	"github.com/pigeatgarlic/goedf/chassis/microservice/instruction/throw"
	"github.com/pigeatgarlic/goedf/chassis/microservice/middleware"
	"github.com/pigeatgarlic/goedf/chassis/util/logger"
)

func InitExceptionMiddleware(logger logger.Logger) *middleware.Middleware {
	return middleware.InitMiddlware("Exception catching", map[string]string{
		"Author": "Pigeatgarlic",
	}, throw.InitThowInstructionSet(logger),
	)

}
