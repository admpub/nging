package factory

type Ranger interface {
	Range(func(model Model) error) error
}

func NewEvent() *Event {
	return &Event{}
}

type Event struct {
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
	case `creating`:
		return e.Creating != nil
	case `created`:
		return e.Created != nil
	case `updating`:
		return e.Updating != nil
	case `updated`:
		return e.Updated != nil
	case `deleting`:
		return e.Deleting != nil
	case `deleted`:
		return e.Deleted != nil
	}
	return false
}

func (e *Event) Call(event string, model Model, editColumns ...string) error {
	if !e.Exists(event) {
		return nil
	}
	switch event {
	case `creating`:
		return e.Creating.Exec(model, editColumns...)
	case `created`:
		return e.Created.Exec(model, editColumns...)
	case `updating`:
		return e.Updating.Exec(model, editColumns...)
	case `updated`:
		return e.Updated.Exec(model, editColumns...)
	case `deleting`:
		return e.Deleting.Exec(model, editColumns...)
	case `deleted`:
		return e.Deleted.Exec(model, editColumns...)
	}
	return nil
}

func (e *Event) On(event string, h EventHandler, async ...bool) *Event {
	switch event {
	case `creating`:
		return e.AddCreating(h, async...)
	case `created`:
		return e.AddCreated(h, async...)
	case `updating`:
		return e.AddUpdating(h, async...)
	case `updated`:
		return e.AddUpdated(h, async...)
	case `deleting`:
		return e.AddDeleting(h, async...)
	case `deleted`:
		return e.AddDeleted(h, async...)
	default:
		panic(`Unsupported event: ` + event)
	}
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
