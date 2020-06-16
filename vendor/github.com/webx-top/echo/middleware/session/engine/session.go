// Session implements middleware for easily using github.com/gorilla/sessions
// within echo. This package was originally inspired from the
// https://github.com/ipfans/echo-session package, and modified to provide more
// functionality
package engine

import (
	"log"

	"github.com/admpub/sessions"
	"github.com/webx-top/echo"
)

const (
	errorFormat = "[sessions] ERROR! %s\n"
)

type Session struct {
	name    string
	context echo.Context
	store   sessions.Store
	session *sessions.Session
	written bool
	preSave []func(echo.Context) error
}

func (s *Session) AddPreSaveHook(hook func(echo.Context) error) {
	s.preSave = append(s.preSave, hook)
}

func (s *Session) SetPreSaveHook(hooks ...func(echo.Context) error) {
	s.preSave = hooks
}

func (s *Session) Get(key string) interface{} {
	return s.Session().Values[key]
}

func (s *Session) Set(key string, val interface{}) echo.Sessioner {
	s.Session().Values[key] = val
	s.setWritten()
	return s
}

func (s *Session) Delete(key string) echo.Sessioner {
	delete(s.Session().Values, key)
	s.setWritten()
	return s
}

func (s *Session) Clear() echo.Sessioner {
	for key := range s.Session().Values {
		if k, ok := key.(string); ok {
			s.Delete(k)
		}
	}
	return s
}

func (s *Session) AddFlash(value interface{}, vars ...string) echo.Sessioner {
	s.Session().AddFlash(value, vars...)
	s.setWritten()
	return s
}

func (s *Session) Flashes(vars ...string) []interface{} {
	flashes := s.Session().Flashes(vars...)
	if len(flashes) > 0 {
		s.setWritten()
	}
	return flashes
}

func (s *Session) SetID(id string) echo.Sessioner {
	s.Session().ID = id
	s.setWritten()
	return s
}

func (s *Session) ID() string {
	return s.Session().ID
}

func (s *Session) Save() error {
	if !s.Written() {
		return nil
	}
	for _, hook := range s.preSave {
		if err := hook(s.context); err != nil {
			return err
		}
	}
	err := s.Session().Save(s.context)
	if err == nil {
		s.written = false
	} else {
		log.Printf(errorFormat, err)
	}
	return err
}

func (s *Session) Session() *sessions.Session {
	if s.session == nil {
		var err error
		s.session, err = s.store.Get(s.context, s.name)
		if err != nil {
			log.Printf(errorFormat, err)
		}
	}
	return s.session
}

func (s *Session) Written() bool {
	return s.written
}

func (s *Session) setWritten() *Session {
	s.written = true
	return s
}
