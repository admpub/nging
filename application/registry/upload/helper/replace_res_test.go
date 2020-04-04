package helper

import (
	"fmt"
	"testing"

	"github.com/webx-top/com"
	"github.com/webx-top/echo/testing/test"
)

func TestPlaceholders(t *testing.T) {
	baseURLs := map[string]string{
		`1`: `https://img1.admpub.com`,
		`2`: `https://img2.admpub.com`,
		`3`: `https://img3.admpub.com`,
		`4`: `https://img4.admpub.com`,
		`5`: `https://img5.admpub.com`,
	}
	repl := func(id string) string {
		return baseURLs[id]
	}
	urlList := map[string]string{}
	filePath := `/1232/1232/ok.jpg`
	for k, v := range baseURLs {
		urlList[v+filePath] = `[storage:`+k+`]`+filePath
	}
	for k, v := range urlList{
		test.Eq(t, k, ReplacePlaceholder(v, repl))
	}
}

func TestEmbedded(t *testing.T) {
	content := `<img src="http://www.admpub.com/test/1.jpg">[example](/test/abc.gif)<style type="text/css">
	.test{background-image:url('/test/bg.png?1')}
	.test2{background-image:url("/test/bg2.png")}
	.test3{background-image:url(/test/bg3.png)}
	</style>`
	expected := []string{
		`http://www.admpub.com/test/1.jpg`,
		`/test/abc.gif`,
		`/test/bg.png?1`,
		`/test/bg2.png`,
		`/test/bg3.png`,
	}
	capture := []string{}
	results := EmbeddedRes(content, func(file string, _ int64) {
		capture = append(capture, file)
	})
	com.Dump(results)
	test.Eq(t, expected, capture)
}

func TestReplaceEmbedded(t *testing.T) {
	content := `<img src="http://www.admpub.com/test/1.jpg#FileID-1">[example](/test/abc.gif#FileID-2)<style type="text/css">
	.test{background-image:url('/test/bg.png?1#FileID-3')}
	.test2{background-image:url("/test/bg2.png#FileID-4")}
	.test3{background-image:url(/test/bg3.png#FileID-5)}
	</style>`
	expected := `<img src="http://img.admpub.com/test/1.jpg">[example](http://img.admpub.com/test/abc.gif)<style type="text/css">
	.test{background-image:url('http://img.admpub.com/test/bg.png?1')}
	.test2{background-image:url("http://img.admpub.com/test/bg2.png")}
	.test3{background-image:url(http://img.admpub.com/test/bg3.png)}
	</style>`
	result := ReplaceEmbeddedResID(content, map[uint64]string{
		1: `http://img.admpub.com/test/1.jpg`,
		2: `http://img.admpub.com/test/abc.gif`,
		3: `http://img.admpub.com/test/bg.png?1`,
		4: `http://img.admpub.com/test/bg2.png`,
		5: `http://img.admpub.com/test/bg3.png`,
	})
	fmt.Println(result)
	test.Eq(t, expected, result)
}

func TestRelated(t *testing.T) {
	var files []string
	var fids []int64
	RelatedRes(`http://www.admpub.com/test/1.jpg#FileID-1`, func(file string, fid int64) {
		files = append(files, file)
		fids = append(fids, fid)
	})

	fmt.Println(files)
	test.Eq(t, []string{`http://www.admpub.com/test/1.jpg`}, files)
	test.Eq(t, []int64{1}, fids)
}

func TestRelated2(t *testing.T) {
	var files []string
	var fids []int64
	RelatedRes(`http://www.admpub.com/test/1.jpg,http://www.admpub.com/test/2.jpg`, func(file string, fid int64) {
		files = append(files, file)
		fids = append(fids, fid)
	}, `,`)

	fmt.Println(files)
	test.Eq(t, []string{`http://www.admpub.com/test/1.jpg`, `http://www.admpub.com/test/2.jpg`}, files)
}
