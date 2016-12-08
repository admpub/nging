package ftp

import (
	"os"

	"github.com/admpub/caddyui/application/dbschema"
)

type Perm struct {
	User         *dbschema.FtpUser
	Group        *dbschema.FtpUserGroup
	defaultUser  string
	defaultGroup string
	defaultMode  os.FileMode
}

func NewPerm(user, group string, mode os.FileMode) *Perm {
	return &Perm{&dbschema.FtpUser{}, &dbschema.FtpUserGroup{}, user, group, mode}
}

func (db *Perm) GetOwner(rPath string) (string, error) {
	return db.defaultUser, nil
}

func (db *Perm) GetGroup(rPath string) (string, error) {
	return db.defaultGroup, nil
}

func (db *Perm) GetMode(rPath string) (os.FileMode, error) {
	return db.defaultMode, nil
}

func (db *Perm) ChOwner(rPath, owner string) error {
	return nil
}

func (db *Perm) ChGroup(rPath, group string) error {
	return nil
}

func (db *Perm) ChMode(rPath string, mode os.FileMode) error {
	return nil
}
