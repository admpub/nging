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

type RequestValidator func() MetaValidator

type MetaHandler struct {
	meta    H
	request RequestValidator
	Handler
}

func (m *MetaHandler) Name() string {
	if v, y := m.Handler.(Name); y {
		return v.Name()
	}
	return HandlerName(m.Handler)
}

type MetaValidator interface {
	Validate(Context) error
	Filters(Context) []FormDataFilter
}

func (m *MetaHandler) Handle(c Context) error {
	if m.request == nil {
		return m.Handler.Handle(c)
	}
	recv := m.request()
	if err := c.MustBind(recv, recv.Filters(c)...); err != nil {
		return err
	}
	if err := recv.Validate(c); err != nil {
		return err
	}
	c.Internal().Set(`validated`, recv)
	return m.Handler.Handle(c)
}

func (m *MetaHandler) Meta() H {
	return m.meta
}
