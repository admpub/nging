// +build gojson

package json

import (
	"github.com/goccy/go-json"
)

var (
	MarshalIndent = json.MarshalIndent
	Marshal       = json.Marshal
	Unmarshal     = json.Unmarshal
	NewDecoder    = json.NewDecoder
	NewEncoder    = json.NewEncoder
)
