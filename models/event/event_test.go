package event

import "testing"

func TestEvent(t *testing.T) {
	event := Event{
		ID:      0,
		Headers: map[string]string{"test": "true"},
		Actions: []Action{Action{
				ID:0,

				Prev: 0,
				Next: 0,

				Service: 0,
				Endpoint: 0,

				Done: true,
				SignedAuthority: make([]string, 0),

				Result: Result{
					Data: map[string]string{},
					Error: "",
				},
			},
		},
	};
	data, err := event.ToProtobytes()
	if err != nil {
		t.Error(err);
		return;
	}
	result, err := FromProtobytes(data);
	if err != nil {
		t.Error(err);
		return;
	}
	if result.ID == event.ID {
		
	}
}