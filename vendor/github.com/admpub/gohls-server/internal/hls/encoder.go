package hls

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type EncodingRequest struct {
	file    string
	segment int64
	res     int64
	data    chan *[]byte
	err     chan error
}

func NewEncodingRequest(file string, segment int64, res int64) *EncodingRequest {
	return &EncodingRequest{file, segment, res, make(chan *[]byte, 1), make(chan error, 1)}
}

func NewWarmupEncodingRequest(file string, segment int64, res int64) *EncodingRequest {
	return &EncodingRequest{file, segment, res, nil, nil}
}

func (r *EncodingRequest) sendError(err error) {
	if r.err != nil {
		r.err <- err
	}
}

func (r *EncodingRequest) sendData(data *[]byte) {
	if r.data != nil {
		r.data <- data
	}
}

func (r *EncodingRequest) getCacheKey() string {
	h := sha1.New()
	h.Write([]byte(r.file))
	return fmt.Sprintf("%x.%v.%v", h.Sum(nil), r.res, r.segment)
}

type Encoder struct {
	cacheDir string
	reqChan  chan EncodingRequest
}

func NewEncoder(cacheDir string, workerCount int) *Encoder {
	rc := make(chan EncodingRequest, 100)
	encoder := &Encoder{cacheDir, rc}
	go func() {
		for {
			r := <-rc
			cache, err := encoder.GetFromCache(r)
			if err != nil {
				r.sendError(err)
				continue
			}
			if cache != nil {
				r.sendData(&cache)
				continue
			}
			log.Debugf("Encoding %v:%v", r.file, r.segment)
			data, err := execute(FFMPEGPath, EncodingArgs(r.file, r.segment, r.res))
			if err != nil {
				r.err <- err
				continue
			}
			r.sendData(&data)
			tmp := encoder.GetCacheFile(r) + ".tmp"
			mkerr := os.MkdirAll(filepath.Join(HomeDir, cacheDirName, encoder.cacheDir), 0777)
			if mkerr != nil {
				log.Errorf("Could not create cache dir")
				continue
			}
			if err2 := ioutil.WriteFile(tmp, data, 0777); err2 == nil {
				os.Rename(tmp, encoder.GetCacheFile(r))
			}
		}
	}()
	return encoder
}

func (e *Encoder) GetFromCache(r EncodingRequest) ([]byte, error) {

	cachePath := e.GetCacheFile(r)
	if _, err := os.Stat(cachePath); err != nil {
		// Cache file could not be openened ....
		if os.IsNotExist(err) {
			// Because segment has not yet been encoded
			return nil, nil
		}
		// Cache file could not be opened because of an underlying error
		return nil, fmt.Errorf("Encoder cache file %v could not be opened because: %v", cachePath, err)
	}
	// The file could be opened, read it's content
	dat, err := ioutil.ReadFile(cachePath)
	if err != nil {
		return nil, fmt.Errorf("Encoder could not read cache file %v because: %v", cachePath, err)
	}
	// file was read successfully
	return dat, nil
}

func (e *Encoder) GetCacheFile(r EncodingRequest) string {
	return filepath.Join(HomeDir, cacheDirName, e.cacheDir, r.getCacheKey())
}

func (e *Encoder) Encode(r EncodingRequest) {
	go func() {
		log.Debugf("Encoding requested %v:%v", r.file, r.segment)
		// This needs to run in it's own go routine because channel writes block
		data, err := e.GetFromCache(r)
		if err != nil {
			r.sendError(err)
			return
		}
		if data != nil {
			r.sendData(&data)
			return
		}
		e.reqChan <- r
		e.reqChan <- *NewWarmupEncodingRequest(r.file, r.segment+1, r.res)
		e.reqChan <- *NewWarmupEncodingRequest(r.file, r.segment+2, r.res)
	}()
}

func EncodingArgs(videoFile string, segment int64, res int64) []string {
	startTime := segment * hlsSegmentLength
	// see http://superuser.com/questions/908280/what-is-the-correct-way-to-fix-keyframes-in-ffmpeg-for-dash
	return []string{
		// Prevent encoding to run longer than 30 seonds
		"-timelimit", "45",

		// TODO: Some stuff to investigate
		// "-probesize", "524288",
		// "-fpsprobesize", "10",
		// "-analyzeduration", "2147483647",
		// "-hwaccel:0", "vda",

		// The start time
		// important: needs to be before -i to do input seeking
		"-ss", fmt.Sprintf("%v.00", startTime),

		// The source file
		"-i", videoFile,

		// Put all streams to output
		// "-map", "0",

		// The duration
		"-t", fmt.Sprintf("%v.00", hlsSegmentLength),

		// TODO: Find out what it does
		//"-strict", "-2",

		// 720p
		"-vf", fmt.Sprintf("scale=-2:%v", res),

		// x264 video codec
		"-vcodec", "libx264",

		// x264 preset
		"-preset", "veryfast",

		// aac audio codec
		"-acodec", "aac",
		//
		"-pix_fmt", "yuv420p",

		//"-r", "25", // fixed framerate

		"-force_key_frames", "expr:gte(t,n_forced*5.000)",

		//"-force_key_frames", "00:00:00.00",
		//"-x264opts", "keyint=25:min-keyint=25:scenecut=-1",

		//"-f", "mpegts",

		"-f", "ssegment",
		"-segment_time", fmt.Sprintf("%v.00", hlsSegmentLength),
		"-initial_offset", fmt.Sprintf("%v.00", startTime),

		"pipe:out%03d.ts",
	}
}
