package factory

import "github.com/webx-top/db"

type Ranger interface {
	Range(func(model Model) error) error
}

func NewEvent() *Event {
	return &Event{}
}

func NewEvents() Events {
	return Events{}
}

type Events map[string]*Event

func (e Events) Call(event string, model Model, mw func(db.Result) db.Result, args ...interface{}) error {
	if len(args) > 0 {
		table := model.Short_()
		events := []*Event{}
		if evt, ok := e[table]; ok {
			events = append(events, evt)
		}
		if evt, ok := e[`*`]; ok {
			events = append(events, evt)
		}
		if len(events) == 0 {
			return nil
		}
		rows := model.NewObjects()
		num := int64(1000)
		cnt, err := model.ListByOffset(&rows, mw, 0, int(num), args...)
		if err != nil {
			return err
		}
		total := cnt()
		if total < 1 {
			return nil
		}
		for i := int64(0); i < total; i += num {
			if i > 0 {
				rows = model.NewObjects()
				_, err := model.ListByOffset(&rows, mw, int(i), int(num), args...)
				if err != nil {
					return err
				}
			}
			return rows.Range(func(m Model) error {
				for _, evt := range events {
					if err := evt.Call(event, m); err != nil {
						return err
					}
				}
				return nil
			})
		}
	}
	return e.call(event, model)
}

func (e Events) call(event string, model Model) error {
	table := model.Short_()
	if evt, ok := e[table]; ok {
		err := evt.Call(event, model)
		if err != nil {
			return err
		}
	}
	if evt, ok := e[`*`]; ok {
		return evt.Call(event, model)
	}
	return nil
}

func (e *Events) Add(event string, h EventHandler, table string) {
	evt, ok := (*e)[table]
	if !ok {
		evt = NewEvent()
		(*e)[table] = evt
	}
	evt.Add(event, h)
}

type EventHandler func(model Model) error

type Event struct {
	// Creating 创建之前
	Creating []EventHandler
	// Created 创建之后
	Created []EventHandler

	// Updating 更新之前
	Updating []EventHandler
	// Updated 更新之后
	Updated []EventHandler

	// Deleting 删除之前
	Deleting []EventHandler
	// Deleted 删除之后
	Deleted []EventHandler
}

func (e *Event) Call(event string, model Model) error {
	switch event {
	case `creating`:
		for _, handler := range e.Creating {
			if err := handler(model); err != nil {
				return err
			}
		}
	case `created`:
		for _, handler := range e.Created {
			if err := handler(model); err != nil {
				return err
			}
		}
	case `updating`:
		for _, handler := range e.Updating {
			if err := handler(model); err != nil {
				return err
			}
		}
	case `updated`:
		for _, handler := range e.Updated {
			if err := handler(model); err != nil {
				return err
			}
		}
	case `deleting`:
		for _, handler := range e.Deleting {
			if err := handler(model); err != nil {
				return err
			}
		}
	case `deleted`:
		for _, handler := range e.Deleted {
			if err := handler(model); err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *Event) Add(event string, h EventHandler) *Event {
	switch event {
	case `creating`:
		return e.AddCreating(h)
	case `created`:
		return e.AddCreated(h)
	case `updating`:
		return e.AddUpdating(h)
	case `updated`:
		return e.AddUpdated(h)
	case `deleting`:
		return e.AddDeleting(h)
	case `deleted`:
		return e.AddDeleted(h)
	}
	return e
}

func (e *Event) AddCreating(h EventHandler) *Event {
	e.Creating = append(e.Creating, h)
	return e
}

func (e *Event) AddCreated(h EventHandler) *Event {
	e.Created = append(e.Created, h)
	return e
}

func (e *Event) AddUpdating(h EventHandler) *Event {
	e.Updating = append(e.Updating, h)
	return e
}

func (e *Event) AddUpdated(h EventHandler) *Event {
	e.Updated = append(e.Updated, h)
	return e
}

func (e *Event) AddDeleting(h EventHandler) *Event {
	e.Deleting = append(e.Deleting, h)
	return e
}

func (e *Event) AddDeleted(h EventHandler) *Event {
	e.Deleted = append(e.Deleted, h)
	return e
}
