/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package echo

import "fmt"

type Data struct {
	Code int
	Info interface{}
	Zone interface{} `json:",omitempty" xml:",omitempty"`
	Data interface{} `json:",omitempty" xml:",omitempty"`
}

func (d *Data) Error() string {
	return fmt.Sprintf(`%v`, d.Info)
}

func (d *Data) SetError(err error, args ...int) *Data {
	if err != nil {
		if len(args) > 0 {
			d.Code = args[0]
		} else {
			d.Code = 0
		}
		d.Info = err.Error()
	} else {
		d.Code = 1
	}
	return d
}

func (d *Data) SetCode(code int) *Data {
	d.Code = code
	return d
}

func (d *Data) SetInfo(info interface{}, args ...int) *Data {
	d.Info = info
	if len(args) > 0 {
		d.Code = args[0]
	}
	return d
}

func (d *Data) SetZone(zone interface{}) *Data {
	d.Zone = zone
	return d
}

func (d *Data) SetData(data interface{}, args ...int) *Data {
	d.Data = data
	if len(args) > 0 {
		d.Code = args[0]
	} else {
		d.Code = 1
	}
	return d
}

// NewData params: Code,Info,Zone,Data
func NewData(code int, args ...interface{}) *Data {
	var info, zone, data interface{}
	switch len(args) {
	case 3:
		data = args[2]
		fallthrough
	case 2:
		zone = args[1]
		fallthrough
	case 1:
		info = args[0]
	}
	return &Data{
		Code: code,
		Info: info,
		Zone: zone,
		Data: data,
	}
}
