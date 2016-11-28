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
func (d *Data) SetCode(code int) *Data {
	d.Code = code
	return d
}
func (d *Data) SetInfo(info interface{}) *Data {
	d.Info = info
	return d
}
func (d *Data) SetZone(zone interface{}) *Data {
	d.Zone = zone
	return d
}
func (d *Data) SetData(data interface{}) *Data {
	d.Data = data
	return d
}

// NewData params: CIZD
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
