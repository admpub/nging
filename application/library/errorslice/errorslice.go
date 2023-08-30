package errorslice

import (
	"strings"
)

const (
	NL  = "\n"
	NLT = "\n\t"
)

func New() Errors {
	return Errors{}
}

// Errors 多个错误信息
type Errors []error

func (e Errors) Error() string {
	return e.Stringify(NL)
}

func (e Errors) ErrorTab() string {
	return e.Stringify(NLT)
}

func (e Errors) Stringify(separator string) string {
	s := make([]string, len(e))
	for k, v := range e {
		s[k] = v.Error()
	}
	return strings.Join(s, separator)
}

func (e Errors) String() string {
	return e.Error()
}

func (e Errors) IsEmpty() bool {
	return len(e) == 0
}

func (e *Errors) Add(err error) {
	*e = append(*e, err)
}

func (e Errors) Unwrap() []error {
	return []error(e)
}

func (e Errors) ToError() error {
	if e.IsEmpty() {
		return nil
	}
	return e
}
