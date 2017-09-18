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
}

func (s *Session) Get(key string) interface{} {
	return s.Session().Values[key]
}

func (s *Session) Set(key string, val interface{}) echo.Sessioner {
	s.Session().Values[key] = val
	s.written = true
	return s
}

func (s *Session) Delete(key string) echo.Sessioner {
	delete(s.Session().Values, key)
	s.written = true
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
	s.written = true
	return s
}

func (s *Session) Flashes(vars ...string) []interface{} {
	s.written = true
	return s.Session().Flashes(vars...)
}

func (s *Session) SetID(id string) echo.Sessioner {
	s.Session().ID = id
	return s
}

func (s *Session) ID() string {
	return s.Session().ID
}

func (s *Session) Save() error {
	if s.Written() {
		e := s.Session().Save(s.context)
		if e == nil {
			s.written = false
		} else {
			log.Printf(errorFormat, e)
		}
		return e
	}
	return nil
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
