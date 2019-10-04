package factory

type EventHandler func(model Model) error

func NewEventHandlers() *EventHandlers {
	return &EventHandlers{}
}

type EventHandlers struct {
	Async []EventHandler
	Sync  []EventHandler
}

func (e *EventHandlers) Exec(model Model) error {
	for _, handler := range e.Async {
		go handler(model)
	}
	for _, handler := range e.Sync {
		if err := handler(model); err != nil {
			return err
		}
	}
	return nil
}
