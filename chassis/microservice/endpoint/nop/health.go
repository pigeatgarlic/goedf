package nop

import (
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/microservice/endpoint"
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/microservice/instruction/nop"
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/microservice/instruction"
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