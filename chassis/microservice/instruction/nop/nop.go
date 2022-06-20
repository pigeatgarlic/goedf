package nop

import (
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/microservice/instruction"
	"github.com/pigeatgarlic/ideacrawler/microservice/models/event"
)

const (
	NOP = "Nop"
)

func InitNopInstruction() *instruction.InstructionSet {
	ret := instruction.InitInstruction(NOP,map[string]string{
		"Author": "Pigeatgarlic",
	})
	ret.DescribeInstruction(NOP, func(prev *event.Result, current *event.Result, ID int, Headers map[string]string) (error) {
		return nil;
	})

	return ret;
}
