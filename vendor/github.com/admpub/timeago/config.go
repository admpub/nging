package timeago

import (
	"errors"
	"time"
)

var (
	language = "ru"
	location = "" // Europe/Kiev
	loc      *time.Location
)

// Set sets configurations parameters to given value
// paramName is the name of the parameter, can be
// `language` or `location`.
// Seconds parameter is value of the configuration.
// For parameter `language` can be `ru` or `en`
func Set(paramName string, value string) error {
	switch paramName {
	case "language":
		if _, ok := translations[value]; !ok {
			return ErrUnsupported
		}
		language = value
	case "location":
		if len(value) > 0 {
			_loc, err := time.LoadLocation(value)
			if err != nil {
				return err
			}
			loc = _loc
		} else {
			loc = nil
		}
		location = value
	default:
	}
	return nil
}

var ErrUnsupported = errors.New(`Unsupported`)
