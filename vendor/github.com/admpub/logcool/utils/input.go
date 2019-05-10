package utils

import (
	"errors"

	"github.com/codegangsta/inject"
)

// Input base type interface.
type TypeInputConfig interface {
	TypeConfig
	Start()
}

// Input base type struct.
type InputConfig struct {
	CommonConfig
}

// InputHandler type interface.
type InputHandler interface{}

var (
	mapInputHandler = map[string]InputHandler{}
)

// Registe InputHandler.
func RegistInputHandler(name string, handler InputHandler) {
	mapInputHandler[name] = handler
}

// Run Inputs.
func (c *Config) RunInputs() (err error) {
	_, err = c.Injector.Invoke(c.runInputs)
	return
}

// run Inputs.
func (c *Config) runInputs(inchan InChan) (err error) {
	inputs, err := c.getInputs(inchan)
	if err != nil {
		return
	}

	for _, input := range inputs {
		go input.Start()
	}
	return
}

// get Inputs.
func (c *Config) getInputs(inchan InChan) (inputs []TypeInputConfig, err error) {
	for _, confraw := range c.InputRaw {
		handler, ok := mapInputHandler[confraw["type"].(string)]
		if !ok {
			err = errors.New(confraw["type"].(string))
			return
		}

		inj := inject.New()
		inj.SetParent(c)
		inj.Map(&confraw)
		inj.Map(inchan)

		refvs, err := inj.Invoke(handler)
		if err != nil {
			return []TypeInputConfig{}, err
		}

		for _, refv := range refvs {
			if !refv.CanInterface() {
				continue
			}
			if conf, ok := refv.Interface().(TypeInputConfig); ok {
				conf.SetInjector(inj)
				inputs = append(inputs, conf)
			}
		}
	}
	return
}
