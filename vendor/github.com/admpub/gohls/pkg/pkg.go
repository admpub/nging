package pkg

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
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
	Client        = &http.Client{}
	IVPlaceholder = []byte{0, 0, 0, 0, 0, 0, 0, 0}
)

func DoRequest(c *http.Client, req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", UserAgent)
	//req.Header.Set("Connection", "Keep-Alive") //http2不支持Keep-Alive
	resp, err := c.Do(req)
	if err != nil {
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
			req, err := http.NewRequest("GET", v.ExtXKey.URI, nil)
			if err != nil {
				log.Println(err)
			}
			resp, err := DoRequest(Client, req)
			if err != nil {
				log.Println(err)
			}
			keyData, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println(err)
			}
			resp.Body.Close()
			(*aes128Keys)[v.ExtXKey.URI] = keyData
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

func DownloadSegment(fn string, dlc chan *Download, recTime time.Duration, prog *Progress) error {
	var out, err = os.Create(fn)
	defer out.Close()

	if err != nil {
		return err
	}
	var (
		data       []byte
		aes128Keys = &map[string][]byte{}
	)

	defer func() {
		if e := recover(); e != nil {
			log.Println(e)
		}
	}()
	for v := range dlc {
		prog.FinishedNum++
		startTime := time.Now()
		req, err := http.NewRequest("GET", v.URI, nil)
		if err != nil {
			return err
		}
		resp, err := DoRequest(Client, req)
		if err != nil {
			log.Print(err)
			continue
		}
		if resp.StatusCode != 200 {
			log.Printf("Received HTTP %v for %v\n", resp.StatusCode, v.URI)
			resp.Body.Close()
			continue
		}

		data, _ = ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		DecryptData(data, v, aes128Keys)
		var size int
		size, err = out.Write(data)
		// _, err = io.Copy(out, resp.Body)
		if err != nil {
			return err
		}

		processSize := len(fmt.Sprint(prog.TotalNum))
		prog.FinishedSize += int64(size)
		prog.SpeedInSecond = float64(size) / time.Now().Sub(startTime).Seconds()
		speedDesc := com.FormatBytes(prog.SpeedInSecond)

		log.Printf("[%0*d/%d][%v/s] Downloaded %v\n", processSize, prog.FinishedNum, prog.TotalNum, speedDesc, v.URI)
		if recTime != 0 {
			log.Printf("Recorded %v of %v\n", v.totalDuration, recTime)
		} else {
			log.Printf("Recorded %v\n", v.totalDuration)
		}
	}
	return nil
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
		msURL, err := root.Parse(msURI)
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
	defer func() {
		if e := recover(); e != nil {
			log.Println(e)
		}
		c.Close()
	}()
	for {
		req, err := http.NewRequest("GET", urlStr, nil)
		if err != nil {
			return err
		}
		resp, err := DoRequest(Client, req)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Duration(3) * time.Second)
		}

		err = cfg.GetPlaylistFromReader(c, resp.Body, dlc)
		if err != nil {
			if err == ErrExit {
				return nil
			}
			return err
		}
	}
}

func (cfg *Config) GetPlaylistFromReader(c *Context, reader io.Reader, dlc chan *Download) error {
	prog := cfg.progress
	recTime := cfg.Duration
	playlist, listType, err := m3u8.DecodeFrom(reader, true)
	if err != nil {
		return err
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
				msURI, err := ParseURI(c.playlistURL, v.URI)
				if err == nil {
					return cfg.GetPlaylist(msURI, dlc)
				}
				log.Println(err)
			}
			for _, v := range mpl.Variants {
				if v == nil {
					continue
				}
				msURI, err := ParseURI(c.playlistURL, v.URI)
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
	prog.TotalNum = len(mpl.Segments)
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
				c.recDuration = time.Now().Sub(c.startTime)
			} else {
				c.recDuration += time.Duration(int64(v.Duration * 1000000000))
			}
			dlc <- &Download{
				URI:           msURI,
				ExtXKey:       mpl.Key,
				SeqNo:         uint64(segmentIndex) + mpl.SeqNo,
				totalDuration: c.recDuration,
			}
		}
		if recTime != 0 && c.recDuration != 0 && c.recDuration >= recTime {
			close(dlc)
			return ErrExit
		}
	}
	if mpl.Closed {
		close(dlc)
		return ErrExit
	}

	time.Sleep(time.Duration(int64(mpl.TargetDuration * 1000000000)))
	return nil
}
