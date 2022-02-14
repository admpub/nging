package events

import (
	"sort"
	"sync"
)

var Default = NewEmitter()

// EmitterOption defines option for Emitter
type EmitterOption struct {
	apply func(*Emitter)
}

// WithEventStategy sets delivery strategy for provided event
func WithEventStategy(event string, strategy DispatchStrategy) EmitterOption {
	return EmitterOption{func(emitter *Emitter) {
		if dispatcher, exists := emitter.dispatchers[event]; exists {
			dispatcher.strategy = strategy
			return
		}

		emitter.dispatchers[event] = NewDispatcher(strategy)
	}}
}

// WithDefaultStrategy sets default delivery strategy for event emitter
func WithDefaultStrategy(strategy DispatchStrategy) EmitterOption {
	return EmitterOption{func(emitter *Emitter) {
		emitter.strategy = strategy
	}}
}

// NewEmitter creates new event emitter
func NewEmitter(options ...EmitterOption) *Emitter {
	emitter := new(Emitter)
	emitter.strategy = Broadcast
	emitter.dispatchers = make(map[string]*Dispatcher)

	for _, option := range options {
		option.apply(emitter)
	}

	return emitter
}

// Emitter
type Emitter struct {
	guard       sync.Mutex
	strategy    DispatchStrategy
	dispatchers map[string]*Dispatcher
}

type Emitterer interface {
	On(string, ...Listener) Emitterer
	AddEventListener(handler Listener, events ...string)
	Off(string) Emitterer
	RemoveEventListener(handler Listener)
	Fire(interface{}) error
	FireByName(name string, options ...EventOption) error
	FireByNameWithMap(name string, data Map) error
	EventNames() []string
	HasEvent(string) bool
}

// On subscribes listeners to provided event and return emitter
// usefull for chain subscriptions
func (emitter *Emitter) On(event string, handlers ...Listener) Emitterer {
	emitter.AddEventListeners(event, handlers...)
	return emitter
}

// AddEventListeners subscribes listeners to provided event
func (emitter *Emitter) AddEventListeners(event string, handlers ...Listener) {
	emitter.guard.Lock()

	if _, exists := emitter.dispatchers[event]; !exists {
		emitter.dispatchers[event] = NewDispatcher(emitter.strategy)
	}
	emitter.dispatchers[event].AddSubscribers(handlers)

	emitter.guard.Unlock()
}

// AddEventListener subscribes listeners to provided events
func (emitter *Emitter) AddEventListener(handler Listener, events ...string) {
	emitter.guard.Lock()
	for _, event := range events {
		if _, exists := emitter.dispatchers[event]; !exists {
			emitter.dispatchers[event] = NewDispatcher(emitter.strategy)
		}

		emitter.dispatchers[event].AddSubscriber(handler)
	}
	emitter.guard.Unlock()
}

// Off unsubscribe all listeners from provided event
func (emitter *Emitter) Off(event string) Emitterer {
	emitter.RemoveEventListeners(event)
	return emitter
}

// RemoveEventListeners unsubscribe all listeners from provided event
func (emitter *Emitter) RemoveEventListeners(event string) {
	emitter.guard.Lock()
	delete(emitter.dispatchers, event)
	emitter.guard.Unlock()
}

// RemoveEventListener unsubscribe provided listener from all events
func (emitter *Emitter) RemoveEventListener(handler Listener) {
	emitter.guard.Lock()
	for _, dispatcher := range emitter.dispatchers {
		dispatcher.RemoveSubscriber(handler)
	}
	emitter.guard.Unlock()
}

// Fire start delivering event to listeners
func (emitter *Emitter) Fire(data interface{}) (err error) {
	event := New(data)
	if dispatcher, ok := emitter.dispatchers[event.Key]; ok {
		err = dispatcher.Dispatch(event)
	}
	return
}

func (emitter *Emitter) FireByName(name string, options ...EventOption) error {
	return emitter.Fire(New(name, options...))
}

func (emitter *Emitter) FireByNameWithMap(name string, data Map) error {
	return emitter.Fire(New(name, WithContext(data)))
}

// EventNames ...
func (emitter *Emitter) EventNames() []string {
	names := make([]string, len(emitter.dispatchers))
	var i int
	for name := range emitter.dispatchers {
		names[i] = name
		i++
	}
	sort.Strings(names)
	return names
}

// HasEvent ...
func (emitter *Emitter) HasEvent(event string) bool {
	_, ok := emitter.dispatchers[event]
	return ok
}
