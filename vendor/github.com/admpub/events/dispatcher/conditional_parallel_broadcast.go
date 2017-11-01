package dispatcher

import (
	"sync"

	"github.com/admpub/events"
)

func ConditionalParallelBroadcastFactory() events.Dispatcher {
	return &ConditionalParallelBroadcastDispatcher{make([]events.Listener, 0)}
}

type ConditionalParallelBroadcastDispatcher struct {
	Subscribers []events.Listener
}

func (dispatcher *ConditionalParallelBroadcastDispatcher) AddSubscribers(subscribers ...events.Listener) {
	dispatcher.Subscribers = append(dispatcher.Subscribers, subscribers...)
}

func (dispatcher *ConditionalParallelBroadcastDispatcher) Dispatch(event events.Event) error {
	var err error
	if _, ok := event.Context["_sync"]; ok {
		delete(event.Context, "_sync")
		for _, subscriber := range dispatcher.Subscribers {
			if event.Aborted() {
				return err
			}
			err = subscriber.Handle(event)
			if err != nil {
				return err
			}
		}
	} else if _, ok := event.Context["_wait"]; ok {
		delete(event.Context, "_wait")
		wg := &sync.WaitGroup{}
		wg.Add(len(dispatcher.Subscribers))
		for _, subscriber := range dispatcher.Subscribers {
			go func(subscriber events.Listener) {
				subscriber.Handle(event)
				wg.Done()
			}(subscriber)
		}
		wg.Wait()
	} else {
		for _, subscriber := range dispatcher.Subscribers {
			go subscriber.Handle(event)
		}
	}
	return err
}
