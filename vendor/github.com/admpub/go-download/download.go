package download

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	defaultGoroutines = 10
	defaultDir        = "go-download"
)

var (
	_           io.Reader = (*File)(nil)
	fileMode              = os.FileMode(0770)
	defaultTime           = time.Time{}
)

// Options contains any specific configuration values
// for downloading/opening a file
type Options struct {
	Concurrency ConcurrencyFn
	Proxy       ProxyFn
	Client      ClientFn
	Request     RequestFn
}

// RequestFn allows for additional information, such as http headers, to the http request
//
// Do not alter the "Range" http headers or the download can become corrupt
type RequestFn func(r *http.Request)

// ClientFn allows for a custom http.Client to be used for the http request
type ClientFn func() http.Client

// ConcurrencyFn is the function used to determine the level of concurrency aka the
// number of goroutines to use. Default concurrency level is 10
//
// if returned value is < 1 then the default value will be used
type ConcurrencyFn func(size int64) int

// ProxyFn is the function used to pass the download io.Reader for proxying.
// eg. displaying a progress bar of the download.
type ProxyFn func(name string, download int, size int64, r io.Reader) io.Reader

// File represents an open file descriptor to a downloaded file(s)
type File struct {
	url      string
	dir      string
	baseName string
	size     int64
	modTime  time.Time
	options  *Options
	readers  []io.ReadCloser
	io.Reader
}

type partialResult struct {
	idx int
	r   io.ReadCloser
	err error
}

// Open downloads and opens the file(s) downloaded by the given url
func Open(url string, options *Options) (*File, error) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return OpenContext(ctx, url, options)
}

