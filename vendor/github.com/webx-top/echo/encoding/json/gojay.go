// +build gojay

package json

import (
	"encoding/json"

	"github.com/francoispqt/gojay"
)

var (
	MarshalIndent = json.MarshalIndent
	Marshal       = gojay.Marshal
	Unmarshal     = gojay.Unmarshal
	NewDecoder    = gojay.NewDecoder
	NewEncoder    = gojay.NewEncoder
)
