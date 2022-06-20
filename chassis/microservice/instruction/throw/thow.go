package throw

import (
	"fmt"

	"github.com/pigeatgarlic/goedf/chassis/microservice/instruction"
	"github.com/pigeatgarlic/goedf/chassis/util/logger"
	"github.com/pigeatgarlic/goedf/models/event"
)


func InitThowInstructionSet(log logger.Logger) instruction.Instruction {
	return func(prev *event.Result,
				current *event.Result,
				ID int,
				Headers map[string]string) error {
		if prev.Error != "" {
			log.Warning("Handled thrown error from previous action: " + prev.Error)
			current.Data = make(map[string]string)
			return fmt.Errorf(prev.Error)
		}
		return nil
	}
}
