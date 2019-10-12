package pkg

import "errors"

var (
	ErrInvalidMediaPlaylist  = errors.New("Not a valid media playlist")
	ErrInvalidMasterPlaylist = errors.New("Not a valid master playlist")
)
