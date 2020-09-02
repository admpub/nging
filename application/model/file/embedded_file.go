package file

import (
	"strings"

	"github.com/webx-top/com"

	"github.com/admpub/nging/application/dbschema"
)

func (f *Embedded) DeleteByInstance(m *dbschema.NgingFileEmbedded) error {
	err := f.Delete(nil, `id`, m.Id)
	if err != nil {
		return err
	}

	ids := strings.Split(m.FileIds, ",")
	return f.DeleteFileByIds(ids)
}

func (f *Embedded) AddFileByIds(fileIds []string, excludeFileIds ...string) error {
	var newIds []interface{}
	if len(excludeFileIds) > 0 {
		for _, v := range fileIds {
			if !com.InSlice(v, excludeFileIds) {
				newIds = append(newIds, v)
			}
		}
	} else {
		newIds = make([]interface{}, len(fileIds))
		for idx, id := range fileIds {
			newIds[idx] = id
		}
	}
	if len(newIds) == 0 {
		return nil
	}
	return f.File.Incr(newIds...)
}

func (f *Embedded) DeleteFileByIds(fileIds []string, excludeFileIds ...string) error {
	var delIds []interface{}
	if len(excludeFileIds) > 0 {
		for _, v := range fileIds {
			if !com.InSlice(v, excludeFileIds) {
				delIds = append(delIds, v)
			}
		}
	} else {
		delIds = make([]interface{}, len(fileIds))
		for idx, id := range fileIds {
			delIds[idx] = id
		}
	}
	if len(delIds) == 0 {
		return nil
	}
	err := f.File.Decr(delIds...)
	return err
}
