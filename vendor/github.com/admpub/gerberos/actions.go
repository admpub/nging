package gerberos

import (
	"errors"
	"fmt"
	"log"
	"time"
)

type Action interface {
	Initialize(r *Rule) error
	Perform(m *Match) error
}

var actions = map[string]func() Action{
	"ban":  NewBanAction,
	"log":  NewLogAction,
	"test": NewTestAction,
}

func RegisterAction(name string, afn func() Action) {
	actions[name] = afn
}

func NewBanAction() Action {
	return &banAction{}
}

func NewLogAction() Action {
	return &logAction{}
}

func NewTestAction() Action {
	return &testAction{}
}

type banAction struct {
	rule     *Rule
	duration time.Duration
}

func (a *banAction) Initialize(r *Rule) error {
	a.rule = r

	if len(r.Action) < 2 {
		return errors.New("missing duration parameter")
	}

	d, err := time.ParseDuration(r.Action[1])
	if err != nil {
		return fmt.Errorf("failed to parse duration parameter: %w", err)
	}
	a.duration = d

	if len(r.Action) > 2 {
		return errors.New("superfluous parameter(s)")
	}

	return nil
}

func (a *banAction) Perform(m *Match) error {
	err := a.rule.runner.backend.Ban(m.IP, m.IPv6, a.duration)
	if err != nil {
		log.Printf(`%s: failed to ban IP %s: %s`, a.rule.name, m.IP, err)
	} else {
		log.Printf(`%s: banned IP %s for %s`, a.rule.name, m.IP, a.duration)
	}

	return err
}

type logAction struct {
	rule     *Rule
	extended bool
}

func (a *logAction) Initialize(r *Rule) error {
	a.rule = r

	if len(r.Action) < 2 {
		return errors.New("missing type parameter")
	}

	switch r.Action[1] {
	case "simple":
		a.extended = false
	case "extended":
		a.extended = true
	default:
		return errors.New("invalid type parameter")
	}

	if len(r.Action) > 2 {
		return errors.New("superfluous parameter(s)")
	}

	return nil
}

func (a *logAction) Perform(m *Match) error {
	var s string
	if a.extended {
		s = m.StringExtended()
	} else {
		s = m.stringSimple()
	}
	log.Printf("%s: %s", a.rule.name, s)

	return nil
}

type testAction struct {
	rule *Rule
}

func (a *testAction) Initialize(r *Rule) error {
	a.rule = r

	return nil
}

func (a *testAction) Perform(m *Match) error {
	return errFault
}
