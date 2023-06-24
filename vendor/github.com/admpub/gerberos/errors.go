package gerberos

import "errors"

var (
	ErrMissingSource            = errors.New("missing source")
	ErrEmptySource              = errors.New("empty source")
	ErrUnknownSource            = errors.New("unknown source")
	ErrMissingAction            = errors.New("missing action")
	ErrEmptyAction              = errors.New("empty action")
	ErrUnknownAction            = errors.New("unknown action")
	ErrMissingIntervalParameter = errors.New("missing interval parameter")
	ErrInvalidIntervalParameter = errors.New("failed to parse interval parameter")
	ErrMissingRegexp            = errors.New("missing regexp")
	ErrEmptyRegexp              = errors.New("empty regexp")
	ErrMissingCountParameter    = errors.New("missing count parameter")
	ErrInvalidCountParameter    = errors.New("failed to parse count parameter")
)
