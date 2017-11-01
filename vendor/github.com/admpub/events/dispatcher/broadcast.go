package dispatcher

import (
	"github.com/admpub/events"
)

func BroadcastFactory() events.Dispatcher {
	return &BroadcastDispatcher{make([]events.Listener, 0)}
}

type BroadcastDispatcher struct {
	Subscribers []events.Listener
}

func (dispatcher *BroadcastDispatcher) AddSubscribers(subscribers ...events.Listener) {
	dispatcher.Subscribers = append(dispatcher.Subscribers, subscribers...)
}

func (dispatcher *BroadcastDispatcher) Dispatch(event events.Event) error {
	var err error
	for _, subscriber := range dispatcher.Subscribers {
		if event.Aborted() {
			return err
		}
		err = subscriber.Handle(event)
		if err != nil {
			return err
		}
	}
	return err
}
