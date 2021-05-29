package json

import (
	"strconv"
	"time"
)

type UnixTime time.Time

func (u UnixTime) MarshalJSON() ([]byte, error) {
	seconds := time.Time(u).Unix()
	return []byte(strconv.FormatInt(seconds, 10)), nil
}

type DateTime time.Time

func (u DateTime) MarshalJSON() ([]byte, error) {
	t := time.Time(u)
	if t.IsZero() {
		return nil, nil
	}
	return []byte(t.Format(`2006-01-02 15:04:05`)), nil
}
