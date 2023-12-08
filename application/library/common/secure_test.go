package common

import (
	"testing"

	"github.com/webx-top/echo/testing/test"
)

func TestSplitSingleMultibytes(t *testing.T) {
	expected := `中 a`
	actual := SplitSingleMutibytes(`中a`)
	test.Eq(t, expected, actual)

	expected = `a 中`
	actual = SplitSingleMutibytes(`a中`)
	test.Eq(t, expected, actual)

	expected = `左边 a 右边`
	actual = SplitSingleMutibytes(`左边a右边`)
	test.Eq(t, expected, actual)
}

func TestSecure(t *testing.T) {
	test.Eq(t, `/c/`, mysqlNetworkRegexp.ReplaceAllString(`//c/`, `/`))
	s := `<p>test<a href="http://www.admpub.com">link</a>test</p>`
	test.Eq(t, `<p>test<a href="http://www.admpub.com" rel="nofollow">link</a>test</p>`, RemoveXSS(s))
	s = `<p>test<a href="http://www.admpub.com"><img src="http://www.admpub.com/test" />link</a>test</p>`
	test.Eq(t, `<p>test<img src="http://www.admpub.com/test"/>linktest</p>`, RemoveXSS(s, true))
	s = `<video src="123">`
	test.Eq(t, `<video src="123">`, RemoveXSS(s, true))
	s = `<pre class="language-javascript">`
	test.Eq(t, `<pre class="language-javascript">`, RemoveXSS(s))
	s = `<ol start="4">`
	test.Eq(t, `<ol start="4">`, RemoveXSS(s))
	s = "<\nimg\n/\nonload=\"alert('OK')\">"
	test.Eq(t, "&lt;\nimg\n/\nonload=&#34;alert(&#39;OK&#39;)&#34;&gt;", RemoveXSS(s))
	s = `<img style="display: block; margin-left: auto; margin-right: auto;" src="http://www.admpub.com/test/">`
	test.Eq(t, [][]string{
		[]string{"display: block; margin-left: auto; margin-right: auto;"},
	}, styleListRegex.FindAllStringSubmatch("display: block; margin-left: auto; margin-right: auto;", -1))
	test.Eq(t, [][]string{
		[]string{" display: block; margin-left: auto; margin-right: auto; "},
	}, styleListRegex.FindAllStringSubmatch(" display: block; margin-left: auto; margin-right: auto; ", -1))
	test.Eq(t, [][]string{
		[]string{" display: block; margin-left: auto; margin-right: auto "},
	}, styleListRegex.FindAllStringSubmatch(" display: block; margin-left: auto; margin-right: auto ", -1))
	test.Eq(t, `<img style="display: block; margin-left: auto; margin-right: auto;" src="http://www.admpub.com/test/">`, RemoveXSS(s))
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

	str = "aaaa\n```bbb\ncode-block\n```b\n```\nk"
	pick, content = MarkdownPickoutCodeblock(str)
	test.Eq(t, []string{"bbb\ncode-block\n```b\n"}, pick)
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
