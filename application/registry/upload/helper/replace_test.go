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

package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFileName(t *testing.T) {
	excepted := []string{
		"/public/upload/test/0/1232aawdwwd.jpg",
		"/public/upload/test/0/1232aawdwwe.jpg",
	}
	r := ParseTemporaryFileName(`AAA` + UploadDir + `test/0/1232aawdwwd.jpgBBB` + UploadDir + `test/0/1232aawdwwe.jpgCCC`)
	assert.Equal(t, excepted, r)

	excepted = []string{
		"/public/upload/test/200/1232aawdwwd.jpg",
		"/public/upload/test/300/1232aawdwwe.jpg",
	}
	r = ParsePersistentFileName(`AAA` + UploadDir + `test/200/1232aawdwwd.jpgBBB` + UploadDir + `test/300/1232aawdwwe.jpgCCC`)
	assert.Equal(t, excepted, r)

	r = ParseAnyFileName(`AAA` + UploadDir + `test/200/1232aawdwwd.jpgBBB` + UploadDir + `test/300/1232aawdwwe.jpgCCC`)
	assert.Equal(t, excepted, r)

	content := ReplaceAnyFileName(`AAA`+UploadDir+`test/200/1232aawdwwd.jpgBBB`+UploadDir+`test/300/1232aawdwwe.jpgCCC`, func(s string) string {
		return `http://coscms.com` + s
	})
	assert.Equal(t, `AAAhttp://coscms.com`+UploadDir+`test/200/1232aawdwwd.jpgBBBhttp://coscms.com`+UploadDir+`test/300/1232aawdwwe.jpgCCC`, content)
}
