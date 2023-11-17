package sevenzip

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/bodgit/plumbing"
	"github.com/bodgit/sevenzip/internal/pool"
	"github.com/bodgit/sevenzip/internal/util"
	"github.com/hashicorp/go-multierror"
	"go4.org/readerutil"
)

var (
	errFormat   = errors.New("sevenzip: not a valid 7-zip file")
	errChecksum = errors.New("sevenzip: checksum error")
	errTooMuch  = errors.New("sevenzip: too much data")
)

//nolint:gochecknoglobals
var newPool pool.Constructor = pool.NewPool

// A Reader serves content from a 7-Zip archive.
type Reader struct {
	r     io.ReaderAt
	start int64
	end   int64
	si    *streamsInfo
	p     string
	File  []*File
	pool  []pool.Pooler

	fileListOnce sync.Once
	fileList     []fileListEntry
}

// A ReadCloser is a Reader that must be closed when no longer needed.
type ReadCloser struct {
	f []*os.File
	Reader
}

// A File is a single file in a 7-Zip archive. The file information is in the
// embedded FileHeader. The file content can be accessed by calling Open.
type File struct {
	FileHeader
	zip    *Reader
	folder int
	offset int64
}

type fileReader struct {
	rc util.SizeReadSeekCloser
	f  *File
	n  int64
}

func (fr *fileReader) Stat() (fs.FileInfo, error) {
	return headerFileInfo{&fr.f.FileHeader}, nil
}

func (fr *fileReader) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	if fr.n <= 0 {
		return 0, io.EOF
	}

	if int64(len(p)) > fr.n {
		p = p[0:fr.n]
	}

	n, err = fr.rc.Read(p)
	fr.n -= int64(n)

	return
}

func (fr *fileReader) Close() error {
	if fr.rc == nil {
		return nil
	}

	offset, err := fr.rc.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}

	if offset == fr.rc.Size() { // EOF reached
		if err := fr.rc.Close(); err != nil {
			return err
		}
	} else {
		f := fr.f
		if _, err := f.zip.pool[f.folder].Put(offset, fr.rc); err != nil {
			return err
		}
	}

	fr.rc = nil

	return nil
}

// Open returns an io.ReadCloser that provides access to the File's contents.
// Multiple files may be read concurrently.
func (f *File) Open() (io.ReadCloser, error) {
	if f.FileHeader.isEmptyStream || f.FileHeader.isEmptyFile {
		// Return empty reader for directory or empty file
		return io.NopCloser(bytes.NewReader(nil)), nil
	}

	var err error

	rc, _ := f.zip.pool[f.folder].Get(f.offset)
	if rc == nil {
		rc, _, err = f.zip.folderReader(f.zip.si, f.folder)
		if err != nil {
			return nil, err
		}
	}

	if _, err = rc.Seek(f.offset, io.SeekStart); err != nil {
		return nil, err
	}

	return &fileReader{
		rc: rc,
		f:  f,
		n:  int64(f.UncompressedSize),
	}, nil
}

// OpenReaderWithPassword will open the 7-zip file specified by name using
// password as the basis of the decryption key and return a ReadCloser. If
// name has a ".001" suffix it is assumed there are multiple volumes and each
// sequential volume will be opened.
//
//nolint:cyclop,funlen
func OpenReaderWithPassword(name, password string) (*ReadCloser, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	info, err := f.Stat()
	if err != nil {
		err = multierror.Append(err, f.Close())

		return nil, err
	}

	var reader io.ReaderAt = f

	size := info.Size()
	files := []*os.File{f}

	if ext := filepath.Ext(name); ext == ".001" {
		sr := []readerutil.SizeReaderAt{io.NewSectionReader(f, 0, size)}

		for i := 2; true; i++ {
			f, err := os.Open(fmt.Sprintf("%s.%03d", strings.TrimSuffix(name, ext), i))
			if err != nil {
				if errors.Is(err, fs.ErrNotExist) {
					break
				}

				for _, file := range files {
					err = multierror.Append(err, file.Close())
				}

				return nil, err
			}

			files = append(files, f)

			info, err = f.Stat()
			if err != nil {
				for _, file := range files {
					err = multierror.Append(err, file.Close())
				}

				return nil, err
			}

			sr = append(sr, io.NewSectionReader(f, 0, info.Size()))
		}

		mr := readerutil.NewMultiReaderAt(sr...)
		reader, size = mr, mr.Size()
	}

	r := new(ReadCloser)
	r.p = password

	if err := r.init(reader, size); err != nil {
		for _, file := range files {
			err = multierror.Append(err, file.Close())
		}

		return nil, err
	}

	r.f = files

	return r, nil
}

