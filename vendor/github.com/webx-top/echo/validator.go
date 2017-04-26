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

// Package echo is a fast and unfancy web framework for Go (Golang)
package echo

import (
	"errors"

	"github.com/webx-top/validation"
)

// Validator is the interface that wraps the Validate method.
type Validator interface {
	Validate(i interface{}, args ...string) error
	ValidateOk(i interface{}, args ...string) bool
	ValidateField(fieldName string, value string, rule string) bool
}

var (
	DefaultNopValidate Validator = &NopValidation{}
	ErrNoSetValidator            = errors.New(`The validator is not set`)
)

type NopValidation struct {
}

func (v *NopValidation) Validate(_ interface{}, _ ...string) error {
	return ErrNoSetValidator
}

func (v *NopValidation) ValidateOk(_ interface{}, _ ...string) bool {
	return false
}

func (v *NopValidation) ValidateField(_ string, _ string, _ string) bool {
	return false
}

func NewValidation() Validator {
	return &Validation{
		validator: validation.New(),
	}
}

type Validation struct {
	validator *validation.Validation
}

func (v *Validation) Validate(i interface{}, args ...string) error {
	_, err := v.validator.Valid(i, args...)
	return err
}

func (v *Validation) ValidateOk(i interface{}, args ...string) bool {
	ok, _ := v.validator.Valid(i, args...)
	return ok
}

func (v *Validation) ValidateField(fieldName string, value string, rule string) bool {
	return v.validator.ValidField(fieldName, value, rule)
}
