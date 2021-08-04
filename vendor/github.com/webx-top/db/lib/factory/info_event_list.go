package factory

import "github.com/webx-top/db"

func NewEvents() Events {
	return Events{}
}

type Events map[string]*Event

func (e Events) getEventList(event string, table string) []*Event {
	events := []*Event{}
	if evt, ok := e[table]; ok {
		if evt.Exists(event) {
			events = append(events, evt)
		}
	}
	if evt, ok := e[`*`]; ok {
		if evt.Exists(event) {
			events = append(events, evt)
		}
	}
	return events
}

func (e Events) Exists(event string, model Model) bool {
	table := model.Short_()
	if evt, ok := e[table]; ok {
		if evt.Exists(event) {
			return true
		}
	}
	if evt, ok := e[`*`]; ok {
		if evt.Exists(event) {
			return true
		}
	}
	return false
}

func (e Events) Call(event string, model Model, editColumns []string, mw func(db.Result) db.Result, args ...interface{}) error {
	if event == EventDeleted {
		return e.call(event, model)
	}
	if len(args) == 0 {
		return e.call(event, model, editColumns...)
	}
	events := e.getEventList(event, model.Short_())
	if len(events) == 0 {
		return nil
	}
	rows := model.NewObjects()
	num := int64(1000)
	cnt, err := model.ListByOffset(rows, mw, 0, int(num), args...)
	if err != nil {
		return err
	}
	total := cnt()
	if total < 1 {
		return nil
	}
	kvset := map[string]interface{}{}
	if len(editColumns) > 0 {
		rowM := model.AsRow()
		for _, key := range editColumns {
			kvset[key] = rowM[key]
		}
	}
	for i := int64(0); i < total; i += num {
		if i > 0 {
			rows = model.NewObjects()
			_, err = model.ListByOffset(rows, mw, int(i), int(num), args...)
			if err != nil {
				return err
			}
		}
		err = rows.Range(func(m Model) error {
			m.CPAFrom(model).FromRow(kvset)
			for _, evt := range events {
				if err := evt.Call(event, m, editColumns...); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (e Events) call(event string, model Model, editColumns ...string) error {
	for _, evt := range e.getEventList(event, model.Short_()) {
		err := evt.Call(event, model, editColumns...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e Events) CallRead(event string, model Model, param *Param, rangers ...Ranger) error {
	table := model.Short_()
	if len(rangers) < 1 { // 单行数据
		if evt, ok := e[table]; ok {
			err := evt.CallRead(event, model, param)
			if err != nil {
				return err
			}
		}
		if evt, ok := e[`*`]; ok {
			return evt.CallRead(event, model, param)
		}
		return nil
	}
	if evt, ok := e[table]; ok {
		err := rangers[0].Range(func(m Model) error {
			m.CPAFrom(model)
			return evt.CallRead(event, m, param)
		})
		if err != nil {
			return err
		}
	}
	if evt, ok := e[`*`]; ok {
		err := rangers[0].Range(func(m Model) error {
			m.CPAFrom(model)
			return evt.CallRead(event, m, param)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Events) On(event string, h EventHandler, table string, async ...bool) {
	evt, ok := (*e)[table]
	if !ok {
		evt = NewEvent()
		(*e)[table] = evt
	}
	evt.On(event, h, async...)
}

func (e *Events) OnRead(event string, h EventReadHandler, table string, async ...bool) {
	evt, ok := (*e)[table]
	if !ok {
		evt = NewEvent()
		(*e)[table] = evt
	}
	evt.OnRead(event, h, async...)
}
