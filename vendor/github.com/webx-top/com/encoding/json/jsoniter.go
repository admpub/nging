// +build jsoniter

package json

import (
	jsoniter "github.com/json-iterator/go"
)

var (
	json          = jsoniter.ConfigCompatibleWithStandardLibrary
	MarshalIndent = json.MarshalIndent
	Marshal       = json.Marshal
	Unmarshal     = json.Unmarshal
	NewDecoder    = json.NewDecoder
	NewEncoder    = json.NewEncoder
)
