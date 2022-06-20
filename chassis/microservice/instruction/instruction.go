package instruction

import (
	"github.com/pigeatgarlic/goedf/models/event"
)

type Instruction func(prev *event.Result, current *event.Result, EventID int, Headers map[string]string) error

