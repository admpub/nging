package webdav

import (
	"database/sql"
	"strings"

	"github.com/webx-top/com"
)

type WebdavPerm struct {
	Readable  sql.NullBool
	Writeable sql.NullBool
	Resource  string
}

func (w *WebdavPerm) String() string {
	var readable, writable string
	if w.Readable.Bool {
		readable = `allow`
	} else {
		readable = `block`
	}
	if w.Writeable.Valid {
		if w.Writeable.Bool {
			writable = `+w`
		} else {
			writable = `-w`
		}
	}
	var perm string
	if strings.Contains(w.Resource, `*`) {
		perm = readable + `_r    "` + com.AddCSlashes(strings.Replace(w.Resource, `*`, `(.*)`, -1), '"') + `"`
	} else if strings.Contains(w.Resource, `|`) {
		perm = readable + `_r    "` + com.AddCSlashes(w.Resource, '"') + `"`
	} else {
		perm = readable + `      "` + com.AddCSlashes(w.Resource, '"') + `"`
	}
	if len(writable) > 0 {
		perm += `      ` + writable
	}
	return perm
}

func (w *WebdavPerm) SetWriteable(v string) {
	switch v {
	case "0":
		w.Writeable.Valid = true
		w.Writeable.Bool = false
	case "1":
		w.Writeable.Valid = true
		w.Writeable.Bool = true
	}
}

func (w *WebdavPerm) SetReadable(r string) {
	switch r {
	case "0":
		w.Readable.Valid = true
		w.Readable.Bool = false
	case "1":
		w.Readable.Valid = true
		w.Readable.Bool = true
	}
}

type WebdavUser struct {
	User      string
	Password  string
	Root      string
	Writeable sql.NullBool
	Perms     []*WebdavPerm
}

func (w *WebdavUser) SetWriteable(v string) {
	switch v {
	case "0":
		w.Writeable.Valid = true
		w.Writeable.Bool = false
	case "1":
		w.Writeable.Valid = true
		w.Writeable.Bool = true
	}
}
