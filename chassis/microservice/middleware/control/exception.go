package control

import (
	"github.com/pigeatgarlic/goedf/chassis/microservice/event-handler"
	controlins "github.com/pigeatgarlic/goedf/chassis/microservice/instruction/control"
	"github.com/pigeatgarlic/goedf/chassis/microservice/middleware"
)



func InitControlMiddleware(handler *eventhandler.EventHandler) *middleware.Middleware {
	return middleware.InitMiddlware("Control",map[string]string{
		"Author": "Pigeatgarlic",
		},controlins.InitControlInstruction(handler),
	)
}



