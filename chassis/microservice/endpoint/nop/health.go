package nop

import (
	"github.com/pigeatgarlic/goedf/chassis/microservice/endpoint"
	"github.com/pigeatgarlic/goedf/chassis/microservice/instruction/nop"
)


func NewNopEndpoint() *endpoint.Endpoint {
	ins := nop.InitNopInstruction()
	return  endpoint.InitEndpoint( "NOP", map[string]string{ "Author" : "Pigeatgarlic", }, ins,)
}