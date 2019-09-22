package file

import (
	"fmt"
	"testing"

	"github.com/webx-top/com"
	"github.com/webx-top/echo/testing/test"
)

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
	result := ReplaceEmbeddedRes(content, map[uint64]string{
		1: `http://img.admpub.com/test/1.jpg`,
		2: `http://img.admpub.com/test/abc.gif`,
		3: `http://img.admpub.com/test/bg.png?1`,
		4: `http://img.admpub.com/test/bg2.png`,
		5: `http://img.admpub.com/test/bg3.png`,
	})
	fmt.Println(result)
	test.Eq(t, expected, result)
}
