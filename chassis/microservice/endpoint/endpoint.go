package endpoint

import (
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/microservice/instruction"
	"github.com/pigeatgarlic/ideacrawler/microservice/models/event"
)

type Endpoint struct {
	ID   int
	Name string
	Tags map[string]string

	handler EndpointFunction
	instructionSets map[string]*instruction.InstructionSet

	steps   map[int]step
}

type step struct {
	InstructionSetName string
	InstructionName    string
}

type EndpointFunction func(event *event.Event) error

func InitEndpoint(name string, 
				  tag map[string]string,
				  instructions map[string]*instruction.InstructionSet) *Endpoint {
	return &Endpoint{
		Name: name,
		Tags: tag,
		instructionSets: instructions,
		steps: make(map[int]step),
	}
}

func (endpoint *Endpoint) describeEndpointFunction(steps map[int]step) {
	endpoint.handler = func(event *event.Event) error {
		for i := 0; i < len(steps); i++ {
			err := endpoint.instructionSets[steps[i].InstructionSetName].InvokeFunction(steps[i].InstructionName, event)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func (endpoint *Endpoint) AddStep(instructionName string, instructionSetName string) *Endpoint {
	new_step := step{
		InstructionSetName: instructionSetName,
		InstructionName:    instructionName,
	}
	endpoint.steps[len(endpoint.steps)] = new_step
	endpoint.describeEndpointFunction(endpoint.steps)
	return endpoint
}

func (endpoint *Endpoint) GetEndpointHandler() EndpointFunction {
	return endpoint.handler
}
