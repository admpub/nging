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
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestSend(t *testing.T) {
	OpenMessage(`testUser`, `testType`)
	_, clientID := OpenClient(`testUser`)
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

func TestSend2(t *testing.T) {
	OpenMessage(`testUser2`, `testType`)
	_, clientID := OpenClient(`testUser2`)
	var wg sync.WaitGroup
	go func() {
		var i = 0
		for {
			b, err := RecvJSON(`testUser`, clientID)
			if err != nil {
				t.Error(err)
			}
			println(string(b))
			assert.Equal(t, `{"client_id":"`+clientID+`","id":null,"type":"testType","title":"testTitle","status":1,"content":"testContent_`+strconv.Itoa(i)+`","mode":"","progress":null}`, string(b))
			i++
		}
	}()
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			Send(`testUser`, NewMessageWithValue(`testType`, `testTitle`, `testContent_`+strconv.Itoa(i)).SetClientID(clientID))
		}(i)
	}
	wg.Wait()
}
