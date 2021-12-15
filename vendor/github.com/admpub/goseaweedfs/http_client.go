package goseaweedfs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path/filepath"
	"strings"

	workerpool "github.com/linxGnu/gumble/worker-pool"
)

type httpClient struct {
	client  *http.Client
	workers *workerpool.Pool
}

func newHTTPClient(client *http.Client) *httpClient {
	c := &httpClient{
		client:  client,
		workers: createWorkerPool(),
	}
	c.workers.Start()
	return c
}

func (c *httpClient) Close() (err error) {
	c.workers.Stop()
	return
}

func (c *httpClient) get(url string, header map[string]string) (body []byte, statusCode int, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err == nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}

		var resp *http.Response
		resp, err = c.client.Do(req)
		if err == nil {
			body, statusCode, err = readAll(resp)
		}
	}

	return
}

func (c *httpClient) delete(url string) (statusCode int, err error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return
	}

	r, err := c.client.Do(req)
	if err != nil {
		return
	}

	body, statusCode, err := readAll(r)
	if err == nil {
		switch r.StatusCode {
		case http.StatusNoContent, http.StatusNotFound, http.StatusAccepted, http.StatusOK:
			err = nil
			return
		}

		m := make(map[string]interface{})
		if e := json.Unmarshal(body, &m); e == nil {
			if s, ok := m["error"].(string); ok {
				err = fmt.Errorf("Delete %s: %v", url, s)
				return
			}
		}

		err = fmt.Errorf("Delete %s. Got response but can not parse. Body:%s Code:%d", url, string(body), r.StatusCode)
	}

	return
}

func (c *httpClient) download(url string, callback func(io.Reader) error) (filename string, err error) {
	r, err := c.client.Get(url)
	if err == nil {
		if r.StatusCode != http.StatusOK {
			drainAndClose(r.Body)
			err = fmt.Errorf("Download %s but error. Status:%s", url, r.Status)
			return
		}

		contentDisposition := r.Header["Content-Disposition"]
		if len(contentDisposition) > 0 {
			if strings.HasPrefix(contentDisposition[0], "filename=") {
				filename = contentDisposition[0][len("filename="):]
				filename = strings.Trim(filename, "\"")
			}
		}

		// execute callback
		err = callback(r.Body)

		// drain and close body
		drainAndClose(r.Body)
	}

	return
}

func (c *httpClient) upload(url string, filename string, fileReader io.Reader, mtype string) (respBody []byte, statusCode int, err error) {
	r, w := io.Pipe()

	// create multipart writer
	mw := multipart.NewWriter(w)

	task := workerpool.NewTask(context.Background(), func(ctx context.Context) (interface{}, error) {
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, normalizeName(filename)))
		if mtype == "" {
			mtype = mime.TypeByExtension(strings.ToLower(filepath.Ext(filename)))
		}
		if mtype != "" {
			h.Set("Content-Type", mtype)
		}

		part, err := mw.CreatePart(h)
		if err == nil {
			_, err = io.Copy(part, fileReader)
		}

		if err == nil {
			if err = mw.Close(); err == nil {
				err = w.Close()
			} else {
				_ = w.Close()
			}
		} else {
			_ = mw.Close()
			_ = w.Close()
		}

		return nil, err
	})
	c.workers.Do(task)

	var resp *http.Response
	resp, err = c.client.Post(url, mw.FormDataContentType(), r)

	// closing reader in case Posting error.
	// This causes pipe writer fail to write and stop above task.
	_ = r.Close()

	if err == nil {
		if respBody, statusCode, err = readAll(resp); err == nil {
			result := <-task.Result()
			err = result.Err
		}
	}

	return
}
