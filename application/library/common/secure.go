package common

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/webx-top/com"
)

var (
	secureStrictPolicy                = bluemonday.StrictPolicy()
	secureUGCPolicy                   = bluemonday.UGCPolicy()
	secureUGCPolicyAllowDataURIImages *bluemonday.Policy
	secureUGCPolicyNoLink             = NoLink()
)

func init() {
	secureUGCPolicyAllowDataURIImages = bluemonday.UGCPolicy()
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

	// "img" is permitted
	p.AllowAttrs("align").Matching(bluemonday.ImageAlign).OnElements("img")
	p.AllowAttrs("alt").Matching(bluemonday.Paragraph).OnElements("img")
	p.AllowAttrs("height", "width").Matching(bluemonday.NumberOrPercent).OnElements("img")
	p.AllowAttrs("src").OnElements("img")

	return p
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

var q = rune('`')

func MarkdownPickoutCodeblock(content string) (repl []string, newContent string) {
	var (
		// reset
		start bool
		n     int
		code  []rune

		// keep
		keep []rune
	)

	for i, b := range content {
		if b == q {
			if start {
				if n == 2 { //终止标记“```”后面必须带换行
					index := i + 1
					if index < len(content)-1 {
						if content[index] == '\r' {
							index++
							if index > len(content)-1 {
								code = append(code, keep[len(keep)-2:]...)
								keep = keep[0 : len(keep)-2]
								code = append(code, b)
								n = 0
								continue
							}
						}
						if content[index] != '\n' {
							code = append(code, keep[len(keep)-2:]...)
							keep = keep[0 : len(keep)-2]
							code = append(code, b)
							n = 0
							continue
						}
					}
				}
			} else {
				if n == 0 { //起始标记“```”前面必须带换行
					if i > 0 {
						if content[i-1] != '\n' {
							keep = append(keep, b)
							continue
						}
					}
				}
			}
			n++
			if n == 3 {
				if start { // end
					keep = append(keep, b)
					repl = append(repl, string(code))

					start = false
					n = 0
					code = nil
					continue
				}
				code = nil
				start = true
				insert := []rune(`{codeblock(` + fmt.Sprint(len(repl)) + `)}`)
				keep = append(keep, b)
				keep = append(keep, insert...)
				n = 0
				continue
			}
			keep = append(keep, b)
			continue
		}
		if start {
			code = append(code, b)
		} else {
			keep = append(keep, b)
		}
	}
	newContent = string(keep)
	return
}

func MarkdownRestorePickout(repl []string, content string) string {
	for i, r := range repl {
		find := "```{codeblock(" + fmt.Sprint(i) + ")}```"
		// r = strings.TrimLeft(r, "\r")
		// if !strings.HasPrefix(r, "\n") {
		// 	r = "\n" + r
		// }
		if !strings.HasSuffix(r, "\n") {
			r += "\n"
		}
		content = strings.Replace(content, find, "```"+r+"```", 1)
	}
	return content
}

func ContentEncode(content string, contypes ...string) string {
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
		var pick []string
		pick, content = MarkdownPickoutCodeblock(content)
		content = RemoveXSS(content)
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
