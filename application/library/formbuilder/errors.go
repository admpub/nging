package formbuilder

func newPostError(err error) *postError {
	return &postError{error: err}
}

type postError struct {
	error
}

func (e *postError) Unwrap() error {
	return e.error
}

func ErrPostFailed(err error) bool {
	_, ok := err.(*postError)
	return ok
}
