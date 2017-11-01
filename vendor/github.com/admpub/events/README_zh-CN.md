events
===========

`events` 是一个golang版本的观察者模式实现 [Observer](https://en.wikipedia.org/wiki/Observer_pattern)

Import
------

`events` 通过 github.com/admpub/events 导入:
```go
import "github.com/admpub/events"
```

用法
-----

### 事件

创建独立的事件对象:
```go
event := events.New("eventName")
event.Meta["key"] = value
```

### 发射器

包 `emiter` 实现了 `events.Emitter` 接口
```go
import (
	"github.com/admpub/events"
	"github.com/admpub/events/emitter"
)
```

#### 创建发射器

发射器可以通过 `events.Emitter` 接口嵌入到其它的结构体内:
```go
type Object struct {
	events.Emitter
}

object := Object{emitter.New()}
```
> 这是一个典型范例,
> 这里简化了结构体，实际情况下应该比这个复杂

可以用特定的调度策略来创建发射器:
```go
import "github.com/admpub/events/dispatcher"
```

``` go
emitter.New(dispatcher.BroadcastFactory)
emitter.New(dispatcher.ParallelBroadcastFactory)
```

#### 发射事件

发射具体事件对象:

```go
em := emitter.New()
em.Fire(events.New("event"), events.ModeAsync)
```

通过标签和参数发送事件(内部会根据标签名自动创建事件对象):
```go
em.Fire("event", events.ModeSync)
// or with event params
em.Fire("event", events.ModeSync, meta.Map{"key": "value"})
// or with plain map
em.Fire("event", events.ModeSync, map[string]interface{}{"key": "value"})
````
> Fire的第二个参数指定执行方式，目前支持以下三个值：
> - `events.ModeAsync`: 异步执行。忽略所有错误
> - `events.ModeSync`: 非异步执行。即依次执行，返回错误时会跳过后续执行
> - `events.ModeWait`: 异步执行，并等待完成。忽略所有错误

> 在并发时小心访问 `event.Meta`

#### 订阅事件

发射器仅支持用“events.Listener”接口来订阅，但可以通过嵌入式类型进行扩展:

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
包 `ticker` 在 `events.Emitter` 之上增加了对定期事件的支持

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
