package storer

import (
	"github.com/admpub/nging/v3/application/dbschema"
	"github.com/admpub/nging/v3/application/library/common"
	"github.com/webx-top/echo"
)

const (
	StorerInfoKey = `NgingStorer`
)

func NewInfo() *Info {
	return &Info{}
}

type Info struct {
	Name  string `json:"name" xml:"name"`
	ID    string `json:"id" xml:"id"`
	cloud *dbschema.NgingCloudStorage
}

func (s *Info) FromStore(v echo.H) *Info {
	s.Name = v.String("name")
	s.ID = v.String("id")
	if s.ID == `0` {
		s.ID = ``
	}
	return s
}

func (s *Info) Cloud(forces ...bool) *dbschema.NgingCloudStorage {
	var force bool
	if len(forces) > 0 {
		force = forces[0]
	}
	if !force && s.cloud != nil {
		return s.cloud
	}
	cloudM := &dbschema.NgingCloudStorage{}
	s.cloud = cloudM
	if len(s.ID) > 0 {
		cloudM.SetContext(common.NewMockContext())
		cloudM.Get(nil, `id`, s.ID)
	}
	return s.cloud
}
