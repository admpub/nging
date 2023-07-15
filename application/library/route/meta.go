package route

import "github.com/webx-top/echo"

type MetaSetter interface {
	SetMeta(meta echo.H) MetaSetter
	SetMetaKV(key string, val interface{}) MetaSetter
}

func newMeta(groupName string, groupInst *Group) meta {
	return meta{groupName: groupName, groupInst: groupInst}
}

type meta struct {
	groupName string
	groupInst *Group
}

func (m meta) SetMeta(meta echo.H) MetaSetter {
	m.groupInst.SetMeta(m.groupName, meta)
	return m
}

func (m meta) SetMetaKV(key string, val interface{}) MetaSetter {
	m.groupInst.SetMetaKV(m.groupName, key, val)
	return m
}