// OpenContext downloads and opens the file(s) downloaded by the given url and is cancellable using the provided context.
// The context provided must be non-nil
func OpenContext(ctx context.Context, url string, options *Options) (*File, error) {

	if ctx == nil {
		panic("nil context")
	}

	f := &File{
		url:      url,
		baseName: filepath.Base(url),
		options:  options,
	}

	req, err := http.NewRequest(http.MethodHead, f.url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if f.options != nil && f.options.Request != nil {
		f.options.Request(req)
	}

	var client http.Client
	if f.options != nil && f.options.Client != nil {
		client = f.options.Client()
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		// not all services support HEAD requests
		// so if this fails just move along to the
		// GET portion, with a warning
		log.Printf("notice: unexpected HEAD response code '%d', proceeding with download.\n", resp.StatusCode)
		err = f.download(ctx)
	} else {
		f.size = resp.ContentLength

		if t := resp.Header.Get("Accept-Ranges"); t == "bytes" {
			err = f.downloadRangeBytes(ctx)
		} else {
			err = f.download(ctx)
		}
	}

	if err != nil {
		f.closeFileHandles()
		return nil, err
	}

	return f, nil
}

func (f *File) download(ctx context.Context) error {

	req, err := http.NewRequest(http.MethodGet, f.url, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	if f.options != nil && f.options.Request != nil {
		f.options.Request(req)
	}

	var client http.Client

	if f.options != nil && f.options.Client != nil {
		client = f.options.Client()
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &InvalidResponseCode{got: resp.StatusCode, expected: http.StatusOK}
	}

	f.dir, err = ioutil.TempDir("", defaultDir)
	if err != nil {
		return err
	}

	fh, err := ioutil.TempFile(f.dir, "")
	if err != nil {
		return err
	}

	f.readers = make([]io.ReadCloser, 1)
	f.readers[0] = fh

	var read io.Reader = resp.Body

	if f.options != nil && f.options.Proxy != nil {
		read = f.options.Proxy(f.baseName, 0, f.size, read)
	}

	_, err = io.Copy(fh, read)
	if err != nil {
		return err
	}

	fh.Seek(0, 0)

	f.Reader = fh
	f.modTime = time.Now()

	return nil
}

func (f *File) downloadRangeBytes(ctx context.Context) (err error) {

	if f.size <= 0 {
		return fmt.Errorf("Invalid content length '%d'", f.size)
	}

	var resume bool

	f.dir = filepath.Join(os.TempDir(), defaultDir+f.generateHash())

	if _, err = os.Stat(f.dir); os.IsNotExist(err) {
		err = os.Mkdir(f.dir, fileMode) // only owner and group have RWX access
		if err != nil {
			return
		}
	} else {
		resume = true
	}

	var goroutines int

	if f.options == nil || f.options.Concurrency == nil {
		goroutines = defaultConcurrencyFn(f.size)
	} else {
		goroutines = f.options.Concurrency(f.size)
		if goroutines < 1 {
			goroutines = defaultConcurrencyFn(f.size)
		}
	}

	chunkSize := f.size / int64(goroutines)
	remainer := f.size % chunkSize
	var pos int64

	chunkSize--

	f.readers = make([]io.ReadCloser, goroutines, goroutines)

	ch := make(chan partialResult)

	var i int

	for ; i < goroutines; i++ {

		if i == goroutines-1 {
			chunkSize += remainer // add remainer to last download
		}

		go f.downloadPartial(ctx, resume, i, pos, pos+chunkSize, ch)

		pos += chunkSize + 1
	}

	for i = 0; i < goroutines; i++ {

		select {
		case <-ctx.Done():

			if ctx.Err() == context.Canceled {
				err = &Canceled{url: f.url}
			} else {
				// context.DeadlineExceeded
				err = &DeadlineExceeded{url: f.url}
			}

			//drain remaining
			for ; i < goroutines; i++ {
				res := <-ch
				f.readers[res.idx] = res.r
				break
			}

		case res := <-ch:

			f.readers[res.idx] = res.r

			if err != nil {
				continue
			}

			if res.err != nil {
				err = res.err
			}
		}
	}

	close(ch)

	readers := make([]io.Reader, len(f.readers))
	for i = 0; i < len(f.readers); i++ {
		readers[i] = f.readers[i]
	}

	f.Reader = io.MultiReader(readers...)
	f.modTime = time.Now()
	return
}

func (f *File) downloadPartial(ctx context.Context, resumeable bool, idx int, start, end int64, ch chan<- partialResult) {

	var err error
	var fh *os.File

	defer func() {
		ch <- partialResult{idx: idx, err: err, r: fh}
	}()

	fPath := filepath.Join(f.dir, strconv.Itoa(idx))

	if resumeable {
		var fi os.FileInfo

		fi, err = os.Stat(fPath)
		if os.IsNotExist(err) {
			fh, err = os.Create(fPath)
		} else {

			// file exists...musts check if partial
			if fi.Size() < (end-start)+1 {

				// lets append/download only the bytes necessary
				start += fi.Size()
				fh, err = os.OpenFile(fPath, os.O_RDWR|os.O_APPEND, fileMode)
			} else {
				fh, err = os.Open(fPath)
				return // if error or not still leaving
			}
		}
	} else {
		fh, err = os.Create(fPath)
	}

	if err != nil {
		return
	}

	var client http.Client
	if f.options != nil && f.options.Client != nil {
		client = f.options.Client()
	}

	var req *http.Request
	if req, err = http.NewRequest(http.MethodGet, f.url, nil); err != nil {
		return
	}
	req = req.WithContext(ctx)
	req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", start, end))

	if f.options != nil && f.options.Request != nil {
		f.options.Request(req)
	}

	var resp *http.Response

	if resp, err = client.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusPartialContent {
		err = &InvalidResponseCode{got: resp.StatusCode, expected: http.StatusPartialContent}
		return
	}

	// check for timeout or cancellation before heaviest operation
	select {
	case <-ctx.Done():
		return
	default:
	}

	var read io.Reader = resp.Body

	if f.options != nil && f.options.Proxy != nil {
		read = f.options.Proxy(f.baseName, idx, (end-start)+1, read)
	}

	if _, err = io.Copy(fh, read); err != nil {
		return
	}

	fh.Seek(0, 0)
}

// Stat returns the FileInfo structure describing file(s). If there is an error, it will be of type *PathError.
func (f *File) Stat() (os.FileInfo, error) {

	if f.modTime.IsZero() {
		return nil, &os.PathError{Op: "stat", Path: filepath.Base(f.url), Err: errors.New("bad file descriptor")}
	}

	return &fileInfo{
		name:    filepath.Base(f.url),
		size:    f.size,
		mode:    fileMode,
		modTime: f.modTime,
	}, nil
}

// Close closes the File(s), rendering it unusable for I/O. It returns an error, if any.
func (f *File) Close() error {

	f.closeFileHandles()
	f.modTime = defaultTime

	return os.RemoveAll(f.dir)
}

func (f *File) closeFileHandles() {
	for i := 0; i < len(f.readers); i++ {
		if f.readers[i] != nil { // possible if cancelled or error occured
			f.readers[i].Close()
		}
	}
}

func (f *File) generateHash() string {

	// Open to a better way, but should not collide
	h := sha1.New()
	io.WriteString(h, f.url)

	return fmt.Sprintf("%x", h.Sum(nil))
}

// chunks up downloads into 2MB chunks, when Accept-Ranges supported
func defaultConcurrencyFn(length int64) int {
	return defaultGoroutines
}
