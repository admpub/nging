package utils

import (
	"context"
	"errors"

	"github.com/codegangsta/inject"
)

// TypeFilterConfig Filter base type interface.
type TypeFilterConfig interface {
	TypeConfig
	Event(LogEvent) LogEvent
}

// FilterConfig Filter base type struct.
type FilterConfig struct {
	CommonConfig
}

// FilterHandler type interface.
type FilterHandler interface{}

var (
	mapFilterHandler = map[string]FilterHandler{}
)

// RegistFilterHandler Registe FilterHandler.
func RegistFilterHandler(name string, handler FilterHandler) {
	mapFilterHandler[name] = handler
}

// RunFilters Run Filters
func (c *Config) RunFilters() (err error) {
	_, err = c.Injector.Invoke(c.runFilters)
	return
}

// run Filetrs.
func (c *Config) runFilters(ctx context.Context, inchan InChan, outchan OutChan) (err error) {
	filters, err := c.getFilters()
	if err != nil {
		return
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(inchan)
				close(outchan)
				return
			case event := <-inchan:
				for _, filter := range filters {
					event = filter.Event(event)
				}
				outchan <- event
			}
		}
	}()
	return
}

// get Filters.
func (c *Config) getFilters() (filters []TypeFilterConfig, err error) {
	for _, confraw := range c.FilterRaw {
		handler, ok := mapFilterHandler[confraw["type"].(string)]
		if !ok {
			err = errors.New(confraw["type"].(string))
			return
		}

		inj := inject.New()
		inj.SetParent(c)
		inj.Map(&confraw)

		refvs, err := inj.Invoke(handler)
		if err != nil {
			return []TypeFilterConfig{}, err
		}

		for _, refv := range refvs {
			if !refv.CanInterface() {
				continue
			}
			if conf, ok := refv.Interface().(TypeFilterConfig); ok {
				conf.SetInjector(inj)
				filters = append(filters, conf)
			}
		}
	}
	return
}
