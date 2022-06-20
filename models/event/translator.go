package event

import (
	proto "github.com/golang/protobuf/proto"
	protoc "github.com/pigeatgarlic/goedf/models/event/protobuf"
)

func (from Event) ToProtobytes() ([]byte, error) {
	event := protoc.Event{
		EventId: int64(from.ID),
		Headers: from.Headers,
	}

	for _, fromAction := range from.Actions {
		action := fromAction.protocTranslate()
		event.Actions = append(event.Actions, &action)
	}
	return proto.Marshal(&event)
}

func FromProtobytes(data []byte) (*Event, error) {
	var protocEvent protoc.Event
	err := proto.Unmarshal(data, &protocEvent)
	if err != nil {
		return nil, err
	}
	finalEvent := eventReverseTranslate(&protocEvent)
	return &finalEvent, nil
}

func (from Event) ProtocTranslate() protoc.Event {
	event := protoc.Event{
		EventId: int64(from.ID),
		Headers: from.Headers,
	}

	for _, fromAction := range from.Actions {
		action := fromAction.protocTranslate()
		event.Actions = append(event.Actions, &action)
	}
	return event
}

func (from Action) protocTranslate() protoc.Action {
	result := from.Result.protocTranslate()
	return protoc.Action{
		ActionId: int64(from.ID),

		NextAction:     int64(from.Next),
		PreviousAction: int64(from.Prev),

		ServiceId:  int64(from.Service),
		EndpointId: int64(from.Endpoint),

		Done:            from.Done,
		SignedAuthority: from.SignedAuthority,
		Result:          &result,
	}
}

func (from Result) protocTranslate() protoc.Result {
	return protoc.Result{
		Error: from.Error,
		Data:  from.Data,
	}
}

func eventReverseTranslate(from *protoc.Event) Event {
	event := Event{
		ID:      int(from.EventId),
		Headers: from.Headers,
	}

	for _, fromAction := range from.Actions {
		action := actionReverseTranslate(fromAction)
		event.Actions = append(event.Actions, action)
	}
	return event
}

func actionReverseTranslate(from *protoc.Action) Action {
	result := resultReverseTranslate(from.Result)
	return Action{
		ID: int(from.ActionId),

		Next: int(from.NextAction),
		Prev: int(from.PreviousAction),

		Service:  int(from.ServiceId),
		Endpoint: int(from.EndpointId),

		Done:            from.Done,
		SignedAuthority: from.SignedAuthority,
		Result:          result,
	}
}

func resultReverseTranslate(from *protoc.Result) Result {
	return Result{
		Error: from.Error,
		Data:  from.Data,
	}

}
