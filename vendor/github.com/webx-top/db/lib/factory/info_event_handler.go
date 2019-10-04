package factory

type EventHandler func(model Model, editColumns ...string) error

func NewEventHandlers() *EventHandlers {
	return &EventHandlers{}
}

type EventHandlers struct {
	Async []EventHandler
	Sync  []EventHandler
}

func (e *EventHandlers) Exec(model Model, editColumns ...string) error {
	for _, handler := range e.Async {
		go handler(model, editColumns...)
	}
	for _, handler := range e.Sync {
		if err := handler(model, editColumns...); err != nil {
			return err
		}
	}
	return nil
}