// OpenReader will open the 7-zip file specified by name and return a
// ReadCloser. If name has a ".001" suffix it is assumed there are multiple
// volumes and each sequential volume will be opened.
func OpenReader(name string) (*ReadCloser, error) {
	return OpenReaderWithPassword(name, "")
}

// NewReaderWithPassword returns a new Reader reading from r using password as
// the basis of the decryption key, which is assumed to have the given size in
// bytes.
func NewReaderWithPassword(r io.ReaderAt, size int64, password string) (*Reader, error) {
	if size < 0 {
		return nil, errors.New("sevenzip: size cannot be negative")
	}

	zr := new(Reader)
	zr.p = password

	if err := zr.init(r, size); err != nil {
		return nil, err
	}

	return zr, nil
}

// NewReader returns a new Reader reading from r, which is assumed to have the
// given size in bytes.
func NewReader(r io.ReaderAt, size int64) (*Reader, error) {
	return NewReaderWithPassword(r, size, "")
}

func (z *Reader) folderReader(si *streamsInfo, f int) (*folderReadCloser, uint32, error) {
	// Create a SectionReader covering all of the streams data
	return si.FolderReader(io.NewSectionReader(z.r, z.start, z.end), f, z.p)
}

//nolint:cyclop,funlen,gocognit
func (z *Reader) init(r io.ReaderAt, size int64) error {
	h := crc32.NewIEEE()
	tra := plumbing.TeeReaderAt(r, h)
	sr := io.NewSectionReader(tra, 0, size) // Will only read first 32 bytes

	var sh signatureHeader
	if err := binary.Read(sr, binary.LittleEndian, &sh); err != nil {
		return err
	}

	signature := []byte{'7', 'z', 0xbc, 0xaf, 0x27, 0x1c}
	if !bytes.Equal(sh.Signature[:], signature) {
		return errFormat
	}

	z.r = r

	h.Reset()

	var (
		err   error
		start startHeader
	)

	if err = binary.Read(sr, binary.LittleEndian, &start); err != nil {
		return err
	}

	// CRC of the start header should match
	if !util.CRC32Equal(h.Sum(nil), sh.CRC) {
		return errChecksum
	}

	// Work out where we are in the file (32, avoiding magic numbers)
	if z.start, err = sr.Seek(0, io.SeekCurrent); err != nil {
		return err
	}

	// Seek over the streams
	if z.end, err = sr.Seek(int64(start.Offset), io.SeekCurrent); err != nil {
		return err
	}

	h.Reset()

	// Bound bufio.Reader otherwise it can read trailing garbage which screws up the CRC check
	br := bufio.NewReader(io.NewSectionReader(tra, z.end, int64(start.Size)))

	id, err := br.ReadByte()
	if err != nil {
		return err
	}

	var (
		header      *header
		streamsInfo *streamsInfo
	)

	switch id {
	case idHeader:
		if header, err = readHeader(br); err != nil {
			return err
		}
	case idEncodedHeader:
		if streamsInfo, err = readStreamsInfo(br); err != nil {
			return err
		}
	default:
		return errUnexpectedID
	}

	// If there's more data to read, we've not parsed this correctly. This
	// won't break with trailing data as the bufio.Reader was bounded
	if n, _ := io.CopyN(io.Discard, br, 1); n != 0 {
		return errTooMuch
	}

	// CRC should match the one from the start header
	if !util.CRC32Equal(h.Sum(nil), start.CRC) {
		return errChecksum
	}

	// If the header was encoded we should have sufficient information now
	// to decode it
	if streamsInfo != nil {
		if streamsInfo.Folders() != 1 {
			return errors.New("sevenzip: expected only one folder in header stream")
		}

		fr, crc, err := z.folderReader(streamsInfo, 0)
		if err != nil {
			return err
		}
		defer fr.Close()

		if header, err = readEncodedHeader(util.ByteReadCloser(fr)); err != nil {
			return err
		}

		if crc != 0 && !util.CRC32Equal(fr.Checksum(), crc) {
			return errChecksum
		}
	}

	z.si = header.streamsInfo

	z.pool = make([]pool.Pooler, z.si.Folders())
	for i := range z.pool {
		if z.pool[i], err = newPool(); err != nil {
			return err
		}
	}

	// spew.Dump(header)

	folder, offset := 0, int64(0)
	z.File = make([]*File, 0, len(header.filesInfo.file))
	j := 0

	for _, fh := range header.filesInfo.file {
		f := new(File)
		f.zip = z
		f.FileHeader = fh

		if f.FileHeader.FileInfo().IsDir() && !strings.HasSuffix(f.FileHeader.Name, "/") {
			f.FileHeader.Name += "/"
		}

		if !fh.isEmptyStream && !fh.isEmptyFile {
			f.folder, _ = header.streamsInfo.FileFolderAndSize(j)

			if f.folder != folder {
				offset = 0
			}

			f.offset = offset
			offset += int64(f.UncompressedSize)
			folder = f.folder
			j++
		}

		z.File = append(z.File, f)
	}

	return nil
}

