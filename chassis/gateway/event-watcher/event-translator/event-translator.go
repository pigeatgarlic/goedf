package eventtranslator

import (
	"strconv"
	"github.com/pigeatgarlic/ideacrawler/microservice/models/event"
	"github.com/pigeatgarlic/ideacrawler/microservice/models/request-response/response"
)

type EventTranslator struct {
}


func InitEventTranslator() *EventTranslator {
	return &EventTranslator{}
}

func (translator *EventTranslator) EventToResponse(event *event.Event) (*response.UserResponse, error) {
	SessionID, ss_err := strconv.ParseInt(event.Headers["SessionID"],10,64);
	RequestID, rq_err := strconv.ParseInt(event.Headers["RequestID"],10,64);

	if ss_err != nil {
		return nil,ss_err
	}
	if rq_err != nil {
		return nil,rq_err
	}

	return &response.UserResponse{
		ID: int(RequestID),
		SessionID: uint64(SessionID),	

		Error: event.CurrentAction().Result.Error,
		Data: event.CurrentAction().Result.Data,
	},nil;
}

