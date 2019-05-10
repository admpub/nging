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

package common

import (
	"encoding/gob"
	"errors"
)

func init() {
	gob.Register(&Success{})
}

var (
	//JSON
	ErrUserNotLoggedIn = errors.New(`User not logged in`)
	ErrUserNotFound    = errors.New(`User does not exist`)
	ErrUserNoPerm      = errors.New(`User has no permission`)
	ErrUserDisabled    = errors.New(`User has been disabled`)

	//Watcher
	ErrIgnoreConfigChange = errors.New(`Ingore file`)
)

var DefaultNopMessage Messager = &NopMessage{}

type NopMessage struct {
}

func (n *NopMessage) Error() string {
	return ``
}
func (n *NopMessage) Success() string {
	return ``
}

func (s *NopMessage) String() string {
	return ``
}

type Messager interface {
	Successor
	error
}

func IsMessage(err interface{}) bool {
	_, y := err.(Messager)
	return y
}

func Message(err interface{}) Messager {
	if v, y := err.(Messager); y {
		return v
	}
	return DefaultNopMessage
}

func NewOk(v string) Successor {
	return &Success{
		Value: v,
	}
}

type Success struct {
	Value string
}

func (s *Success) Success() string {
	return s.Value
}

func (s *Success) String() string {
	return s.Value
}

type Successor interface {
	Success() string
}

func IsError(err interface{}) bool {
	_, y := err.(error)
	return y
}

func IsOk(err interface{}) bool {
	_, y := err.(Successor)
	return y
}

func OkString(err interface{}) string {
	if v, y := err.(Successor); y {
		return v.Success()
	}
	return ``
}
