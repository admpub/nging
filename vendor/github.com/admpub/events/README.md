events
===========

Package `events` is implementation of [Observer](https://en.wikipedia.org/wiki/Observer_pattern)

[![GoDoc](https://godoc.org/github.com/admpub/events?status.svg)](https://godoc.org/github.com/admpub/events)

Import
------

`go.events` available through [gopkg.in](http://labix.org/gopkg.in) interface:

```go
import "github.com/admpub/events"
```

or directly from github:

```go
import "github.com/adone/go.events.v2"
```

Usage
-----

### Event

Create event:

```go
event := events.New("eventName")
event.Context["key"] = value
```

Or create with predefined context:

```go
data := events.Map{"foo": "bar"} // or map[string]interface{}{"foo": "bar"}
event := events.New("eventName", events.WithContext(data))
```

### Emitter

#### Create

`Emitter` can be embed into other struct:

```go
type Object struct {
	*events.Emitter
}

object := Object{events.NewEmitter()}
```

> it is a preferable use of `Emitter`

`Emitter` supports different delivery strategies:

```go
events.NewEmitter(events.WithDefaultStrategy(Broadcast))

events.NewEmitter(events.WithEventStrategy("event.for.parallel.processing", ParallelBroadcast))
```

You can define custom strategies by implementing `DispatchStrategy` function:

```go
customStrategy := func(event events.Event, listeners map[events.Listener]struct{}) error {
	// ...
	return nil
}
```

#### Fire event

```go
em := events.NewEmitter()

em.Fire(events.New("event"))
```

fire event with parameters

```go
em.Fire("event")
//or
em.Fire(events.New("event"))
````

#### Subscription on event

* channels

```go
channel := make(chan events.Event)

object.On("event", events.Stream(channel))
```

* handlers

```go
type Handler struct {}

func (Handler) Handle (events.Event) error {
	// handle events
	return nil
}

object.On("event", Handler{})
// or
object.On("event", Handler{}, Handler{}).On("anotherEvent", Handler{})
```

* functions

```go
object.On("event", events.Callback(func(event events.Event) error {
	// handle event
	return nil
}))
```

### Ticker

`PeriodicEmitter` adds support for periodic events on `Emitter`

```go
import (
	"github.com/admpub/events"
	"time"
)
```

```go
tick := events.NewTicker(events.NewEmitter())
tick.RegisterEvent("periodic.event.1", 5*time.Second)
// or
tick.RegisterEvent("periodic.event.2", time.NewTicker(5*time.Second))
// or
tick.RegisterEvent("periodic.event.3", 5*time.Second, Handler{})
```
