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
