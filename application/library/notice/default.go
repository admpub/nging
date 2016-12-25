/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package notice

var DefaultUserNotices = NewUserNotices()

func Send(user string, message *Message) {
	DefaultUserNotices.Send(user, message)
}

func Recv(user string) <-chan *Message {
	return DefaultUserNotices.Recv(user)
}

func RecvJSON(user string) []byte {
	return DefaultUserNotices.RecvJSON(user)
}

func RecvXML(user string) []byte {
	return DefaultUserNotices.RecvXML(user)
}

func CloseClient(user string) bool {
	return DefaultUserNotices.CloseClient(user)
}

func OpenClient(user string) {
	DefaultUserNotices.OpenClient(user)
}

func CloseMessage(user string, types ...string) {
	DefaultUserNotices.CloseMessage(user, types...)
}

func OpenMessage(user string, types ...string) {
	DefaultUserNotices.OpenMessage(user, types...)
}

func Clear() {
	DefaultUserNotices.Clear()
}
