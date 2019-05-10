package utils

import (
	"context"
	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/inject"
)

// Output base type interface.
type TypeOutputConfig interface {
	TypeConfig
	Event(ctx context.Context, event LogEvent) (err error)
}

// Output base type struct.
type OutputConfig struct {
	CommonConfig
}

// OutputHandler type interface.
type OutputHandler interface{}

var (
	mapOutputHandler = map[string]OutputHandler{}
)

// Registe OutputHandler.
func RegistOutputHandler(name string, handler OutputHandler) {
	mapOutputHandler[name] = handler
}

// Run Outputs.
func (c *Config) RunOutputs() (err error) {
	_, err = c.Injector.Invoke(c.runOutputs)
	return
}

// run Outputs.
func (c *Config) runOutputs(ctx context.Context, outchan OutChan, logger *logrus.Logger) (err error) {
	outputs, err := c.getOutputs()
	if err != nil {
		return
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-outchan:
				for _, output := range outputs {
					go func(o TypeOutputConfig, e LogEvent) {
						if err = o.Event(ctx, e); err != nil {
							logger.Errorf("output failed: %v\n", err)
						}
					}(output, event)
				}
			}
		}
	}()
	return
}

// get Outputs.
func (c *Config) getOutputs() (outputs []TypeOutputConfig, err error) {
	for _, confraw := range c.OutputRaw {
		handler, ok := mapOutputHandler[confraw["type"].(string)]
		if !ok {
			err = errors.New(confraw["type"].(string))
			return
		}

		inj := inject.New()
		inj.SetParent(c)
		inj.Map(&confraw)

		refvs, err := inj.Invoke(handler)
		if err != nil {
			return []TypeOutputConfig{}, err
		}

		for _, refv := range refvs {
			if !refv.CanInterface() {
				continue
			}
			if conf, ok := refv.Interface().(TypeOutputConfig); ok {
				conf.SetInjector(inj)
				outputs = append(outputs, conf)
			}
		}
	}
	return
}
