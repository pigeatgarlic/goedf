package nop

import (
	"github.com/pigeatgarlic/goedf/chassis/microservice/endpoint"
	"github.com/pigeatgarlic/goedf/chassis/microservice/instruction"
	"github.com/pigeatgarlic/goedf/chassis/microservice/instruction/nop"
)

const (
	Nop = "NOP"
)

func NewNopEndpoint() *endpoint.Endpoint {
	ins := nop.InitNopInstruction()
	ret := endpoint.InitEndpoint( Nop, map[string]string{
		"Author" : "Pigeatgarlic",
	}, map[string]*instruction.InstructionSet{
		Nop: ins,
	})

	ret.AddStep(Nop, nop.NOP)
	return ret
}