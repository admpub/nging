package perm

import (
	"testing"

	"github.com/admpub/nging/v3/application/registry/navigate"
	"github.com/webx-top/echo"

	"github.com/stretchr/testify/assert"
)

var testNavigate = &navigate.List{
	{
		Display: true,
		Name:    `设置`,
		Action:  `manager`,
		Icon:    `gear`,
		Children: &navigate.List{
			{
				Display: true,
				Name:    `系统设置`,
				Action:  `settings`,
			},
			{
				Display: true,
				Name:    `用户管理`,
				Action:  `user`,
				Children: &navigate.List{
					{
						Display: true,
						Name:    `用户设置`,
						Action:  `settings`,
					},
				},
			},
			{
				Display: false,
				Name:    `删除验证码`,
				Action:  `verification/delete`,
				Children: &navigate.List{
					{
						Display: true,
						Name:    `用户设置`,
						Action:  `settings`,
					},
					{
						Display:   true,
						Name:      `任何人有权2`,
						Action:    `everone2`,
						Unlimited: true,
					},
				},
			},
			{
				Display: false,
				Name:    `上传图片`,
				Action:  `upload/:type`,
			},
			{
				Display:   false,
				Name:      `任何人有权`,
				Action:    `everone`,
				Unlimited: true,
			},
		},
	},
}

func TestParse(t *testing.T) {
	navTree := NewMap()
	navTree.Import(testNavigate)
	m := NewMap()
	m.Parse(`manager/user,manager/settings,manager/upload/*,manager/verification/delete`, navTree)
	echo.Dump(navTree)
	//echo.Dump(m)
	//echo.Dump(m.V["manager"].V["upload"])
	//echo.Dump(m.V["manager"].V["verification"])

	assert.Equal(t, navTree.V["manager"].V["user"].Nav, m.V["manager"].V["user"].Nav)
	assert.Equal(t, navTree.V["manager"].V["settings"].Nav, m.V["manager"].V["settings"].Nav)
	assert.Equal(t, navTree.V["manager"].V["verification/delete"].Nav, m.V["manager"].V["verification"].V["delete"].Nav)

	assert.True(t, m.Check(`manager/verification/delete`, navTree))
	assert.True(t, m.Check(`manager/upload/:type`, navTree))

	assert.True(t, m.Check(`manager/verification/delete/everone2`, navTree))
	assert.True(t, m.Check(`manager/everone`, navTree))
	assert.False(t, m.Check(`manager/user/settings`, navTree))
	assert.False(t, m.Check(`manager/verification/delete/settings`, navTree))

	assert.True(t, m.Check(`manager/upload/*`, navTree))
	assert.True(t, m.Check(`manager/user`, navTree))
	assert.True(t, m.Check(`manager/settings`, navTree))
	assert.False(t, m.Check(`manager/verification/delete/settings`, navTree))
}
