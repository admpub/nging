package formbuilder

func newPostError(err error) *postError {
	return &postError{err: err}
}

type postError struct {
	err error
}

func (e *postError) Unwrap() error {
	return e.err
}

func ErrPostFailed(err error) bool {
	_, ok := err.(*postError)
	return ok
}
