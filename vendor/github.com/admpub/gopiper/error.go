package gopiper

import "errors"

var (
	ErrJsonparseNeedSubItem    = errors.New("Pipe type jsonparse need one subItem")
	ErrArrayNeedSubItem        = errors.New("Pipe type array need one subItem")
	ErrNotSupportPipeType      = errors.New("Not support pipe type")
	ErrUnknowHTMLAttr          = errors.New("Unknow html attr")
	ErrUnsupportText2boolType  = errors.New("Unsupport text2bool type")
	ErrUnsupportText2floatType = errors.New("Unsupport text2float type")
	ErrUnsupportText2intType   = errors.New("Unsupport text2int type")
	ErrTrimNilParams           = errors.New("Filter trim nil params")
	ErrSplitNilParams          = errors.New("Filter split nil params")
	ErrJoinNilParams           = errors.New("Filter join nil params")
	ErrFetcherNotRegistered    = errors.New("Fetcher not registered")
	ErrStorerNotRegistered     = errors.New("Storer not registered")
	ErrInvalidContent          = errors.New("Invalid content")
)
