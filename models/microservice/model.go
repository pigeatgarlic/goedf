package microservice



type MicroService struct {
	ID   int
	Name string
	Tags map[string]string


	Endpoints []Endpoint
}

type Endpoint struct {
	ID   int
	Name string
	Tags map[string]string

	InstructionSet []Instruction
	Order int

	MicroserviceID int
}

type Instruction struct {
	ID   int
	Name string
	Tags map[string]string

	Adapters []string
	Steps []string
}