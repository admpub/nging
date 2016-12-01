package errors

func IsError(err interface{}) bool {
	_, y := err.(error)
	return y
}
