package event

func (event *Event) finalAction(action *Action) *Action {
	next := event.FindActionByID(action.Next) 
	if next.Done {
		return action;
	} else {
		return event.finalAction(next);
	}
}

func (event *Event) CurrentAction() *Action {
	for i := 0; i < len(event.Actions); i++ {
		if event.Actions[i].Done {
			return event.finalAction(&event.Actions[i]);
		}
	}
	return nil;
}

func (event *Event) PreviousAction()  *Action {
	current := event.CurrentAction();
	if current == nil {
		return nil;
	} else {
		return event.FindActionByID(current.Prev);
	}
}

func (event *Event) FindActionByID(id int) *Action {
	for i := 0; i < len(event.Actions); i++ {
		if event.Actions[i].ID == id {
			return &event.Actions[i];
		}
	}
	return nil;
}