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
