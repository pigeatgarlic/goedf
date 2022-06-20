package microservice

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/microservice/endpoint"
	nopendpoint "github.com/pigeatgarlic/ideacrawler/microservice/chassis/microservice/endpoint/nop"
	eventhandler "github.com/pigeatgarlic/ideacrawler/microservice/chassis/microservice/event-handler"
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/microservice/instruction"
	controlinstruction "github.com/pigeatgarlic/ideacrawler/microservice/chassis/microservice/instruction/control"
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/microservice/middleware"
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/microservice/middleware/control"
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/microservice/middleware/exception-catching"
	loggermiddleware "github.com/pigeatgarlic/ideacrawler/microservice/chassis/microservice/middleware/logger"
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/util/config"
	eventpusher "github.com/pigeatgarlic/ideacrawler/microservice/chassis/util/event/pusher"
	eventquerier "github.com/pigeatgarlic/ideacrawler/microservice/chassis/util/event/querier"
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/util/event/provider/kafka"
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/util/logger"
	"github.com/pigeatgarlic/ideacrawler/microservice/models/event"
)

type MicroService struct {
	ID   int               `json:"ID"`
	Name string            `json:"Name"`
	Tags map[string]string `json:"Tags"`

	endpoints    []*endpoint.Endpoint
	instructions []*instruction.InstructionSet
	middlewares  map[int]*middleware.Middleware

	handler *eventhandler.EventHandler
	final   middleware.MiddlewareHandler

	logger logger.Logger
}

func InitMicroService(log logger.Logger,
	mq config.MQInterface) (*MicroService, error) {
	var err error
	var ret MicroService
	ret.middlewares = make(map[int]*middleware.Middleware)

	nopEndpoint := nopendpoint.NewNopEndpoint()
	controlInstructionSet := controlinstruction.InitControlInstruction(ret.handler)

	controlMiddleware := control.InitControlMiddleware(controlInstructionSet)
	loggerMiddleware := loggermiddleware.InitLoggerMiddleware(log)
	exceptionMiddleware := exception.InitExceptionMiddleware(log)

	ret.AddInstructionSet(controlInstructionSet)

	ret.AddMiddleware(controlMiddleware)
	ret.AddMiddleware(loggerMiddleware)
	ret.AddMiddleware(exceptionMiddleware)

	ret.AddEndpoint(nopEndpoint)

	var querier eventquerier.EventListener
	var pusher eventpusher.EventSpeaker
	switch mq.GetProvider() {
	case "kafka":
		querier, err = kafka.InitKafkaQuerier(log, mq)
		pusher, err = kafka.InitKafkaPusher(log, mq)
	case "redis":
	case "rabbitmq":
	case "nats":
	}
	ret.handler, err = eventhandler.InitEventHandler(log, querier, pusher)
	return &ret, err
}

func (service *MicroService) Start() error {
	// chain middleware by order specified in map
	service.final = middleware.ChainMap(service.middlewares)
	service.handler.Start()

	go func() {
		for {
			event := service.handler.WaitIncomingEvent()
			go service.findAndExecute(event)
		}
	}()

	// Interrupt handler.
	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()
	return <-errc
}

func (service *MicroService) findAndExecute(event *event.Event) {
	var err error
	var end *endpoint.Endpoint
	var exec eventhandler.ExecuteEndpoint

	// find matched endpoint
	end, err = service.findSuitableEndpoint(event)
	if err != nil {
		service.logger.Error(err.Error())
	}
	exec, err = service.getExecuteEndpoints(end)
	if err != nil {
		service.logger.Error(err.Error())
	}

	// run endpoint in go rountine
	go service.handler.ProcessEvent(event, exec)
}

func (service *MicroService) findSuitableEndpoint(event *event.Event) (*endpoint.Endpoint, error) {
	currentEndpoint := event.CurrentAction().Endpoint
	for i := 0; i < len(service.endpoints); i++ {
		if currentEndpoint == service.endpoints[i].ID {
			return service.endpoints[i], nil
		}
	}
	return nil, errors.New("cannot find suitable endpoint")
}

// return executeEndpoint with chained middleware
func (service *MicroService) getExecuteEndpoints(endpoint *endpoint.Endpoint) (eventhandler.ExecuteEndpoint, error) {
	return func(event *event.Event) *event.Event {
		err := service.final(endpoint.GetEndpointHandler())(event)
		event.CurrentAction().MarkAsDone(err, "Microservice "+service.Name)
		return event
	}, nil
}

func (service *MicroService) AddMiddleware(mdware *middleware.Middleware) *MicroService {
	for i := 0; i < len(service.middlewares); i++ {
		if service.middlewares[i] == nil {
			service.middlewares[i] = mdware
		}
	}
	return service
}

func (service *MicroService) AddEndpoint(endpoint *endpoint.Endpoint) *MicroService {
	service.endpoints = append(service.endpoints, endpoint)
	return service
}

func (service *MicroService) AddInstructionSet(instruction *instruction.InstructionSet) *MicroService {
	service.instructions = append(service.instructions, instruction)
	return service
}
