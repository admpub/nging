package subdir

import "strings"

type Params struct {
	Subdir     string
	Table      string
	Field      string
	SubdirInfo *SubdirInfo
	TableFlag  bool // uploadType中是否传递的表名称的方式：table.field
}

func (p *Params) MustGetSubdir() string {
	subdir := p.Subdir
	if len(subdir) == 0 {
		subdir = p.Table
	}
	return subdir
}

func (p *Params) MustGetTable() string {
	table := p.Table
	if len(table) == 0 {
		table = p.SubdirInfo.TableName()
	}
	return table
}

func (p *Params) IsAllowed() bool {
	if p.SubdirInfo == nil {
		return false
	}
	return p.SubdirInfo.Allowed
}

func NewParams() *Params {
	return &Params{}
}

func getSeperator(content string) string {
	for _, v := range content {
		if v == '.' {
			return `.`
		}
		if v == ':' {
			return `:`
		}
	}
	return ``
}

// ParseUploadType 根据updateType值获取SubdirInfo数据
// uploadType: subdir:fieldName
func ParseUploadType(uploadType string) *Params {
	params := NewParams()
	seperator := getSeperator(uploadType)
	if len(seperator) > 0 {
		types := strings.SplitN(uploadType, seperator, 2)
		switch len(types) {
		case 2:
			params.Field = types[1]
			fallthrough
		case 1:
			if seperator == `.` { // table.field
				params.Table = types[0]
				params.TableFlag = true
			} else { // subdir:field
				params.Subdir = types[0]
			}
		}
	} else {
		params.Subdir = uploadType
	}
	if len(params.Subdir) > 0 {
		params.SubdirInfo, _ = subdirs[params.Subdir]
	} else {
		params.SubdirInfo = GetByTable(params.Table)
		if params.SubdirInfo != nil {
			params.Subdir = params.SubdirInfo.Key
		}
	}
	return params
}
