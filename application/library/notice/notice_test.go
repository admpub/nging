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

package notice

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	SetDebug(true)
}

func TestOpenMessage(t *testing.T) {
	OpenMessage(`testUser`, `testType`)
	user := Default().user[`testUser`]
	assert.Equal(t, 1, user.Notice.types.Size())
	assert.True(t, user.Notice.types.Has(`testType`))

	clientID := OpenClient(`testUser`)
	assert.Equal(t, 1, user.CountClient())

	CloseClient(`testUser`, clientID)
	assert.Equal(t, 0, user.CountClient())
}

func TestSend(t *testing.T) {
	OpenMessage(`testUser`, `testType`)
	clientID := OpenClient(`testUser`)
	var wg sync.WaitGroup
	go func() {
		defer wg.Done()
		b, err := RecvJSON(`testUser`, clientID)
		if err != nil {
			t.Error(err)
		}
		println(string(b))
		assert.Equal(t, `{"client_id":"`+clientID+`","id":null,"type":"testType","title":"testTitle","status":1,"content":"testContent","mode":"","progress":null}`, string(b))
	}()
	wg.Add(1)
	Send(`testUser`, NewMessageWithValue(`testType`, `testTitle`, `testContent`).SetClientID(clientID))
	wg.Wait()
}
