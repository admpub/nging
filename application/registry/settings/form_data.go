package settings

import (
	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

type DataDecoder func(v *dbschema.NgingConfig) (pointer interface{}, err error)

func (d DataDecoder) Register(name string) {
	RegisterDecoder(name, func(v *dbschema.NgingConfig, r echo.H) error {
		jsonData, err := d(v)
		if err != nil {
			return err
		}
		if len(v.Value) > 0 {
			err = com.JSONDecode(com.Str2bytes(v.Value), jsonData)
		}
		r[`ValueObject`] = jsonData
		return err
	})
}

type DataDecoders map[string]DataDecoder

func (d DataDecoders) Register(group string) {
	for name, initor := range d {
		if len(name) > 0 {
			name = group + `.` + name
		} else {
			name = group
		}
		initor.Register(name)
	}
}

type DataEncoder func(v *dbschema.NgingConfig, r echo.H) (pointer interface{}, err error)

func (d DataEncoder) Register(name string) {
	RegisterEncoder(name, func(v *dbschema.NgingConfig, r echo.H) ([]byte, error) {
		cfg, err := d(v, r)
		if err != nil {
			return nil, err
		}
		return com.JSONEncode(cfg)
	})
}

type DataEncoders map[string]DataEncoder

func (d DataEncoders) Register(group string) {
	for name, from := range d {
		if len(name) > 0 {
			name = group + `.` + name
		} else {
			name = group
		}
		from.Register(name)
	}
}
