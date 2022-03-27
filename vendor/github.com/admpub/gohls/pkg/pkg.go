package pkg

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/golang/groupcache/lru"
	"github.com/grafov/m3u8"
	"github.com/webx-top/com"
)

var (
	UserAgent     = "Mozilla/5.0 (X11; Linux x86_64; rv:38.0) Gecko/38.0 Firefox/38.0"
	IVPlaceholder = []byte{0, 0, 0, 0, 0, 0, 0, 0}
)

func DoRequest(c *http.Client, req *http.Request, maxRetries int) (*http.Response, error) {
	req.Header.Set("User-Agent", UserAgent)
	//req.Header.Set("Connection", "Keep-Alive") //http2不支持Keep-Alive

	var i int

RETRY:
	resp, err := c.Do(req)
	if err != nil {
		i++
		if i <= maxRetries {
			wait := time.Second * time.Duration(i)
			log.Printf("[%v] %v => %v: wait %v and try again (%d/%d)\n", req.Method, req.URL.String(), err, wait, i, maxRetries)
			time.Sleep(wait)
			goto RETRY
		}
		return nil, err
	}

	// Maybe in the future it will force connection to stay opened for "Connection: close"
	resp.Close = false
	resp.Request.Close = false

	return resp, err
}

type Download struct {
	URI           string
	SeqNo         uint64
	ExtXKey       *m3u8.Key
	totalDuration time.Duration
}

func DecryptData(data []byte, v *Download, aes128Keys *map[string][]byte) error {
	var (
		iv          *bytes.Buffer
		keyData     []byte
		cipherBlock cipher.Block
	)

	if v.ExtXKey != nil && (v.ExtXKey.Method == "AES-128" || v.ExtXKey.Method == "aes-128") {

		keyData = (*aes128Keys)[v.ExtXKey.URI]

		if keyData == nil {
			resp, err := Request().Get(v.ExtXKey.URI)
			if err != nil {
				return err
			}
			(*aes128Keys)[v.ExtXKey.URI] = resp.Body()
		}

		if len(v.ExtXKey.IV) == 0 {
			iv = bytes.NewBuffer(IVPlaceholder)
			binary.Write(iv, binary.BigEndian, v.SeqNo)
		} else {
			iv = bytes.NewBufferString(v.ExtXKey.IV)
		}

		cipherBlock, _ = aes.NewCipher((*aes128Keys)[v.ExtXKey.URI])
		cipher.NewCBCDecrypter(cipherBlock, iv.Bytes()).CryptBlocks(data, data)
	}
	return nil
}

func DownloadSegment(ctx context.Context, cfg *Config, dlc chan *Download) error {
	prog := cfg.progress
	var (
		out *os.File
		err error
	)
	if cfg.ForceRerownload || prog.FinishedNum <= 0 {
		out, err = os.Create(cfg.OutputFile)
	} else {
		out, err = os.OpenFile(cfg.OutputFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	}
	if err != nil {
		return fmt.Errorf(`%v: %w`, cfg.OutputFile, err)
	}
	defer out.Close()

	var (
		data       []byte
		aes128Keys = &map[string][]byte{}
	)

	defer func() {
		if e := recover(); e != nil {
			log.Println(e)
		}
		os.WriteFile(cfg.OutputFile+`._prog_`, prog.JSONBytes(), 0666)
	}()
	var i int
	var exec = func(v *Download) error {
		i++
		if i <= prog.FinishedNum {
			return nil
		}
		startTime := time.Now()
		resp, err := Request().Get(v.URI)
		if err != nil {
			return err
		}
		if !resp.IsSuccess() {
			err = fmt.Errorf("Received HTTP %v for %v", resp.StatusCode(), v.URI)
			return err
		}
		data = resp.Body()
		err = DecryptData(data, v, aes128Keys)
		if err != nil {
			return err
		}
		var size int
		size, err = out.Write(data)
		// _, err = io.Copy(out, resp.Body)
		if err != nil {
			return err
		}

		processSize := len(fmt.Sprint(prog.TotalNum))
		prog.FinishedSize += int64(size)
		prog.SpeedInSecond = float64(size) / time.Since(startTime).Seconds()
		speedDesc := com.FormatBytes(prog.SpeedInSecond)

		prog.FinishedNum++
		log.Printf("[%0*d/%d][%v/s] Downloaded %v\n", processSize, prog.FinishedNum, prog.TotalNum, speedDesc, v.URI)
		if cfg.Duration != 0 {
			log.Printf("Recorded %v of %v\n", v.totalDuration, cfg.Duration)
		} else {
			log.Printf("Recorded %v\n", v.totalDuration)
		}
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return ErrContextCancelled
		case v, ok := <-dlc:
			if !ok {
				return err
			}
			if err := exec(v); err != nil {
				log.Println(err)
				return err
			}
		}
	}
}

func IsFullURL(url string) bool {
	if len(url) < 8 {
		return false
	}
	switch strings.ToLower(url[0:7]) {
	case `https:/`, `http://`:
		return true
	default:
		return false
	}
}

