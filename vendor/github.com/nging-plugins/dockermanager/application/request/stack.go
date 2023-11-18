package request

import (
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

type StackAdd struct {
	Name    string `validate:"required"`
	WorkDir string
	File    string
	Content string
}

func (c *StackAdd) BeforeValidate(ctx echo.Context) error {
	c.Name = strings.TrimSpace(c.Name)
	c.WorkDir = strings.TrimSpace(c.WorkDir)
	c.File = strings.TrimSpace(c.File)
	if len(c.File) > 0 {
		if !com.FileExists(c.File) {
			return ctx.NewError(code.DataNotFound, `文件%q没有找到`, c.File).SetZone(`file`)
		}
	} else if len(c.Content) == 0 {
		return ctx.NewError(code.InvalidParameter, `YAML文件路径和YAML内容不能同时为空`).SetZone(`content`)
	}
	return nil
}

func (c *StackAdd) AfterValidate(ctx echo.Context) error {
	if !com.IsAlphaNumericUnderscoreHyphen(c.Name) {
		return ctx.NewError(code.InvalidParameter, `名称%q只能由字母、数字、下划线或短横“-”组成`, c.File).SetZone(`name`)
	}
	return nil
}

type StackEdit struct {
	Content string `validate:"required"`
	WorkDir string
}

func (c *StackEdit) BeforeValidate(ctx echo.Context) error {
	c.WorkDir = strings.TrimSpace(c.WorkDir)
	return nil
}
