package upload

import (
	"io"
	"strings"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/matchers"
	"github.com/h2non/filetype/types"
)

func ReadHeadBytes(r io.Reader) ([]byte, error) {
	// We only have to pass the file header = first 261 bytes
	head := make([]byte, 261)
	_, err := r.Read(head)
	return head, err
}

func IsImage(b []byte) bool {
	return filetype.IsImage(b)
}

func IsVideo(b []byte) bool {
	return filetype.IsVideo(b)
}

func IsAudio(b []byte) bool {
	return filetype.IsAudio(b)
}

func IsFont(b []byte) bool {
	return filetype.IsFont(b)
}

func IsArchive(b []byte) bool {
	return filetype.IsArchive(b)
}

func IsDocument(b []byte) bool {
	return filetype.IsDocument(b)
}

func IsApplication(b []byte) bool {
	return filetype.IsApplication(b)
}

func IsType(b []byte, expected types.Type) bool {
	kind, _ := filetype.Match(b)
	return kind == expected
}

func IsTypeString(b []byte, expected string) bool {
	switch expected {
	case `image`:
		return IsImage(b)
	case `video`:
		return IsVideo(b)
	case `audio`:
		return IsAudio(b)
	case `archive`:
		return IsArchive(b)
	case `document`, `office`:
		return IsDocument(b)
	case `doc`:
		return IsType(b, matchers.TypeDoc) || IsType(b, matchers.TypeDocx)
	case `ppt`:
		return IsType(b, matchers.TypePpt) || IsType(b, matchers.TypePptx)
	case `xls`:
		return IsType(b, matchers.TypeXls) || IsType(b, matchers.TypeXlsx)
	case `file`:
		return IsApplication(b)
	case `font`:
		return IsFont(b)
	case `pdf`:
		return IsType(b, matchers.TypePdf)
	case `photoshop`:
		return IsType(b, matchers.TypePsd)
	default:
		return false
	}
}

func IsSupported(extension string) bool {
	extension = strings.TrimPrefix(extension, `.`)
	return filetype.IsSupported(extension)
}

func IsMIMESupported(mime string) bool {
	return filetype.IsMIMESupported(mime)
}
