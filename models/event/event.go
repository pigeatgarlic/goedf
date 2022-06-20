package event

const (
	GRPC = "GRPC Gateway"
	Crawler = "Event crawler"
	Microservice = "Event crawler"
)

type Event struct {
	ID		int 				`json:"ID"`
	Headers	map[string]string   `json:"Headers"`

	Actions []Action			`json:"Action"`
}

type Action struct {
	ID	      int		`json:"ID"`

	Next      int		`json:"Next"`
	Prev      int		`json:"Prev"`

	Service   int		`json:"Service"`
	Endpoint  int		`json:"Endpoint"` 
	
	Done				bool		`json:"Finish"`
	SignedAuthority     []string	`json:"Signed"`

	Result    Result	`json:"Result"`
}

type Result struct {
	Error	string				`json:"Error"`
	Data	map[string]string	`json:"Data"`
}


func (action *Action) MarkAsDone(err error,authorities ...string){
	action.Done = true;
	action.SignedAuthority = append(action.SignedAuthority, authorities...);
	
	// add error infor to upstream event
	action.Result.Error = err.Error();
}

