package goseaweedfs

import (
	"io"
	"mime"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// FilePart file wrapper with reader and some metadata
type FilePart struct {
	Reader     io.ReadCloser
	FileName   string
	FileSize   int64
	MimeType   string
	ModTime    int64 //in seconds
	Collection string

	// TTL Time to live.
	// 3m: 3 minutes
	// 4h: 4 hours
	// 5d: 5 days
	// 6w: 6 weeks
	// 7M: 7 months
	// 8y: 8 years
	TTL string

	Server string
	FileID string
}

// Close underlying openned file.
func (f *FilePart) Close() (err error) {
	err = f.Reader.Close()
	return
}

// NewFilePartFromReader new file part from file reader.
// fileName and fileSize must be known
func NewFilePartFromReader(reader io.ReadCloser, fileName string, fileSize int64) *FilePart {
	ret := FilePart{
		Reader:   reader,
		FileSize: fileSize,
		FileName: fileName,
	}

	ext := strings.ToLower(path.Ext(fileName))
	if ext != "" {
		ret.MimeType = mime.TypeByExtension(ext)
	}

	return &ret
}

// NewFilePart new file path from real file dir
func NewFilePart(fullPathFilename string) (*FilePart, error) {
	fh, openErr := os.Open(fullPathFilename)
	if openErr != nil {
		return nil, openErr
	}

	ret := FilePart{
		Reader:   fh,
		FileName: filepath.Base(fullPathFilename),
	}

	if fi, fiErr := fh.Stat(); fiErr == nil {
		ret.ModTime = fi.ModTime().UTC().Unix()
		ret.FileSize = fi.Size()
	} else {
		return nil, fiErr
	}

	ext := strings.ToLower(path.Ext(ret.FileName))
	if ext != "" {
		ret.MimeType = mime.TypeByExtension(ext)
	}

	return &ret, nil
}

// NewFileParts create many file part at once.
func NewFileParts(fullPathFilenames []string) (ret []*FilePart, err error) {
	ret = make([]*FilePart, 0, len(fullPathFilenames))
	for _, file := range fullPathFilenames {
		if fp, err := NewFilePart(file); err == nil {
			ret = append(ret, fp)
		} else {
			closeFileParts(ret)
			return nil, err
		}
	}
	return
}

func closeFileParts(fps []*FilePart) {
	for i := range fps {
		_ = fps[i].Close()
	}
}
