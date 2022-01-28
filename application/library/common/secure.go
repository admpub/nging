package common

import (
	"bytes"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/webx-top/com"
)

func NewUGCPolicy() *bluemonday.Policy {
	p := bluemonday.UGCPolicy()
	allowMedia(p)
	return p
}

func NewStrictPolicy() *bluemonday.Policy {
	p := bluemonday.StrictPolicy()
	return p
}

var (
	secureStrictPolicy                = NewStrictPolicy()
	secureUGCPolicy                   = NewUGCPolicy()
	secureUGCPolicyAllowDataURIImages *bluemonday.Policy
	secureUGCPolicyNoLink             = NoLink()
)

func init() {
	secureUGCPolicyAllowDataURIImages = NewUGCPolicy()
	secureUGCPolicyAllowDataURIImages.AllowDataURIImages()
}

// ClearHTML 清除所有HTML标签及其属性，一般用处理文章标题等不含HTML标签的字符串
func ClearHTML(title string) string {
	return secureStrictPolicy.Sanitize(title)
}

// RemoveXSS 清除不安全的HTML标签和属性，一般用于处理文章内容
func RemoveXSS(content string, noLinks ...bool) string {
	if len(noLinks) > 0 && noLinks[0] {
		return secureUGCPolicyNoLink.Sanitize(content)
	}
	return secureUGCPolicy.Sanitize(content)
}

