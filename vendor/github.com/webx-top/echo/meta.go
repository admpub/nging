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

type MetaValidator interface {
	Methods() []string
	Filters(Context) []FormDataFilter
}

type RequestValidator func() MetaValidator

func NewBaseRequestValidator(data interface{}, method ...string) *BaseRequestValidator {
	return &BaseRequestValidator{data: data, methods: method}
}

type BaseRequestValidator struct {
	methods []string
	data    interface{}
}

func (b *BaseRequestValidator) SetStruct(data interface{}) *BaseRequestValidator {
	b.data = data
	return b
}

func (b *BaseRequestValidator) Methods() []string {
	return b.methods
}

func (b *BaseRequestValidator) Filters(c Context) []FormDataFilter {
	return nil
}

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

func (m *MetaHandler) Meta() H {
	return m.meta
}

func (m *MetaHandler) Handle(c Context) error {
	if m.request == nil {
		return m.Handler.Handle(c)
	}
	recv := m.request()
	methods := recv.Methods()
	var data interface{}
	if bs, ok := recv.(*BaseRequestValidator); ok {
		data = bs.data
	} else {
		data = recv
	}
	if len(methods) > 0 && !InSliceFold(c.Method(), methods) {
		return m.Handler.Handle(c)
	}
	if err := c.MustBindAndValidate(data, recv.Filters(c)...); err != nil {
		return err
	}
	c.Internal().Set(`validated`, data)
	return m.Handler.Handle(c)
}

func GetValidated(c Context, defaults ...interface{}) interface{} {
	return c.Internal().Get(`validated`, defaults...)
}
