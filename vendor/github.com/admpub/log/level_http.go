package log

import (
	"fmt"

	"github.com/admpub/color"
)

func NewHttpLevel(code int, level Leveler) *httpLevel {
	lvName := HTTPStatusLevelName(code)
	lv, ok := Levels[lvName]
	if !ok {
		lv = LevelInfo
	}
	return &httpLevel{
		Code:       code,
		colorLevel: lv,
		Leveler:    level,
	}
}

type httpLevel struct {
	Code       int     // HTTP Status Code
	colorLevel Leveler // 用于显示颜色
	Leveler            // 用于判断是否显示
}

func (h httpLevel) Tag() string {
	return `[ ` + fmt.Sprint(h.Code) + ` ]`
}

func (l httpLevel) Color() *color.Color {
	return colorBrushes[l.Leveler]
}
