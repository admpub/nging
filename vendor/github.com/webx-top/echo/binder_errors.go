package echo

import "errors"

var (
	ErrBreak              = errors.New("[BREAK]")
	ErrContinue           = errors.New("[CONTINUE]")
	ErrExit               = errors.New("[EXIT]")
	ErrReturn             = errors.New("[RETURN]")
	ErrSliceIndexTooLarge = errors.New("the slice index value of the form field is too large")
)
