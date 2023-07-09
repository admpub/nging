package progressbar

import (
	"fmt"
	"io"

	download "github.com/admpub/go-download/v2"
	"github.com/vbauerster/mpb/v6"
	"github.com/vbauerster/mpb/v6/decor"
)

func New(opt *download.Options, width ...int) *mpb.Progress {
	var w int
	if len(width) > 0 {
		w = width[0]
	}
	if w <= 0 {
		w = 80
	}
	progress := mpb.New(mpb.WithWidth(w))
	//defer progress.Wait()
	opt.Proxy = func(name string, download int, size int64, r io.Reader) io.Reader {
		bar := AddBar(progress, name, size, download)
		return bar.ProxyReader(r)
	}
	return progress
}

func AddBar(progress *mpb.Progress, name string, size int64, partNo ...int) *mpb.Bar {
	if len(partNo) > 0 {
		name = fmt.Sprintf("%s-%d", name, partNo[0])
	}
	return progress.AddBar(
		size,
		mpb.PrependDecorators(
			decor.Name(name, decor.WC{W: len(name) + 1, C: decor.DidentRight}),
			decor.CountersNoUnit(`%3d / %3d`, decor.WC{W: 18}),
		),
		mpb.AppendDecorators(
			decor.Percentage(decor.WC{W: 5}),
		),
	)
}
