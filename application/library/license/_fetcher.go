package license

import (
	"io"
	"sync/atomic"

	"github.com/jpillora/overseer/fetcher"
	"github.com/webx-top/echo/defaults"
)

type HTTPFecher struct {
	version atomic.Value
	*fetcher.HTTP
}

func (h *HTTPFecher) SetVersion(version string) {
	h.version.Store(version)
}

func (h *HTTPFecher) SetURL(urlStr string) {
	h.URL = urlStr
}

func (h *HTTPFecher) Init() error {
	h.SetURL(`none`)
	h.SetVersion(``)
	return h.HTTP.Init()
}

func (h *HTTPFecher) Fetch() (io.Reader, error) {
	ctx := defaults.NewMockContext()
	info, err := LatestVersion(ctx, h.version.Load().(string), false)
	if err != nil {
		return nil, err
	}
	h.SetURL(info.DownloadURL)
	return h.HTTP.Fetch()
}
