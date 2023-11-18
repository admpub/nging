package notice

type NProgressor interface {
	Send(message interface{}, statusCode int) error
	Success(message interface{}) error
	Failure(message interface{}) error
	Add(n int64) NProgressor
	Done(n int64) NProgressor
	AutoComplete(on bool) NProgressor
	Complete() NProgressor
}
