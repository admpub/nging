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
	"fmt"
)

var DefaultNopSession Sessioner = &NopSession{}

// Options stores configuration for a session or session store.
// Fields are a subset of http.Cookie fields.
type SessionOptions struct {
	Engine string //Store Engine
	Name   string //Session Name
	*CookieOptions
}

// Wraps thinly gorilla-session methods.
// Session stores the values and optional configuration for a session.
type Sessioner interface {
	// Get returns the session value associated to the given key.
	Get(key string) interface{}
	// Set sets the session value associated to the given key.
	Set(key string, val interface{}) Sessioner
	SetId(id string) Sessioner
	Id() string
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

	Options(SessionOptions) Sessioner

	// Save saves all sessions used during the current request.
	Save() error
}

type NopSession struct {
}

func (n *NopSession) Get(name string) interface{} {
	fmt.Println(`NopSession#Get:`, name)
	return nil
}

func (n *NopSession) Set(name string, value interface{}) Sessioner {
	fmt.Println(`NopSession#Set:`, name, `Value:`, value)
	return n
}

func (n *NopSession) SetId(id string) Sessioner {
	fmt.Println(`NopSession#SetId:`, id)
	return n
}

func (n *NopSession) Id() string {
	fmt.Println(`NopSession#Id()`)
	return ``
}

func (n *NopSession) Delete(name string) Sessioner {
	fmt.Println(`NopSession#Delete:`, name)
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
