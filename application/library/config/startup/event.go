package startup

func NewEvent() *Event {
	return &Event{
		Before: []EventFunc{},
		After:  []EventFunc{},
	}
}

type EventFunc func()

type Event struct {
	Before []EventFunc
	After  []EventFunc
}

func (e *Event) AddBefore(fn EventFunc) *Event {
	e.Before = append(e.Before, fn)
	return e
}

func (e *Event) AddAfter(fn EventFunc) *Event {
	e.After = append(e.After, fn)
	return e
}

func (e *Event) RunBefore() *Event {
	for _, fn := range e.Before {
		fn()
	}
	return e
}

func (e *Event) RunAfter() *Event {
	for _, fn := range e.After {
		fn()
	}
	return e
}
