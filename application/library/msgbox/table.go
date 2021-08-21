/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package msgbox

import (
	"os"

	"github.com/admpub/go-pretty/v6/table"
	"github.com/admpub/go-pretty/v6/text"
)

func Table(title, data interface{}, width ...int) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	if title != nil {
		t.AppendHeader(table.Row{title})
	}
	if len(width) > 0 {
		t.SetAllowedRowLength(width[0])
	}
	switch v := data.(type) {
	case map[string]interface{}:
		for key, val := range v {
			t.AppendRow([]interface{}{key, val})
		}
	case [][]interface{}:
		for _, row := range v {
			t.AppendRow(row)
		}
	case []interface{}:
		t.AppendRow(v)
	}
	t.SetStyle(table.StyleColoredRedWhiteOnBlack)
	headerColor := text.Colors{text.BgBlue, text.FgHiWhite, text.Bold}
	t.Style().Color.Header = headerColor
	t.Style().Color.Footer = text.Colors{text.BgWhite, text.FgBlack, text.Italic}
	t.Style().Color.Row = text.Colors{text.BgWhite, text.FgBlack}
	t.Style().Color.RowAlternate = text.Colors{text.BgWhite, text.FgBlack}
	t.Render()
}
