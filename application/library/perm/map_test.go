package perm

import (
	"testing"

	"github.com/admpub/nging/application/registry/navigate"
	"github.com/webx-top/echo"

	"github.com/stretchr/testify/assert"
)

var testNavigate = &navigate.List{
	{
		DisplayOnMenu: true,
		Name:          `设置`,
		Action:        `manager`,
		Icon:          `gear`,
		Children: &navigate.List{
			{
				DisplayOnMenu: true,
				Name:          `系统设置`,
				Action:        `settings`,
			},
			{
				DisplayOnMenu: true,
				Name:          `用户管理`,
				Action:        `user`,
				Children: &navigate.List{
					{
						DisplayOnMenu: true,
						Name:          `用户设置`,
						Action:        `settings`,
					},
				},
			},
			{
				DisplayOnMenu: false,
				Name:          `删除验证码`,
				Action:        `verification/delete`,
				Children: &navigate.List{
					{
						DisplayOnMenu: true,
						Name:          `用户设置`,
						Action:        `settings`,
					},
				},
			},
			{
				DisplayOnMenu: false,
				Name:          `上传图片`,
				Action:        `upload/:type`,
			},
		},
	},
}

func TestParse(t *testing.T) {
	navTree := NewMap()
	navTree.Import(testNavigate)
	m := NewMap()
	m.Parse(`manager/user,manager/settings,manager/upload/*,manager/verification/delete`, navTree)
	//echo.Dump(navTree)
	//echo.Dump(m)
	echo.Dump(m.V["manager"].V["upload"])
	echo.Dump(m.V["manager"].V["verification"])

	assert.Equal(t, navTree.V["manager"].V["user"].Nav, m.V["manager"].V["user"].Nav)
	assert.Equal(t, navTree.V["manager"].V["settings"].Nav, m.V["manager"].V["settings"].Nav)
	assert.Equal(t, navTree.V["manager"].V["verification/delete"].Nav, m.V["manager"].V["verification"].V["delete"].Nav)

	assert.True(t, m.Check(`manager/upload/*`))
	assert.True(t, m.Check(`manager/upload/:type`))
	assert.False(t, m.Check(`manager/user/settings`))
	assert.False(t, m.Check(`manager/verification/delete/settings`))
	assert.True(t, m.Check(`manager/user`))
	assert.True(t, m.Check(`manager/settings`))
	assert.True(t, m.Check(`manager/verification/delete`))
}
