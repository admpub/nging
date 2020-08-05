package echo

import (
	"github.com/admpub/events"
	"github.com/admpub/events/emitter"
)

func On(name string, cb func(H) error, onErrorDie ...bool) {
	var onErrorStop bool
	if len(onErrorDie) > 0 {
		onErrorStop = onErrorDie[0]
	}
	emitter.DefaultCondEmitter.On(name, events.Callback(func(e events.Event) error {
		if onErrorStop {
			err := cb(e.Context)
			if err != nil {
				e.Abort()
			}
			return err
		}
		return cb(e.Context)
	}))
}

func Off(name string) {
	emitter.DefaultCondEmitter.Off(name)
}

const (
	EventAsyncMode = emitter.Async
	EventSyncMode  = emitter.Sync
	EventCondMode  = emitter.Cond
	EventNoneMode  = -1
)

func Fire(e interface{}, mode int, context ...H) error {
	return emitter.DefaultCondEmitter.Fire(e, mode, context...)
}

func Events() map[string]events.Dispatcher {
	return emitter.DefaultCondEmitter.Events()
}

func HasEvent(name string) bool {
	return emitter.DefaultCondEmitter.HasEvent(name)
}
