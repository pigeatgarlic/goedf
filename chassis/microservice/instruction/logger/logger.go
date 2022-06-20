package loggerinstrcution

import (
	"fmt"

	"github.com/pigeatgarlic/goedf/chassis/microservice/instruction"
	"github.com/pigeatgarlic/goedf/chassis/util/logger"
	"github.com/pigeatgarlic/goedf/models/event"
)

const (
	Infor = "Infor"
)

func InitLoggerInstruction(logger logger.Logger) *instruction.InstructionSet {
	ret := instruction.InitInstruction("Logging", map[string]string{
		"Author": "Pigeatgarlic",
	})

	ret.DescribeInstruction(Infor, func(prev *event.Result,
		current *event.Result,
		ID int,
		Headers map[string]string) error {
		logger.Infor(fmt.Sprintf("Got event %d in event name %s, username %s",
			ID,
			Headers["Name"],
			Headers["Username"]))
		return nil
	})

	return ret
}
