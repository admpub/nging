package common

import (
	"testing"

	"github.com/webx-top/echo/testing/test"
)

func TestSecure(t *testing.T) {
	s := `<p>test<a href="http://www.admpub.com">link</a>test</p>`
	test.Eq(t, `<p>test<a href="http://www.admpub.com" rel="nofollow">link</a>test</p>`, RemoveXSS(s))
	s = `<p>test<a href="http://www.admpub.com"><img src="http://www.admpub.com/test" />link</a>test</p>`
	test.Eq(t, `<p>test<img src="http://www.admpub.com/test"/>linktest</p>`, RemoveXSS(s, true))
}

func TestPickCodeblock(t *testing.T) {
	str := "aaaa\n```bbb\ncode-block\n```\nb"
	pick, content := MarkdownPickoutCodeblock(str)
	test.Eq(t, []string{"bbb\ncode-block\n"}, pick)
	test.Eq(t, "aaaa\n```{codeblock(0)}```\nb", content)
	content = MarkdownRestorePickout(pick, content)
	test.Eq(t, "aaaa\n```bbb\ncode-block\n```\nb", content)

	str += "ccc\r\n```ddd\ncode-block2\n```\r\neee"
	pick, content = MarkdownPickoutCodeblock(str)
	test.Eq(t, []string{"bbb\ncode-block\n", "ddd\ncode-block2\n"}, pick)
	test.Eq(t, "aaaa\n```{codeblock(0)}```\nbccc\r\n```{codeblock(1)}```\r\neee", content)
	content = MarkdownRestorePickout(pick, content)
	test.Eq(t, "aaaa\n```bbb\ncode-block\n```\nbccc\r\n```ddd\ncode-block2\n```\r\neee", content)

	str = "aaaa\n```bbb\ncode-block\n```b```\nk"
	pick, content = MarkdownPickoutCodeblock(str)
	test.Eq(t, []string{"bbb\ncode-block\n```b"}, pick)
	test.Eq(t, "aaaa\n```{codeblock(0)}```\nk", content)
	content = MarkdownRestorePickout(pick, content)
	test.Eq(t, "aaaa\n```bbb\ncode-block\n```b\n```\nk", content)

	content = ContentEncode(`[普通链接带标题](http://localhost/ &#34;普通链接带标题&#34;)`, `markdown`)
	test.Eq(t, `[普通链接带标题](http://localhost/ "普通链接带标题")`, content)
	content = ContentEncode(`[普通链接带标题](javascript:000 &#34;普通链接带标题&#34;) 123`, `markdown`)
	test.Eq(t, `[普通链接带标题](-javascript-000 "普通链接带标题") 123`, content)
	content = ContentEncode(`[普通链接带标题](Javascript:000) abc`, `markdown`)
	test.Eq(t, `[普通链接带标题](-Javascript-000) abc`, content)
	content = ContentEncode(`![普通链接带标题](Javascript:000) abc`, `markdown`)
	test.Eq(t, `![普通链接带标题](-Javascript-000) abc`, content)
	content = ContentEncode(`![普通链接带标题](Javascript:000)abc[普通链接带标题](Javascript:111)`, `markdown`)
	test.Eq(t, `![普通链接带标题](-Javascript-000)abc[普通链接带标题](-Javascript-111)`, content)

	content = ContentEncode("123\n&gt; 123", `markdown`)
	test.Eq(t, "123\n> 123", content)

	content = ContentEncode("123\n &gt;123", `markdown`)
	test.Eq(t, "123\n >123", content)

	content = ContentEncode("&gt; 123", `markdown`)
	test.Eq(t, "> 123", content)
}
