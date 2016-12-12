// Copyright 2013 Dmitry Chestnykh. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package htmlmin minifies HTML.
package minify

import (
	"bytes"
	"io"

	"golang.org/x/net/html"
)

type Options struct {
	MinifyScripts bool // if true, use jsmin to minify contents of script tags.
	MinifyStyles  bool // if true, use cssmin to minify contents of style tags and inline styles.
}

var DefaultOptions = &Options{
	MinifyScripts: false,
	MinifyStyles:  false,
}

var FullOptions = &Options{
	MinifyScripts: true,
	MinifyStyles:  true,
}

// Minify returns minified version of the given HTML data.
// If passed options is nil, uses default options.
func MinifyHTML(data []byte, options *Options) (out []byte, err error) {
	if options == nil {
		options = DefaultOptions
	}
	var b bytes.Buffer
	z := html.NewTokenizer(bytes.NewReader(data))
	raw := false
	javascript := false
	style := false
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			err := z.Err()
			if err == io.EOF {
				return b.Bytes(), nil
			}
			return nil, err
		case html.StartTagToken, html.SelfClosingTagToken:
			tagName, hasAttr := z.TagName()
			switch string(tagName) {
			case "script":
				javascript = true
				raw = true
			case "style":
				style = true
				raw = true
			case "pre", "code", "textarea":
				raw = true
			default:
				raw = false
			}
			b.WriteByte('<')
			b.Write(tagName)
			var k, v []byte
			isFirst := true
			for hasAttr {
				k, v, hasAttr = z.TagAttr()
				if javascript && string(k) == "type" && string(v) != "text/javascript" {
					javascript = false
				}
				if string(k) == "style" && options.MinifyStyles {
					v = []byte("a{" + string(v) + "}") // simulate "full" CSS
					v = MinifyCSS(v)
					v = v[2 : len(v)-1] // strip simulation
				}
				if isFirst {
					b.WriteByte(' ')
					isFirst = false
				}
				b.Write(k)
				b.WriteByte('=')
				if quoteChar := valueQuoteChar(v); quoteChar != 0 {
					// Quoted value.
					b.WriteByte(quoteChar)
					b.WriteString(html.EscapeString(string(v)))
					b.WriteByte(quoteChar)
				} else {
					// Unquoted value.
					b.Write(v)
				}
				if hasAttr {
					b.WriteByte(' ')
				}
			}
			b.WriteByte('>')
		case html.EndTagToken:
			tagName, _ := z.TagName()
			raw = false
			if javascript && string(tagName) == "script" {
				javascript = false
			}
			if style && string(tagName) == "style" {
				style = false
			}
			b.Write([]byte("</"))
			b.Write(tagName)
			b.WriteByte('>')
		case html.CommentToken:
			if bytes.HasPrefix(z.Raw(), []byte("<!--[if")) ||
				bytes.HasPrefix(z.Raw(), []byte("<!--//")) {
				// Preserve IE conditional and special style comments.
				b.Write(z.Raw())
			}
			// ... otherwise, skip.
		case html.TextToken:
			if javascript && options.MinifyScripts {
				min, err := MinifyJS(z.Raw())
				if err != nil {
					// Just write it as is.
					b.Write(z.Raw())
				} else {
					b.Write(min)
				}
			} else if style && options.MinifyStyles {
				b.Write(MinifyCSS(z.Raw()))
			} else if raw {
				b.Write(z.Raw())
			} else {
				b.Write(trimTextToken(z.Raw()))
			}
		default:
			b.Write(z.Raw())
		}

	}
}

func trimTextToken(b []byte) (out []byte) {
	out = make([]byte, 0)
	seenSpace := false
	for _, c := range b {
		switch c {
		case ' ', '\n', '\r', '\t':
			if !seenSpace {
				out = append(out, c)
				seenSpace = true
			}
		default:
			out = append(out, c)
			seenSpace = false
		}
	}
	return out
}

func valueQuoteChar(b []byte) byte {
	if len(b) == 0 || bytes.IndexAny(b, "'`=<> \n\r\t\b") != -1 {
		return '"' // quote with quote mark
	}
	if bytes.IndexByte(b, '"') != -1 {
		return '\'' // quote with apostrophe
	}
	return 0 // do not quote
}
