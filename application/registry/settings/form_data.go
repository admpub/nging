package settings

import (
	"github.com/admpub/nging/v3/application/dbschema"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

type DataInitor func(v *dbschema.NgingConfig) (pointer interface{}, err error)

type DataInitors map[string]DataInitor

func (d DataInitors) Register(group string) {
	for name, initor := range d {
		if len(name) > 0 {
			name = group + `.` + name
		} else {
			name = group
		}
		RegisterDecoder(name, func(v *dbschema.NgingConfig, r echo.H) error {
			jsonData, err := initor(v)
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
}

type DataFrom func(v *dbschema.NgingConfig, r echo.H) (pointer interface{}, err error)

type DataFroms map[string]DataFrom

func (d DataFroms) Register(group string) {
	for name, from := range d {
		if len(name) > 0 {
			name = group + `.` + name
		} else {
			name = group
		}
		RegisterEncoder(name, func(v *dbschema.NgingConfig, r echo.H) ([]byte, error) {
			cfg, err := from(v, r)
			if err != nil {
				return nil, err
			}
			return com.JSONEncode(cfg)
		})
	}
}
