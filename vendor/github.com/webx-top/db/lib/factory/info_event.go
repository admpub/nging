package factory

const (
	EventCreating = "creating" // 创建前
	EventCreated  = "created"  // 创建后
	EventUpdating = "updating" // 更新前
	EventUpdated  = "updated"  // 更新后
	EventDeleting = "deleting" // 删除前
	EventDeleted  = "deleted"  // 删除后

	EventReading = "reading" // 读取前
	EventReaded  = "readed"  // 读取后
)

type Ranger interface {
	Range(func(model Model) error) error
}

func NewEvent() *Event {
	return &Event{}
}

var (
	// AllAfterWriteEvents 所有写事件
	AllAfterWriteEvents  = []string{EventCreated, EventUpdated, EventDeleted}
	AllBeforeWriteEvents = []string{EventCreating, EventUpdating, EventDeleting}

	// AllAfterReadEvents 所有读事件
	AllAfterReadEvents  = []string{EventReaded}
	AllBeforeReadEvents = []string{EventReading}
)

type Event struct {
	// Reading 读取之前
	Reading *EventReadHandlers
	// Readed 读取之后
	Readed *EventReadHandlers

	// Creating 创建之前
	Creating *EventHandlers
	// Created 创建之后
	Created *EventHandlers

	// Updating 更新之前
	Updating *EventHandlers
	// Updated 更新之后
	Updated *EventHandlers

	// Deleting 删除之前
	Deleting *EventHandlers
	// Deleted 删除之后
	Deleted *EventHandlers
}

func (e *Event) Exists(event string) bool {
	switch event {
	case EventReading:
		return e.Reading != nil
	case EventReaded:
		return e.Readed != nil
	case EventCreating:
		return e.Creating != nil
	case EventCreated:
		return e.Created != nil
	case EventUpdating:
		return e.Updating != nil
	case EventUpdated:
		return e.Updated != nil
	case EventDeleting:
		return e.Deleting != nil
	case EventDeleted:
		return e.Deleted != nil
	}
	return false
}

func (e *Event) CallRead(event string, model Model, param *Param) error {
	if !e.Exists(event) {
		return nil
	}
	switch event {
	case EventReading:
		return e.Reading.Exec(model, param)
	case EventReaded:
		return e.Readed.Exec(model, param)
	}
	return nil
}

func (e *Event) Call(event string, model Model, editColumns ...string) error {
	if !e.Exists(event) {
		return nil
	}
	switch event {
	case EventCreating:
		return e.Creating.Exec(model, editColumns...)
	case EventCreated:
		return e.Created.Exec(model, editColumns...)
	case EventUpdating:
		return e.Updating.Exec(model, editColumns...)
	case EventUpdated:
		return e.Updated.Exec(model, editColumns...)
	case EventDeleting:
		return e.Deleting.Exec(model, editColumns...)
	case EventDeleted:
		return e.Deleted.Exec(model, editColumns...)
	}
	return nil
}

func (e *Event) OnRead(event string, h EventReadHandler, async ...bool) *Event {
	switch event {
	case EventReading:
		return e.AddReading(h, async...)
	case EventReaded:
		return e.AddReaded(h, async...)
	default:
		panic(`Unsupported event: ` + event)
	}
}

func (e *Event) On(event string, h EventHandler, async ...bool) *Event {
	switch event {
	case EventCreating:
		return e.AddCreating(h, async...)
	case EventCreated:
		return e.AddCreated(h, async...)
	case EventUpdating:
		return e.AddUpdating(h, async...)
	case EventUpdated:
		return e.AddUpdated(h, async...)
	case EventDeleting:
		return e.AddDeleting(h, async...)
	case EventDeleted:
		return e.AddDeleted(h, async...)
	default:
		panic(`Unsupported event: ` + event)
	}
}

func (e *Event) AddReading(h EventReadHandler, async ...bool) *Event {
	if e.Reading == nil {
		e.Reading = NewEventReadHandlers()
	}
	if len(async) > 0 && async[0] {
		e.Reading.Async = append(e.Reading.Async, h)
	} else {
		e.Reading.Sync = append(e.Reading.Sync, h)
	}
	return e
}

func (e *Event) AddReaded(h EventReadHandler, async ...bool) *Event {
	if e.Readed == nil {
		e.Readed = NewEventReadHandlers()
	}
	if len(async) > 0 && async[0] {
		e.Readed.Async = append(e.Readed.Async, h)
	} else {
		e.Readed.Sync = append(e.Readed.Sync, h)
	}
	return e
}

func (e *Event) AddCreating(h EventHandler, async ...bool) *Event {
	if e.Creating == nil {
		e.Creating = NewEventHandlers()
	}
	if len(async) > 0 && async[0] {
		e.Creating.Async = append(e.Creating.Async, h)
	} else {
		e.Creating.Sync = append(e.Creating.Sync, h)
	}
	return e
}

func (e *Event) AddCreated(h EventHandler, async ...bool) *Event {
	if e.Created == nil {
		e.Created = NewEventHandlers()
	}
	if len(async) > 0 && async[0] {
		e.Created.Async = append(e.Created.Async, h)
	} else {
		e.Created.Sync = append(e.Created.Sync, h)
	}
	return e
}

func (e *Event) AddUpdating(h EventHandler, async ...bool) *Event {
	if e.Updating == nil {
		e.Updating = NewEventHandlers()
	}
	if len(async) > 0 && async[0] {
		e.Updating.Async = append(e.Updating.Async, h)
	} else {
		e.Updating.Sync = append(e.Updating.Sync, h)
	}
	return e
}

func (e *Event) AddUpdated(h EventHandler, async ...bool) *Event {
	if e.Updated == nil {
		e.Updated = NewEventHandlers()
	}
	if len(async) > 0 && async[0] {
		e.Updated.Async = append(e.Updated.Async, h)
	} else {
		e.Updated.Sync = append(e.Updated.Sync, h)
	}
	return e
}

func (e *Event) AddDeleting(h EventHandler, async ...bool) *Event {
	if e.Deleting == nil {
		e.Deleting = NewEventHandlers()
	}
	if len(async) > 0 && async[0] {
		e.Deleting.Async = append(e.Deleting.Async, h)
	} else {
		e.Deleting.Sync = append(e.Deleting.Sync, h)
	}
	return e
}

func (e *Event) AddDeleted(h EventHandler, async ...bool) *Event {
	if e.Deleted == nil {
		e.Deleted = NewEventHandlers()
	}
	if len(async) > 0 && async[0] {
		e.Deleted.Async = append(e.Deleted.Async, h)
	} else {
		e.Deleted.Sync = append(e.Deleted.Sync, h)
	}
	return e
}
