package lib

import (
	"context"
	"mime"
	"os"
	"path"
	"path/filepath"

	"golang.org/x/net/webdav"
)

// NoSniffFileInfo wraps any generic FileInfo interface and bypasses mime type sniffing.
type NoSniffFileInfo struct {
	os.FileInfo
}

const MIMEDefault = "application/octet-stream"

func (w NoSniffFileInfo) ContentType(ctx context.Context) (contentType string, err error) {
	if mimeType := mime.TypeByExtension(path.Ext(w.FileInfo.Name())); len(mimeType) > 0 {
		// We can figure out the mime from the extension.
		return mimeType, nil
	}
	// We can't figure out the mime type without sniffing, call it an octet stream.
	return MIMEDefault, nil
}

type WebDavDir struct {
	webdav.Dir
	User    *User
	NoSniff bool
}

func (d WebDavDir) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	// Skip wrapping if NoSniff is off
	if !d.NoSniff {
		return d.Dir.Stat(ctx, name)
	}

	info, err := d.Dir.Stat(ctx, name)
	if err != nil {
		return nil, err
	}

	if name != `/` && d.User != nil {
		if !d.User.Allowed(name, true) {
			return nil, filepath.SkipDir
		}
	}

	return NoSniffFileInfo{info}, nil
}

func (d WebDavDir) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	// Skip wrapping if NoSniff is off
	if !d.NoSniff {
		return d.Dir.OpenFile(ctx, name, flag, perm)
	}

	file, err := d.Dir.OpenFile(ctx, name, flag, perm)
	if err != nil {
		return nil, err
	}

	return WebDavFile{File: file}, nil
}

type WebDavFile struct {
	webdav.File
}

func (f WebDavFile) Stat() (os.FileInfo, error) {
	info, err := f.File.Stat()
	if err != nil {
		return nil, err
	}

	return NoSniffFileInfo{info}, nil
}

func (f WebDavFile) Readdir(count int) (fis []os.FileInfo, err error) {
	fis, err = f.File.Readdir(count)
	if err != nil {
		return nil, err
	}

	for i, fi := range fis {
		fis[i] = NoSniffFileInfo{fi}
	}
	return fis, nil
}
