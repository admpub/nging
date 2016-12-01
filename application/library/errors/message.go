package errors

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
