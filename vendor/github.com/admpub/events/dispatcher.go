package events

var value = struct{}{}

type Dispatcherer interface {
	AddSubscriber(handler Listener)
	AddSubscribers([]Listener)
	RemoveSubscriber(handler Listener)
	Dispatch(Event) error
}

// NewDispatcher creates new dispatcher
func NewDispatcher(strategy DispatchStrategy) *Dispatcher {
	dispatcher := new(Dispatcher)
	dispatcher.strategy = strategy
	dispatcher.idm = make(map[string]Listener)
	dispatcher.subscribers = make(map[Listener]struct{})

	return dispatcher
}

// Dispatcher stores event listeners of concrete event
type Dispatcher struct {
	strategy    DispatchStrategy
	idm         map[string]Listener
	subscribers map[Listener]struct{}
}

// AddSubscriber adds one listener
func (dispatcher *Dispatcher) AddSubscriber(handler Listener) {
	dispatcher.add(handler)
}

func (dispatcher *Dispatcher) add(handler Listener) {
	if idi, ok := handler.(ID); ok {
		id := idi.ID()
		if len(id) > 0 {
			if h, y := dispatcher.idm[id]; y {
				dispatcher.RemoveSubscriber(h)
			}
			dispatcher.idm[id] = handler
		}
	}
	dispatcher.subscribers[handler] = value
}

// AddSubscribers adds slice of listeners
func (dispatcher *Dispatcher) AddSubscribers(handlers []Listener) {
	for _, handler := range handlers {
		dispatcher.add(handler)
	}
}

// RemoveSubscriber removes listener
func (dispatcher *Dispatcher) RemoveSubscriber(handler Listener) {
	delete(dispatcher.subscribers, handler)
}

// Dispatch deliver event to listeners using strategy
func (dispatcher *Dispatcher) Dispatch(event Event) error {
	return dispatcher.strategy(event, dispatcher.subscribers)
}
