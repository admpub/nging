// Copyright 2013 com authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package com

import (
	"html"
	"regexp"
	"strings"
)

var html2jsReplacer = strings.NewReplacer(
	`\`, `\\`,
	"\n", `\n`,
	"\r", "",
	`"`, `\"`,
)

// HTML2JS converts []byte type of HTML content into JS format.
func HTML2JS(data []byte) []byte {
	s := string(data)
	s = html2jsReplacer.Replace(s)
	return []byte(s)
}

// HTMLEncode encode html chars to string
func HTMLEncode(str string) string {
	return html.EscapeString(str)
}

// HTMLDecode decode string to html chars
func HTMLDecode(str string) string {
	return html.UnescapeString(str)
}

// HTMLDecodeAll decode string to html chars
func HTMLDecodeAll(text string) string {
	original := text
	text = HTMLDecode(text)
	if original == text {
		return text
	}
	return HTMLDecodeAll(text)
}

var (
	regexpAnyHTMLTag    = regexp.MustCompile(`<[\S\s]+?>`)
	regexpStyleHTMLTag  = regexp.MustCompile(`<(?i:style)[\S\s]+?</(?i:style)[^>]*>`)
	regexpScriptHTMLTag = regexp.MustCompile(`<(?i:script)[\S\s]+?</(?i:script)[^>]*>`)
	regexpMoreSpace     = regexp.MustCompile(`([\s]){2,}`)
	regexpMoreNewline   = regexp.MustCompile("(\n){2,}")
	regexpAnyHTMLAttr   = regexp.MustCompile(`<[/]?[\S]+[^>]*>`)
	regexpBrHTMLTag     = regexp.MustCompile("<(?i:br)[^>]*>")
)

// ClearHTMLAttr clear all attributes
func ClearHTMLAttr(src string) string {
	src = regexpAnyHTMLAttr.ReplaceAllString(src, "<$1$2>")
	return src
}

// TextLine Single line of text
func TextLine(src string) string {
	src = StripTags(src)
	return RemoveEOL(src)
}

// CleanMoreNl remove all \n(2+)
func CleanMoreNl(src string) string {
	return regexpMoreNewline.ReplaceAllString(src, "$1")
}

// CleanMoreSpace remove all spaces(2+)
func CleanMoreSpace(src string) string {
	return regexpMoreSpace.ReplaceAllString(src, "$1")
}

// StripTags strip tags in html string
func StripTags(src string) string {
	//将HTML标签全转换成小写
	//src = regexpAnyHTMLTag.ReplaceAllStringFunc(src, strings.ToLower)

	//remove tag <style>
	src = regexpStyleHTMLTag.ReplaceAllString(src, "")

	//remove tag <script>
	src = regexpScriptHTMLTag.ReplaceAllString(src, "")

	//replace all html tag into \n
	src = regexpAnyHTMLTag.ReplaceAllString(src, "\n")
	src = CleanMoreSpace(src)

	return strings.TrimSpace(src)
}

var nl2brReplacer = strings.NewReplacer(
	"\r", "",
	"\n", "<br />",
)

// Nl2br change \n to <br/>
func Nl2br(str string) string {
	return nl2brReplacer.Replace(str)
}

// Br2nl change <br/> to \n
func Br2nl(str string) string {
	return regexpBrHTMLTag.ReplaceAllString(str, "\n")
}