// Close closes the 7-zip file or volumes, rendering them unusable for I/O.
func (rc *ReadCloser) Close() error {
	var err *multierror.Error
	for _, f := range rc.f {
		err = multierror.Append(err, f.Close())
	}

	return err.ErrorOrNil()
}

type fileListEntry struct {
	name  string
	file  *File
	isDir bool
	isDup bool
}

type fileInfoDirEntry interface {
	fs.FileInfo
	fs.DirEntry
}

func (e *fileListEntry) stat() (fileInfoDirEntry, error) {
	if e.isDup {
		return nil, errors.New(e.name + ": duplicate entries in 7-zip file")
	}

	if !e.isDir {
		return headerFileInfo{&e.file.FileHeader}, nil
	}

	return e, nil
}

func (e *fileListEntry) Name() string {
	_, elem := split(e.name)

	return elem
}

func (e *fileListEntry) Size() int64       { return 0 }
func (e *fileListEntry) Mode() fs.FileMode { return fs.ModeDir | 0o555 }
func (e *fileListEntry) Type() fs.FileMode { return fs.ModeDir }
func (e *fileListEntry) IsDir() bool       { return true }
func (e *fileListEntry) Sys() interface{}  { return nil }

func (e *fileListEntry) ModTime() time.Time {
	if e.file == nil {
		return time.Time{}
	}

	return e.file.FileHeader.Modified.UTC()
}

func (e *fileListEntry) Info() (fs.FileInfo, error) { return e, nil }

