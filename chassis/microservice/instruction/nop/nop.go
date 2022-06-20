package nop

import (
	"github.com/pigeatgarlic/goedf/chassis/microservice/instruction"
	"github.com/pigeatgarlic/goedf/models/event"
)


func InitNopInstruction() instruction.Instruction{
	return func(prev *event.Result, current *event.Result, ID int, Headers map[string]string) error {
		return nil;
	}
}
