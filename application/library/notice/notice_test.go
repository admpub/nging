/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

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

package notice

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo/defaults"
)

func init() {
	SetDebug(true)
}

func TestOpenMessage(t *testing.T) {
	OpenMessage(`testUser`, `testType`)
	user, _ := Default().users.GetOk(`testUser`)
	assert.Equal(t, 1, user.Notice.types.Size())
	assert.True(t, user.Notice.types.Has(`testType`))

	_, clientID := OpenClient(`testUser`)
	assert.Equal(t, 1, user.CountClient())

	CloseClient(`testUser`, clientID)
	assert.Equal(t, 0, user.CountClient())
}

func TestNoticeProgress(t *testing.T) {
	ctx := context.Background()
	eCtx := defaults.NewMockContext()
	noticer := NewP(eCtx, `databaseImport`, `username`, ctx)
	noticer.AutoComplete(true)
	noticer.Add(2)
	assert.Equal(t, int64(2), noticer.prog.Total)
	noticer.Done(1)
	assert.Equal(t, int64(1), noticer.prog.Finish)
	noticer.Done(1)
	assert.Equal(t, int64(2), noticer.prog.Finish)
	assert.True(t, noticer.prog.Complete)
}
