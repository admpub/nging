package model

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"

	"github.com/admpub/goseaweedfs/libs"
)

// File structure according to filer API at https://github.com/chrislusf/seaweedfs/wiki/Filer-Server-API.
type File struct {
	FileID string `json:"fid"`
	Name   string `json:"name"`
}

// Dir directory of filer. According to https://github.com/chrislusf/seaweedfs/wiki/Filer-Server-API.
type Dir struct {
	Path    string `json:"Directory"`
	Files   []*File
	Subdirs []*File `json:"Subdirectories"`
}

// Filer client
type Filer struct {
	URL        string `json:"url"`
	HTTPClient *libs.HTTPClient
}

// FilerUploadResult upload result which responsed from filer server. According to https://github.com/chrislusf/seaweedfs/wiki/Filer-Server-API.
type FilerUploadResult struct {
	Name    string `json:"name,omitempty"`
	FileURL string `json:"url,omitempty"`
	FileID  string `json:"fid,omitempty"`
	Size    int64  `json:"size,omitempty"`
	Error   string `json:"error,omitempty"`
}

// NewFiler new filer with filer server's url
func NewFiler(url string, httpClient *libs.HTTPClient) *Filer {
	if !strings.HasPrefix(url, "http:") && !strings.HasPrefix(url, "https:") {
		url = "http://" + url
	}

	return &Filer{
		URL:        url,
		HTTPClient: httpClient,
	}
}

// Dir list in directory
func (f *Filer) Dir(pathname string) (result *Dir, err error) {
	if !strings.HasPrefix(pathname, "/") {
		pathname = "/" + pathname
	}
	if !strings.HasSuffix(pathname, "/") {
		pathname = pathname + "/"
	}

	data, _, err := f.HTTPClient.GetWithHeaders(f.URL+pathname, map[string]string{"Accept": "application/json"})
	if err != nil {
		return nil, err
	}

	result = &Dir{}
	if err = json.Unmarshal(data, result); err != nil {
		return
	}

	return
}

func (f *Filer) DownloadFile(fileURL string) (string, []byte, error) {
	fileName, rc, err := f.HTTPClient.DownloadFromURL(f.URL + fileURL)
	if err != nil {
		return "", nil, err
	}
	defer rc.Close()
	fileData, err := ioutil.ReadAll(rc)
	if err != nil {
		return "", nil, err
	}
	return fileName, fileData, nil
}

// Download download file from url.
// Note: rc must be closed after finishing as other ReadCloser.
func (f *Filer) Download(fileURL string) (string, io.ReadCloser, error) {
	fileName, rc, err := f.HTTPClient.DownloadFromURL(f.URL + fileURL)
	return fileName, rc, err
}

// UploadFile a file
func (f *Filer) UploadFile(filePath, newFilerPath, collection, ttl string) (result *FilerUploadResult, err error) {
	fp, err := NewFilePart(filePath)
	if err != nil {
		return
	}
	fp.Collection = collection
	fp.TTL = ttl

	if !strings.HasPrefix(newFilerPath, "/") {
		newFilerPath = "/" + newFilerPath
	}

	data, _, err := f.HTTPClient.Upload(f.URL+newFilerPath, filePath, fp.Reader, fp.IsGzipped, fp.MimeType)
	if err != nil {
		return
	}

	result = &FilerUploadResult{}
	if err = json.Unmarshal(data, result); err != nil {
		return
	}

	return
}

// Upload file by reader
func (f *Filer) Upload(reader io.Reader, fileSize int64, newFilerPath, collection, ttl string) (result *FilerUploadResult, err error) {
	fp := NewFilePartFromReader(reader, newFilerPath, fileSize)
	fp.Collection = collection
	fp.TTL = ttl

	if !strings.HasPrefix(newFilerPath, "/") {
		newFilerPath = "/" + newFilerPath
	}

	data, _, err := f.HTTPClient.Upload(f.URL+newFilerPath, newFilerPath, fp.Reader, fp.IsGzipped, fp.MimeType)
	if err != nil {
		return
	}

	result = &FilerUploadResult{}
	if err = json.Unmarshal(data, result); err != nil {
		return
	}

	return
}

// Delete a file/dir
func (f *Filer) Delete(pathname string, recursive ...bool) (err error) {
	if !strings.HasPrefix(pathname, "/") {
		pathname = "/" + pathname
	}
	if len(recursive) > 0 && recursive[0] {
		_, err = f.HTTPClient.Delete(f.URL + pathname + "?recursive=true")
		return
	}
	_, err = f.HTTPClient.Delete(f.URL + pathname)
	return
}
