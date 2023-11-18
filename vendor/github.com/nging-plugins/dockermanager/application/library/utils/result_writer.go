/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

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