func ParseURI(root *url.URL, uri string) (string, error) {
	msURI, err := url.QueryUnescape(uri)
	if err != nil {
		return msURI, err
	}
	if !IsFullURL(msURI) {
		var msURL *url.URL
		msURL, err = root.Parse(msURI)
		if err != nil {
			return msURI, err
		}
		msURI, err = url.QueryUnescape(msURL.String())
	}
	return msURI, err
}

type Context struct {
	playlistURL *url.URL
	startTime   time.Time
	recDuration time.Duration
	cache       *lru.Cache
}

func (c *Context) Close() error {
	c.cache.Clear()
	return nil
}

func NewContext(urlStr string, bufferSize int) (*Context, error) {
	playlistURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	return &Context{
		playlistURL: playlistURL,
		startTime:   time.Now(),
		cache:       lru.New(bufferSize),
	}, nil
}

func (cfg *Config) GetPlaylist(urlStr string, dlc chan *Download) error {
	c, err := NewContext(urlStr, 1024)
	if err != nil {
		return err
	}
	defer c.Close()
	for {
		resp, err := Request().SetDoNotParseResponse(true).Get(urlStr)
		if err != nil {
			return err
		}
		if !resp.IsSuccess() {
			resp.RawBody().Close()
			err = fmt.Errorf("%v: [%v]%v", urlStr, resp.StatusCode(), resp.Status())
			log.Println(err)
			time.Sleep(time.Duration(3) * time.Second)
			continue
		}
		err = cfg.GetPlaylistFromReader(c, resp.RawBody(), dlc)
		resp.RawBody().Close()
		if err != nil {
			if err == ErrExit {
				return nil
			}
			return err
		}
	}
}

func (cfg *Config) GetPlaylistFromReader(c *Context, reader io.Reader, dlc chan *Download) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%w: %v", ErrPanic, e)
		}
	}()
	prog := cfg.progress
	recTime := cfg.Duration
	var playlist m3u8.Playlist
	var listType m3u8.ListType
	playlist, listType, err = m3u8.DecodeFrom(reader, true)
	if err != nil {
		return
	}

	if listType != m3u8.MEDIA {
		if listType == m3u8.MASTER {
			mpl := playlist.(*m3u8.MasterPlaylist)
			var maxBandwidth uint32
			index := -1
			for i, v := range mpl.Variants {
				if v == nil {
					continue
				}
				if v.Bandwidth > maxBandwidth {
					maxBandwidth = v.Bandwidth
					index = i
				}
			}
			if index > -1 {
				v := mpl.Variants[index]
				var msURI string
				msURI, err = ParseURI(c.playlistURL, v.URI)
				if err == nil {
					return cfg.GetPlaylist(msURI, dlc)
				}
				log.Println(err)
			}
			for _, v := range mpl.Variants {
				if v == nil {
					continue
				}
				var msURI string
				msURI, err = ParseURI(c.playlistURL, v.URI)
				if err != nil {
					log.Println(err)
					continue
				}
				return cfg.GetPlaylist(msURI, dlc)
			}
			return ErrInvalidMasterPlaylist
		}
		return ErrInvalidMediaPlaylist
	}
	mpl := playlist.(*m3u8.MediaPlaylist)
	if mpl.Key != nil {
		mpl.Key.URI, err = ParseURI(c.playlistURL, mpl.Key.URI)
		if err != nil {
			return err
		}
	}
	prog.TotalNum = len(mpl.Segments)
	dll := make([]*Download, 0, len(mpl.Segments))
	saveHistory := func() {
		b, _ := json.Marshal(dll)
		os.WriteFile(cfg.OutputFile+`._seg_`, b, 0666)
	}
	for segmentIndex, v := range mpl.Segments {
		if v == nil {
			prog.TotalNum--
			continue
		}
		msURI, err := ParseURI(c.playlistURL, v.URI)
		if err != nil {
			log.Println(err)
			prog.TotalNum--
			continue
		}
		_, hit := c.cache.Get(msURI)
		if !hit {
			c.cache.Add(msURI, nil)
			if cfg.UseLocalTime {
				c.recDuration = time.Since(c.startTime)
			} else {
				c.recDuration += time.Duration(int64(v.Duration * 1000000000))
			}
			d := &Download{
				URI:           msURI,
				ExtXKey:       mpl.Key,
				SeqNo:         uint64(segmentIndex) + mpl.SeqNo,
				totalDuration: c.recDuration,
			}
			dll = append(dll, d)
			dlc <- d
		}
		if recTime != 0 && c.recDuration != 0 && c.recDuration >= recTime {
			close(dlc)
			saveHistory()
			return ErrExit
		}
	}
	if mpl.Closed {
		close(dlc)
		saveHistory()
		return ErrExit
	}

	time.Sleep(time.Duration(int64(mpl.TargetDuration * 1000000000)))
	return nil
}
