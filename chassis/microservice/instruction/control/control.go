package control

import (
	"encoding/json"

	eventhandler "github.com/pigeatgarlic/goedf/chassis/microservice/event-handler"
	"github.com/pigeatgarlic/goedf/chassis/microservice/instruction"
	"github.com/pigeatgarlic/goedf/models/event"
	"github.com/pigeatgarlic/goedf/models/microservice"
)

const (
	Filter = "Filter"
)

func InitControlInstruction(handler *eventhandler.EventHandler) *instruction.InstructionSet {
	ret := instruction.InitInstruction("Control", map[string]string{
		"Author": "Pigeatgarlic",
	},
	)

	ret.Tags["Author"] = "Pigeatgarlic"

	ret.DescribeInstruction(Filter, func(prev *event.Result, current *event.Result, ID int, Headers map[string]string) error {
		switch Headers["ControlEvent"] {
		case "UpdateGrid":
			var config map[int]string
			json.Unmarshal([]byte(prev.Data["Service"]), &config)

			handler.ConfigTopic(config)

			current.Data["Service"] = prev.Data["Service"]
		case "NewService":
			var svc microservice.MicroService
			json.Unmarshal([]byte(prev.Data["Service"]), &svc)

			current.Data["Service"] = prev.Data["Service"]
		}
		return nil
	})

	return ret
}
