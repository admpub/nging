package com

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	ErrNoHeaderContentLength = errors.New(`No Content-Length Provided`)
	ErrMd5Unmatched          = errors.New("WARNING: MD5 Sums don't match")
)

func RangeDownload(url string, saveTo string, args ...int) error {
	threads := 10
	if len(args) > 0 {
		threads = args[0]
	}
	defer timeTrack(time.Now(), "Full download")
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	contentLength := resp.Header.Get("Content-Length")
	if len(contentLength) < 1 {
		return ErrNoHeaderContentLength
	}
	var startByte int64
	outfile, err := os.OpenFile(saveTo, os.O_RDWR, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			outfile, err = os.Create(saveTo)
		}
	} else {
		stat, err := outfile.Stat()
		if err == nil {
			outfile.Seek(stat.Size(), 0)
			startByte = stat.Size()
		}
	}
	if outfile != nil {
		defer outfile.Close()
	}
	if err != nil {
		return err
	}
	contentSize, _ := strconv.ParseInt(contentLength, 10, 64)
	if resp.Header.Get("Accept-Ranges") == "bytes" {
		var wg sync.WaitGroup
		log.Println("Ranges Supported!")
		log.Println("Content Size:", contentLength, `(`+FormatByte(contentSize)+`)`)
		if contentSize <= startByte {
			log.Println("Download Complete! Total Size:", contentSize, `(`+FormatByte(contentSize)+`)`)
			return nil
		}
		contentSize -= startByte
		calculatedChunksize := contentSize / int64(threads)
		log.Println("Chunk Size: ", calculatedChunksize, `(`+FormatByte(calculatedChunksize)+`)`)
		var endByte int64
		chunks := 0
		completedChunks := 0
		totalChunks := threads
		if math.Mod(float64(contentSize), float64(threads)) > 0 {
			totalChunks++
		}
		lengthStr := strconv.Itoa(len(strconv.Itoa(totalChunks)))
		completedChunkCallback := func() {
			completedChunks++
			log.Println(`Completed`, saveTo, `chunks:`, fmt.Sprintf(`%`+lengthStr+`d`, completedChunks), `/`, totalChunks)
		}
		for i := 0; i < threads; i++ {
			wg.Add(1)
			endByte = startByte + calculatedChunksize
			go fetchChunk(startByte, endByte, url, outfile, &wg, completedChunkCallback)
			startByte = endByte
			chunks++
		}
		if endByte < contentSize {
			wg.Add(1)
			startByte = endByte
			endByte = contentSize
			go fetchChunk(startByte, endByte, url, outfile, &wg, completedChunkCallback)
			chunks++
		}
		wg.Wait()
		log.Println("Download Complete! Total Size:", contentSize, `(`+FormatByte(contentSize)+`)`)
		log.Println("Building File...")
		defer timeTrack(time.Now(), "File Assembled")
		//Verify file size
		filestats, err := outfile.Stat()
		if err != nil {
			return err
		}
		actualFileSize := filestats.Size()
		if actualFileSize != contentSize {
			return errors.New(fmt.Sprint("Actual Size: ", actualFileSize, " Expected: ", contentSize))
		}
		//Verify Md5
		fileHash := resp.Header.Get("X-File-Hash")
		if len(fileHash) == 0 {
			if len(resp.Header["X-Goog-Hash"]) > 1 {
				if len(resp.Header["X-Goog-Hash"][1]) > 4 {
					fileHash = resp.Header["X-Goog-Hash"][1][4:]
				}
			}
		}
		if len(fileHash) > 0 {
			contentMd5, err := hex.DecodeString(fileHash)
			if err != nil {
				return err
			}
			barray, _ := os.ReadFile(saveTo)
			computedHash := md5.Sum(barray)
			computedSlice := computedHash[0:]
			if bytes.Compare(computedSlice, contentMd5) != 0 {
				return ErrMd5Unmatched
			}
			//log.Println("File MD5 Matches!")
		}
		log.Println("File Build Complete!")
		return nil
	}
	log.Println("Range Download unsupported")
	log.Println("Beginning full download...")
	err = fetchChunk(0, contentSize, url, outfile, nil, nil)
	log.Println("Download Complete")
	return err
}

func assembleChunk(filename string, outfile *os.File) error {
	chunkFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer chunkFile.Close()
	_, err = io.Copy(outfile, chunkFile)
	if err != nil {
		return err
	}
	return os.Remove(filename)
}

func fetchChunk(startByte, endByte int64, url string, outfile *os.File, wg *sync.WaitGroup, callback func()) error {
	if wg != nil {
		defer wg.Done()
	}
	client := new(http.Client)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			log.Println(err)
			return
		}
		//log.Println("Finished Downloading byte ", startByte, `(`+FormatByte(startByte)+`)`)
		if callback != nil {
			callback()
		}
	}()

	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", startByte, endByte-1))
	res, err := client.Do(req)
	/*
		var retry int = 3
		var res *http.Response
		for i := retry; i > 0; i-- {
			res, err = client.Do(req)
			if res.StatusCode == 200 {
				retry = 3
				break
			}
			retry = i
		}
		if retry == 0 && res == nil {
			log.Fatal(err)
			return
		}
	*/
	if err != nil {
		return err
	}
	defer res.Body.Close()
	ra, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	_, err = outfile.WriteAt(ra, startByte)
	return err
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
