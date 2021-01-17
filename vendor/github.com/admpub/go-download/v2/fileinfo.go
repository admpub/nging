package download

import (
	"os"
	"time"
)

var _ os.FileInfo = (*fileInfo)(nil)

type fileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name returns the base name of the file
func (f *fileInfo) Name() string {
	return f.name
}

// Size returns file(s) length in bytes
func (f *fileInfo) Size() int64 {
	return f.size
}

// Mode returns the file mode bits
func (f *fileInfo) Mode() os.FileMode {
	return f.mode
}

// ModTime returns the file(s) modifications time
// NOTE: in the case of multiple files, from a split download,
// the Modtime will be that of the last partial file downloaded
func (f *fileInfo) ModTime() time.Time {
	return f.modTime
}

// IsDir returns if the file is a directory
func (f *fileInfo) IsDir() bool {
	return false
}

// Sys returns the underlying data source (can return nil)
func (f *fileInfo) Sys() interface{} {
	return nil
}