func toValidName(name string) string {
	name = strings.ReplaceAll(name, `\`, `/`)

	p := strings.TrimPrefix(path.Clean(name), "/")

	for strings.HasPrefix(p, "../") {
		p = p[len("../"):]
	}

	return p
}

//nolint:cyclop,gocognit
func (z *Reader) initFileList() {
	z.fileListOnce.Do(func() {
		files := make(map[string]int)
		knownDirs := make(map[string]int)

		dirs := make(map[string]struct{})

		for _, file := range z.File {
			isDir := len(file.Name) > 0 && file.Name[len(file.Name)-1] == '/'

			name := toValidName(file.Name)
			if name == "" {
				continue
			}

			if idx, ok := files[name]; ok {
				z.fileList[idx].isDup = true

				continue
			}

			if idx, ok := knownDirs[name]; ok {
				z.fileList[idx].isDup = true

				continue
			}

			for dir := path.Dir(name); dir != "."; dir = path.Dir(dir) {
				dirs[dir] = struct{}{}
			}

			idx := len(z.fileList)
			entry := fileListEntry{
				name:  name,
				file:  file,
				isDir: isDir,
			}
			z.fileList = append(z.fileList, entry)
			if isDir {
				knownDirs[name] = idx
			} else {
				files[name] = idx
			}
		}
		for dir := range dirs {
			if _, ok := knownDirs[dir]; !ok {
				if idx, ok := files[dir]; ok {
					z.fileList[idx].isDup = true
				} else {
					entry := fileListEntry{
						name:  dir,
						file:  nil,
						isDir: true,
					}
					z.fileList = append(z.fileList, entry)
				}
			}
		}

		sort.Slice(z.fileList, func(i, j int) bool { return fileEntryLess(z.fileList[i].name, z.fileList[j].name) })
	})
}

func fileEntryLess(x, y string) bool {
	xdir, xelem := split(x)
	ydir, yelem := split(y)

	return xdir < ydir || xdir == ydir && xelem < yelem
}

// Open opens the named file in the 7-zip archive, using the semantics of
// fs.FS.Open: paths are always slash separated, with no leading / or ../
// elements.
func (z *Reader) Open(name string) (fs.File, error) {
	z.initFileList()

	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrInvalid}
	}

	e := z.openLookup(name)
	if e == nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}

	if e.isDir {
		return &openDir{e, z.openReadDir(name), 0}, nil
	}

	rc, err := e.file.Open()
	if err != nil {
		return nil, err
	}

	return rc.(fs.File), nil //nolint:forcetypeassert
}

func split(name string) (dir, elem string) {
	if len(name) > 0 && name[len(name)-1] == '/' {
		name = name[:len(name)-1]
	}

	i := len(name) - 1
	for i >= 0 && name[i] != '/' {
		i--
	}

	if i < 0 {
		return ".", name
	}

	return name[:i], name[i+1:]
}

//nolint:gochecknoglobals
var dotFile = &fileListEntry{name: "./", isDir: true}

func (z *Reader) openLookup(name string) *fileListEntry {
	if name == "." {
		return dotFile
	}

	dir, elem := split(name)

	files := z.fileList
	i := sort.Search(len(files), func(i int) bool {
		idir, ielem := split(files[i].name)

		return idir > dir || idir == dir && ielem >= elem
	})

	if i < len(files) {
		fname := files[i].name
		if fname == name || len(fname) == len(name)+1 && fname[len(name)] == '/' && fname[:len(name)] == name {
			return &files[i]
		}
	}

	return nil
}

func (z *Reader) openReadDir(dir string) []fileListEntry {
	files := z.fileList

	i := sort.Search(len(files), func(i int) bool {
		idir, _ := split(files[i].name)

		return idir >= dir
	})

	j := sort.Search(len(files), func(j int) bool {
		jdir, _ := split(files[j].name)

		return jdir > dir
	})

	return files[i:j]
}

type openDir struct {
	e      *fileListEntry
	files  []fileListEntry
	offset int
}

func (d *openDir) Close() error               { return nil }
func (d *openDir) Stat() (fs.FileInfo, error) { return d.e.stat() }

func (d *openDir) Read([]byte) (int, error) {
	return 0, &fs.PathError{Op: "read", Path: d.e.name, Err: errors.New("is a directory")}
}

func (d *openDir) ReadDir(count int) ([]fs.DirEntry, error) {
	n := len(d.files) - d.offset
	if count > 0 && n > count {
		n = count
	}

	if n == 0 {
		if count <= 0 {
			return nil, nil
		}

		return nil, io.EOF
	}

	list := make([]fs.DirEntry, n)
	for i := range list {
		s, err := d.files[d.offset+i].stat()
		if err != nil {
			return nil, err
		}

		list[i] = s
	}

	d.offset += n

	return list, nil
}
