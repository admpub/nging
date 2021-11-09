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
)

func newNoticeMessages() *noticeMessages {
	return &noticeMessages{messages: map[string]chan *Message{}}
}

type noticeMessages struct {
	messages map[string]chan *Message
	lock     sync.RWMutex
}

func (n *noticeMessages) Size() int {
	return len(n.messages)
}

func (n *noticeMessages) Delete(clientID string) {
	n.lock.Lock()
	if msg, ok := n.messages[clientID]; ok {
		close(msg)
		delete(n.messages, clientID)
	}
	n.lock.Unlock()
}

func (n *noticeMessages) Clear() {
	n.lock.Lock()
	for key, msg := range n.messages {
		close(msg)
		delete(n.messages, key)
	}
	n.lock.Unlock()
}

func (n *noticeMessages) Add(clientID string) {
	n.lock.Lock()
	if _, ok := n.messages[clientID]; !ok {
		n.messages[clientID] = make(chan *Message)
	}
	n.lock.Unlock()
}

func (n *noticeMessages) Send(message *Message) error {
	n.lock.RLock()
	msg, ok := n.messages[message.ClientID]
	n.lock.RUnlock()
	if ok {
		msg <- message
		return nil
	}
	return ErrClientIDNotOnline
}

func (n *noticeMessages) Recv(clientID string) chan *Message {
	n.lock.RLock()
	message, ok := n.messages[clientID]
	n.lock.RUnlock()
	if ok {
		return message
	}
	return nil
}
