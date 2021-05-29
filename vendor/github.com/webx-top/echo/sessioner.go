/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

package echo

import (
	"log"
)

var (
	DefaultNopSession     Sessioner = &NopSession{}
	DefaultDebugSession   Sessioner = &DebugSession{}
	DefaultSession                  = DefaultNopSession
	DefaultSessionOptions           = NewSessionOptions(`cookie`, `SID`)
)

func NewSessionOptions(engine string, name string, args ...*CookieOptions) *SessionOptions {
	cookieOptions := DefaultCookieOptions
	if len(args) > 0 {
		cookieOptions = args[0]
	}
	return &SessionOptions{
		Engine:        engine,
		Name:          name,
		CookieOptions: cookieOptions,
	}
}

// SessionOptions stores configuration for a session or session store.
// Fields are a subset of http.Cookie fields.
type SessionOptions struct {
	Engine string //Store Engine
	Name   string //Session Name
	*CookieOptions
}

func (s *SessionOptions) Clone() *SessionOptions {
	clone := *s
	if s.CookieOptions != nil {
		clone.CookieOptions = s.CookieOptions.Clone()
	}
	return &clone
}

// Sessioner Wraps thinly gorilla-session methods.
// Session stores the values and optional configuration for a session.
type Sessioner interface {
	// Get returns the session value associated to the given key.
	Get(key string) interface{}
	// Set sets the session value associated to the given key.
	Set(key string, val interface{}) Sessioner
	SetID(id string, notReload ...bool) error
	ID() string
	MustID() string
	RemoveID(sessionID string) error
	// Delete removes the session value associated to the given key.
	Delete(key string) Sessioner
	// Clear deletes all values in the session.
	Clear() Sessioner
	// AddFlash adds a flash message to the session.
	// A single variadic argument is accepted, and it is optional: it defines the flash key.
	// If not defined "_flash" is used by default.
	AddFlash(value interface{}, vars ...string) Sessioner
	// Flashes returns a slice of flash messages from the session.
	// A single variadic argument is accepted, and it is optional: it defines the flash key.
	// If not defined "_flash" is used by default.
	Flashes(vars ...string) []interface{}
	// Save saves all sessions used during the current request.
	Save() error
	AddPreSaveHook(func(Context) error)
	SetPreSaveHook(...func(Context) error)
}

type NopSession struct {
}

func (n *NopSession) Get(name string) interface{} {
	return nil
}

func (n *NopSession) Set(name string, value interface{}) Sessioner {
	return n
}

func (n *NopSession) SetID(id string, notReload ...bool) error {
	return nil
}

func (n *NopSession) ID() string {
	return ``
}

func (n *NopSession) MustID() string {
	return ``
}

func (n *NopSession) RemoveID(sessionID string) error {
	return nil
}

func (n *NopSession) Delete(name string) Sessioner {
	return n
}

func (n *NopSession) Clear() Sessioner {
	return n
}

func (n *NopSession) AddFlash(_ interface{}, _ ...string) Sessioner {
	return n
}

func (n *NopSession) Flashes(_ ...string) []interface{} {
	return []interface{}{}
}

func (n *NopSession) Options(_ SessionOptions) Sessioner {
	return n
}

func (n *NopSession) Save() error {
	return nil
}

func (n *NopSession) AddPreSaveHook(_ func(Context) error) {
}

func (n *NopSession) SetPreSaveHook(_ ...func(Context) error) {
}

type DebugSession struct {
}

func (n *DebugSession) Get(name string) interface{} {
	log.Println(`DebugSession.Get`, name)
	return nil
}

func (n *DebugSession) Set(name string, value interface{}) Sessioner {
	log.Println(`DebugSession.Set`, name, value)
	return n
}

func (n *DebugSession) SetID(id string, notReload ...bool) error {
	log.Println(`DebugSession.SetID`, id, `notReload:`, notReload)
	return nil
}

func (n *DebugSession) ID() string {
	log.Println(`DebugSession.ID`)
	return ``
}

func (n *DebugSession) MustID() string {
	log.Println(`DebugSession.MustID`)
	return ``
}

func (n *DebugSession) RemoveID(sessionID string) error {
	log.Println(`DebugSession.RemoveID`, sessionID)
	return nil
}

func (n *DebugSession) Delete(name string) Sessioner {
	log.Println(`DebugSession.Delete`, name)
	return n
}

func (n *DebugSession) Clear() Sessioner {
	log.Println(`DebugSession.Clear`)
	return n
}

func (n *DebugSession) AddFlash(value interface{}, args ...string) Sessioner {
	log.Println(`DebugSession.AddFlash`, value, args)
	return n
}

func (n *DebugSession) Flashes(args ...string) []interface{} {
	log.Println(`DebugSession.Flashes`, args)
	return []interface{}{}
}

func (n *DebugSession) Options(options SessionOptions) Sessioner {
	log.Println(`DebugSession.Options`, options)
	return n
}

func (n *DebugSession) Save() error {
	log.Println(`DebugSession.Save`)
	return nil
}

func (n *DebugSession) AddPreSaveHook(_ func(Context) error) {
}

func (n *DebugSession) SetPreSaveHook(_ ...func(Context) error) {
}
