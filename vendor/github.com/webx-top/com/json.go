package com

import (
	"fmt"
	"io"
	"os"

	"github.com/webx-top/com/encoding/json"
)

// GetJSON Read json data, writes in struct f
func GetJSON(dat *string, s interface{}) error {
	return json.Unmarshal([]byte(*dat), s)
}

// UnmarshalStream .
func UnmarshalStream(r io.Reader, m interface{}, fn func()) error {
	dec := json.NewDecoder(r)
	for {
		if err := dec.Decode(m); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		fn()
	}
	return nil
}

// MarshalStream .
func MarshalStream(w io.Writer, m interface{}) error {
	enc := json.NewEncoder(w)
	return enc.Encode(m)
}

// SetJSON Struct s will be converted to json format
func SetJSON(s interface{}) (string, error) {
	dat, err := json.Marshal(s)
	return string(dat), err
}

// ReadJSON Json data read from the specified file
func ReadJSON(path string, s interface{}) error {
	dat, err1 := os.ReadFile(path)
	if err1 != nil {
		return err1
	}
	return json.Unmarshal(dat, s)
}

// WriteJSON The json data is written to the specified file
func WriteJSON(path string, dat *string) error {
	_, err0 := os.Stat(path)
	if err0 != nil || !os.IsExist(err0) {
		os.Create(path)
	}
	return os.WriteFile(path, []byte(*dat), 0644)
}

// Dump 输出对象和数组的结构信息
func Dump(m interface{}, args ...bool) (r string) {
	v, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	r = string(v)
	l := len(args)
	if l < 1 || args[0] {
		fmt.Println(r)
	}
	return
}
