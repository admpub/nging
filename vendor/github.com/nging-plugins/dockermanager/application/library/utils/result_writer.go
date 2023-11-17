package utils

import (
	"encoding/json"
)

func NewResultWriter(cb func(*Result) error) *resultWriter {
	return &resultWriter{
		buf: []byte{},
		cb:  cb,
	}
}

type resultWriter struct {
	buf []byte
	cb  func(*Result) error
}

func (r *resultWriter) Write(p []byte) (int, error) {
	var err error
	for _, v := range p {
		if v == '\n' {
			result := &Result{}
			err = json.Unmarshal(r.buf, result)
			if err != nil {
				return 0, err
			}
			if err = r.cb(result); err != nil {
				return 0, err
			}
			r.buf = r.buf[0:0]
			continue
		}
		r.buf = append(r.buf, v)
	}
	return len(p), err
}

func (r *resultWriter) Flush() error {
	if len(r.buf) == 0 {
		return nil
	}
	result := &Result{}
	err := json.Unmarshal(r.buf, result)
	if err != nil {
		return err
	}
	result.SetCompleted(true)
	if err = r.cb(result); err != nil {
		return err
	}
	r.buf = r.buf[0:0:0]
	return err
}