func NoLink() *bluemonday.Policy {
	p := HTMLFilter()
	p.AllowStandardAttributes()

	////////////////////////////////
	// Declarations and structure //
	////////////////////////////////

	// "xml" "xslt" "DOCTYPE" "html" "head" are not permitted as we are
	// expecting user generated content to be a fragment of HTML and not a full
	// document.

	//////////////////////////
	// Sectioning root tags //
	//////////////////////////

	// "article" and "aside" are permitted and takes no attributes
	p.AllowElements("article", "aside")

	// "body" is not permitted as we are expecting user generated content to be a fragment
	// of HTML and not a full document.

	// "details" is permitted, including the "open" attribute which can either
	// be blank or the value "open".
	p.AllowAttrs(
		"open",
	).Matching(regexp.MustCompile(`(?i)^(|open)$`)).OnElements("details")

	// "fieldset" is not permitted as we are not allowing forms to be created.

	// "figure" is permitted and takes no attributes
	p.AllowElements("figure")

	// "nav" is not permitted as it is assumed that the site (and not the user)
	// has defined navigation elements

	// "section" is permitted and takes no attributes
	p.AllowElements("section")

	// "summary" is permitted and takes no attributes
	p.AllowElements("summary")

	//////////////////////////
	// Headings and footers //
	//////////////////////////

	// "footer" is not permitted as we expect user content to be a fragment and
	// not structural to this extent

	// "h1" through "h6" are permitted and take no attributes
	p.AllowElements("h1", "h2", "h3", "h4", "h5", "h6")

	// "header" is not permitted as we expect user content to be a fragment and
	// not structural to this extent

	// "hgroup" is permitted and takes no attributes
	p.AllowElements("hgroup")

	/////////////////////////////////////
	// Content grouping and separating //
	/////////////////////////////////////

	// "blockquote" is permitted, including the "cite" attribute which must be
	// a standard URL.
	p.AllowAttrs("cite").OnElements("blockquote")

	// "br" "div" "hr" "p" "span" "wbr" are permitted and take no attributes
	p.AllowElements("br", "div", "hr", "p", "span", "wbr")

	// "area" is permitted along with the attributes that map image maps work
	p.AllowAttrs("name").Matching(
		regexp.MustCompile(`^([\p{L}\p{N}_-]+)$`),
	).OnElements("map")
	p.AllowAttrs("alt").Matching(bluemonday.Paragraph).OnElements("area")
	p.AllowAttrs("coords").Matching(
		regexp.MustCompile(`^([0-9]+,)+[0-9]+$`),
	).OnElements("area")
	p.AllowAttrs("rel").Matching(bluemonday.SpaceSeparatedTokens).OnElements("area")
	p.AllowAttrs("shape").Matching(
		regexp.MustCompile(`(?i)^(default|circle|rect|poly)$`),
	).OnElements("area")
	p.AllowAttrs("usemap").Matching(
		regexp.MustCompile(`(?i)^#[\p{L}\p{N}_-]+$`),
	).OnElements("img")

	// "link" is not permitted

	/////////////////////
	// Phrase elements //
	/////////////////////

	// The following are all inline phrasing elements
	p.AllowElements("abbr", "acronym", "cite", "code", "dfn", "em",
		"figcaption", "mark", "s", "samp", "strong", "sub", "sup", "var")

	// "q" is permitted and "cite" is a URL and handled by URL policies
	p.AllowAttrs("cite").OnElements("q")

	// "time" is permitted
	p.AllowAttrs("datetime").Matching(bluemonday.ISO8601).OnElements("time")

	////////////////////
	// Style elements //
	////////////////////

	// block and inline elements that impart no semantic meaning but style the
	// document
	p.AllowElements("b", "i", "pre", "small", "strike", "tt", "u")

	// "style" is not permitted as we are not yet sanitising CSS and it is an
	// XSS attack vector

	//////////////////////
	// HTML5 Formatting //
	//////////////////////

	// "bdi" "bdo" are permitted
	p.AllowAttrs("dir").Matching(bluemonday.Direction).OnElements("bdi", "bdo")

	// "rp" "rt" "ruby" are permitted
	p.AllowElements("rp", "rt", "ruby")

	///////////////////////////
	// HTML5 Change tracking //
	///////////////////////////

	// "del" "ins" are permitted
	p.AllowAttrs("cite").Matching(bluemonday.Paragraph).OnElements("del", "ins")
	p.AllowAttrs("datetime").Matching(bluemonday.ISO8601).OnElements("del", "ins")

	///////////
	// Lists //
	///////////

	p.AllowLists()

	////////////
	// Tables //
	////////////

	p.AllowTables()

	///////////
	// Forms //
	///////////

	// By and large, forms are not permitted. However there are some form
	// elements that can be used to present data, and we do permit those
	//
	// "button" "fieldset" "input" "keygen" "label" "output" "select" "datalist"
	// "textarea" "optgroup" "option" are all not permitted

	// "meter" is permitted
	p.AllowAttrs(
		"value",
		"min",
		"max",
		"low",
		"high",
		"optimum",
	).Matching(bluemonday.Number).OnElements("meter")

	// "progress" is permitted
	p.AllowAttrs("value", "max").Matching(bluemonday.Number).OnElements("progress")

	//////////////////////
	// Embedded content //
	//////////////////////

	// Vast majority not permitted
	// "audio" "canvas" "embed" "iframe" "object" "param" "source" "svg" "track"
	// "video" are all not permitted
	allowMedia(p)

	// "img" is permitted
	p.AllowAttrs("align").Matching(bluemonday.ImageAlign).OnElements("img")
	p.AllowAttrs("alt").Matching(bluemonday.Paragraph).OnElements("img")
	p.AllowAttrs("height", "width").Matching(bluemonday.NumberOrPercent).OnElements("img")
	p.AllowAttrs("src").OnElements("img")

	return p
}

func allowMedia(p *bluemonday.Policy) {
	p.AllowElements("picture")
	//<video webkit-playsinline="true" x-webkit-airplay="true" playsinline="true" x5-video-player-type="h5" x5-video-orientation="h5" x5-video-player-fullscreen="true" preload="auto" class="evaluate-video" src="'+source+'" poster="'+source+'?vframe/jpg/offset/1"></video>
	p.AllowAttrs("src", "controls", "width", "height", "autoplay", "muted", "loop", "poster", "preload", "playsinline", "webkit-playsinline", "x-webkit-airplay", "x5-video-player-type", "x5-video-orientation", "x5-video-player-fullscreen").OnElements("video")
	p.AllowAttrs("src", "controls", "width", "height", "autoplay", "muted", "loop", "preload").OnElements("audio")
	p.AllowAttrs("src", "type", "srcset", "media").OnElements("source")
}

func RemoveBytesXSS(content []byte, noLinks ...bool) []byte {
	if len(noLinks) > 0 && noLinks[0] {
		return secureUGCPolicyNoLink.SanitizeBytes(content)
	}
	return secureUGCPolicy.SanitizeBytes(content)
}

