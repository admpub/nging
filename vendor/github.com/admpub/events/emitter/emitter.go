package emitter

import (
	"sync"

	"github.com/admpub/events"
	"github.com/admpub/events/dispatcher"
	"github.com/webx-top/echo/param"
)

var (
	DefaultDispatcherFactory = dispatcher.BroadcastFactory
	DefaultAsyncEmitter      = New(dispatcher.ParallelBroadcastFactory)
	DefaultSyncEmitter       = New(dispatcher.BroadcastFactory)
	DefaultCondEmitter       = New(dispatcher.ConditionalParallelBroadcastFactory)
)

const (
	Async = iota
	Sync
	Cond
)

func NewWithType(typ int) *Emitter {
	switch typ {
	case Sync:
		return New(dispatcher.BroadcastFactory)
	case Cond:
		return New(dispatcher.ConditionalParallelBroadcastFactory)
	default:
		return New(dispatcher.ParallelBroadcastFactory)
	}
}

func New(factory ...events.DispatcherFactory) *Emitter {
	emitter := new(Emitter)
	emitter.Dispatchers = make(map[string]events.Dispatcher)
	if len(factory) > 0 {
		emitter.DispatcherFactory = factory[0]
	} else {
		emitter.DispatcherFactory = DefaultDispatcherFactory
	}
	return emitter
}

type Emitter struct {
	sync.Mutex
	DispatcherFactory events.DispatcherFactory
	Dispatchers       map[string]events.Dispatcher
}

func (emitter Emitter) On(event string, handlers ...events.Listener) events.Emitter {
	emitter.Lock()
	if _, exists := emitter.Dispatchers[event]; !exists {
		emitter.Dispatchers[event] = emitter.DispatcherFactory()
	}
	emitter.Dispatchers[event].AddSubscribers(handlers...)
	emitter.Unlock()
	return emitter
}

func (emitter Emitter) Off(event string) events.Emitter {
	emitter.Lock()
	delete(emitter.Dispatchers, event)
	emitter.Unlock()
	return emitter
}

func (emitter Emitter) HasEvent(event string) bool {
	emitter.Lock()
	_, exists := emitter.Dispatchers[event]
	emitter.Unlock()
	return exists
}

func (emitter Emitter) Fire(e interface{}, mode int, context ...param.Store) error {
	emitter.Lock()
	var (
		event events.Event
		err   error
	)

	switch e := e.(type) {
	case string:
		event = events.New(e)
	case events.Event:
		event = e
	}

	if len(context) > 0 {
		event.Context.DeepMerge(context[0])
	}

	if mode > -1 {
		switch mode {
		case events.ModeSync:
			event.Context["_sync"] = struct{}{}
		case events.ModeWait:
			event.Context["_wait"] = struct{}{}
		}
	}

	if dispatcher, ok := emitter.Dispatchers[event.Key]; ok {
		err = dispatcher.Dispatch(event)
	}
	emitter.Unlock()
	return err
}

func (emitter Emitter) Events() map[string]events.Dispatcher {
	return emitter.Dispatchers
}
