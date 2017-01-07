package charset

import (
	"errors"

	//"github.com/admpub/chardet"
	sc "github.com/admpub/mahonia"
)

func Convert(fromEnc string, toEnc string, b []byte) ([]byte, error) {
	if !Validate(fromEnc) || !Validate(toEnc) {
		return nil, errors.New(`Unsuppored encoding.`)
	}
	dec := sc.NewDecoder(fromEnc)
	s := dec.ConvertString(string(b))
	enc := sc.NewEncoder(toEnc)
	s = enc.ConvertString(s)
	b = []byte(s)
	return b, nil
}

func Validate(enc string) bool {
	switch enc {
	case `utf-8`, `gbk`:
		return true
	}
	return false
}