func RemoveReaderXSS(reader io.Reader, noLinks ...bool) *bytes.Buffer {
	if len(noLinks) > 0 && noLinks[0] {
		return secureUGCPolicyNoLink.SanitizeReader(reader)
	}
	return secureUGCPolicy.SanitizeReader(reader)
}

// HTMLFilter 构建自定义的HTML标签过滤器
func HTMLFilter() *bluemonday.Policy {
	return bluemonday.NewPolicy()
}

func MyRemoveXSS(content string) string {
	return com.RemoveXSS(content)
}

func MyCleanText(value string) string {
	value = com.StripTags(value)
	value = com.RemoveEOL(value)
	return value
}

func MyCleanTags(value string) string {
	value = com.StripTags(value)
	return value
}

var (
	q                           = rune('`')
	markdownLinkWithDoubleQuote = regexp.MustCompile(`(\]\([^ \)]+ )&#34;([^"\)]+)&#34;(\))`)
	markdownLinkWithSingleQuote = regexp.MustCompile(`(\]\([^ \)]+ )&#39;([^'\)]+)&#39;(\))`)
	markdownLinkWithScript      = regexp.MustCompile(`(?i)(\]\()(javascript):([^\)]*\))`)
	markdownQuoteTag            = regexp.MustCompile("((\n|^)[ ]{0,3})&gt;")
	markdownCodeBlock           = regexp.MustCompile("(?s)([\r\n]|^)```([\\w]*[\r\n].*?[\r\n])```([\r\n]|$)")
)

func MarkdownPickoutCodeblock(content string) (repl []string, newContent string) {
	newContent = markdownCodeBlock.ReplaceAllStringFunc(content, func(found string) string {
		placeholder := `{codeblock(` + strconv.Itoa(len(repl)) + `)}`
		leftIndex := strings.Index(found, "```")
		rightIndex := strings.LastIndex(found, "```")
		repl = append(repl, found[leftIndex+3:rightIndex])
		return found[0:leftIndex+3] + placeholder + found[rightIndex:]
	})
	//echo.Dump([]interface{}{repl, newContent, content})
	return
}

func MarkdownRestorePickout(repl []string, content string) string {
	for i, r := range repl {
		if strings.Count(r, "\n") < 2 {
			r = strings.TrimLeft(r, "\r")
			if !strings.HasPrefix(r, "\n") {
				r = "\n" + r
			}
		}
		if !strings.HasSuffix(r, "\n") {
			r += "\n"
		}
		find := "```{codeblock(" + strconv.Itoa(i) + ")}```"
		content = strings.Replace(content, find, "```"+r+"```", 1)
	}
	return content
}

func ContentEncode(content string, contypes ...string) string {
	if len(content) == 0 {
		return content
	}
	var contype string
	if len(contypes) > 0 {
		contype = contypes[0]
	}
	switch contype {
	case `html`:
		content = RemoveXSS(content)

	case `url`, `image`, `video`, `audio`, `file`, `id`:
		content = MyCleanText(content)

	case `text`:
		content = com.StripTags(content)

	case `json`:
		// pass

	case `markdown`:
		// 提取代码块
		var pick []string
		pick, content = MarkdownPickoutCodeblock(content)

		// - 删除XSS

		// 删除HTML中的XSS代码
		content = RemoveXSS(content)
		// 拦截Markdown链接中的“javascript:”
		content = markdownLinkWithScript.ReplaceAllString(content, `${1}-${2}-${3}`)

		// - 还原

		// 还原双引号
		content = markdownLinkWithDoubleQuote.ReplaceAllString(content, `${1}"${2}"${3}`)
		// 还原单引号
		content = markdownLinkWithSingleQuote.ReplaceAllString(content, `${1}'${2}'${3}`)
		// 还原引用标识
		content = markdownQuoteTag.ReplaceAllString(content, `${1}>`)
		// 还原代码块
		content = MarkdownRestorePickout(pick, content)

	case `list`:
		content = MyCleanText(content)
		content = strings.TrimSpace(content)
		content = strings.Trim(content, `,`)

	default:
		content = com.StripTags(content)
	}
	content = strings.TrimSpace(content)
	return content
}
