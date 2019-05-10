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
