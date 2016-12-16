package driver

import "encoding/gob"

type DbAuth struct {
	Driver   string
	Username string
	Password string
	Host     string
	Db       string
}

func (d *DbAuth) CopyFrom(auth *DbAuth) *DbAuth {
	d.Driver = auth.Driver
	d.Username = auth.Username
	d.Password = auth.Password
	d.Host = auth.Host
	d.Db = auth.Db
	return d
}

func init() {
	gob.Register(&DbAuth{})
}
