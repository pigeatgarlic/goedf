package instruction

import (
	"fmt"

	"github.com/pigeatgarlic/ideacrawler/microservice/models/event"
)


type InstructionSet struct {
	ID   int
	Name string
	Tags  map[string]string

	instructions map[string]instruction
}

func InitInstruction(name string, tags map[string]string) *InstructionSet{
	return &InstructionSet{
		Name: name,
		Tags: tags,
		instructions: make(map[string]instruction),
	}
}

type instruction func(*event.Event) error

type Instruction func(prev *event.Result, current *event.Result, EventID int, Headers map[string]string) (error)

func (service *InstructionSet) DescribeInstruction(key string, handler Instruction) *InstructionSet {
	service.instructions[key] = func(processing_event *event.Event) error {
		err := handler(
			&processing_event.PreviousAction().Result,
			&processing_event.CurrentAction().Result,
			processing_event.ID,
			processing_event.Headers)

		current_action := processing_event.CurrentAction()
		current_action.SignedAuthority = append(current_action.SignedAuthority, fmt.Sprintf("Instruction set %s, Instruction %s", service.Name, key))

		return err
	}
	return service
}

func (service *InstructionSet) InvokeFunction(key string, event *event.Event) error {
	return service.instructions[key](event)
}
