package events

import "github.com/webx-top/echo/param"

type Map = param.Store

// DispatchStrategy defines strategy of delivery event to handlers
type DispatchStrategy func(Event, map[Listener]struct{}) error

// Listener defines event handler interface
type Listener interface {
	Handle(Event) error
}

type ID interface {
	ID() string
}

// Stream implements Listener interface on channel
type Stream chan Event

// Handle Listener
func (stream Stream) Handle(event Event) error {
	stream <- event
	return nil
}

type Streamer interface {
	Listener
	ID
	Chan() <-chan Event
}

func StreamWithID(ch chan Event, id string) Streamer {
	return &stream{
		ch: ch,
		id: id,
	}
}

// Stream implements Listener interface on channel
type stream struct {
	ch chan Event
	id string
}

// Handle Listener
func (s *stream) Handle(event Event) error {
	s.ch <- event
	return nil
}

func (s *stream) ID() string {
	return s.id
}

func (s *stream) Chan() <-chan Event {
	return s.ch
}

// Callback implements Listener interface on function
func Callback(function func(Event) error, id ...string) Listener {
	var _id string
	if len(id) > 0 {
		_id = id[0]
	}
	return callback{function: &function, id: _id}
}

type callback struct {
	function *func(Event) error
	id       string
}

// Handle Listener
func (c callback) Handle(event Event) error {
	return (*c.function)(event)
}

// ID Listener ID
func (c callback) ID() string {
	return c.id
}

func WithID(l Listener, id string) Listener {
	return &listenerWithID{Listener: l, id: id}
}

type listenerWithID struct {
	Listener
	id string
}

func (l listenerWithID) ID() string {
	return l.id
}
