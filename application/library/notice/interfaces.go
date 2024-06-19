package notice

import "io"

type NProgressor interface {
	Notifier
	Progressor
}

type Notifier interface {
	Send(message interface{}, statusCode int) error
	Success(message interface{}) error
	Failure(message interface{}) error
}

type Progressor interface {
	Add(n int64) NProgressor
	Done(n int64) NProgressor
	AutoComplete(on bool) NProgressor
	Complete() NProgressor
	Reset()
	ProxyReader(r io.Reader) io.ReadCloser
	ProxyWriter(w io.Writer) io.WriteCloser
	Callback(total int64, exec func(callback func(strLen int)) error) error
}

type IOnlineUser interface {
	GetUser() string
	HasMessageType(messageTypes ...string) bool
	Send(message *Message, openDebug ...bool) error
	Recv(clientID string) <-chan *Message
	ClearMessage()
	ClearMessageType(types ...string)
	OpenMessageType(types ...string)
	CountType() int
	CountClient() int
	CloseClient(clientID string)
	OpenClient(clientID string)
}

type IOnlineUsers interface {
	GetOk(user string, noLock ...bool) (IOnlineUser, bool)
	OnlineStatus(users ...string) map[string]bool
	Set(user string, oUser IOnlineUser)
	Delete(user string)
	Clear()
}

type NoticeMessager interface {
	Size() int
	Delete(clientID string)
	Clear()
	Add(clientID string)
	Send(message *Message) error
	Recv(clientID string) <-chan *Message
}

type NoticeTyper interface {
	Has(types ...string) bool
	Size() int
	Clear(types ...string)
	Open(types ...string)
}
