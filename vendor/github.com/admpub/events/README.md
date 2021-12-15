events [简体中文](https://github.com/admpub/events/blob/master/README_zh-CN.md)
===========

`events` is a small [Observer](https://en.wikipedia.org/wiki/Observer_pattern) implemetation for golang

Import
------

`events` available through github.com/admpub/events
interface:
```go
import "github.com/admpub/events"
```

Usage
-----

### Event

Creating standalone event object:
```go
event := events.New("eventName")
event.Meta["key"] = value
```

### Emitter

Package `emiter` implements `events.Emitter` interface
```go
import (
	"github.com/admpub/events"
	"github.com/admpub/events/emitter"
)
```

#### Create

Emitter could combined with other structs via common `events.Emitter` interface:
```go
type Object struct {
	events.Emitter
}

object := Object{emitter.New()}
```
> it's preferable usage example,
> it simplify test cases of base structs

Emitter could be created with specific dispatch strategy:
```go
import "github.com/admpub/events/dispatcher"
```

``` go
emitter.New(dispatcher.BroadcastFactory)
emitter.New(dispatcher.ParallelBroadcastFactory)
```

#### Emmit event

Emit concrete event object:

```go
em := emitter.New()
em.Fire(events.New("event"), events.ModeAsync)
```

Emit event with label & params:
```go
em.Fire("event", events.ModeSync)
// or with event params
em.Fire("event", events.ModeSync, meta.Map{"key": "value"})
// or with plain map
em.Fire("event", events.ModeSync, map[string]interface{}{"key": "value"})
````
> Be carefully with concurrent access to `event.Meta`

#### Subscribe for event

Emitter supports only `events.Listener` interface for subscription, but it can be extended by embedded types:

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
Package `ticker` adds support of periodic events on top of events.Emitter

```go
import (
	"github.com/admpub/events/emitter/ticker"
	"github.com/admpub/events/emitter"
	"time"
)
```

```go
tick := ticker.New(emitter.New())
tick.RegisterEvent("periodicEvent1", 5*time.Second)
// or
tick.RegisterEvent("periodicEvent2", time.NewTicker(5*time.Second))
// or directly with handlers
tick.RegisterEvent("periodicEvent3", 5*time.Second, Handler{})
```
