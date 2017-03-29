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

// Html2JS converts []byte type of HTML content into JS format.
func Html2JS(data []byte) []byte {
	s := string(data)
	s = strings.Replace(s, `\`, `\\`, -1)
	s = strings.Replace(s, "\n", `\n`, -1)
	s = strings.Replace(s, "\r", "", -1)
	s = strings.Replace(s, "\"", `\"`, -1)
	s = strings.Replace(s, "<table>", "&lt;table>", -1)
	return []byte(s)
}

// encode html chars to string
func HtmlEncode(str string) string {
	return html.EscapeString(str)
}

// decode string to html chars
func HtmlDecode(str string) string {
	return html.UnescapeString(str)
}

var (
	regexpAnyHtmlTag    = regexp.MustCompile("\\<[\\S\\s]+?\\>")
	regexpStyleHtmlTag  = regexp.MustCompile("\\<(?i:style)[\\S\\s]+?\\</(?i:style)\\>")
	regexpScriptHtmlTag = regexp.MustCompile("\\<(?i:script)[\\S\\s]+?\\</(?i:script)\\>")
	regexpMoreSpace     = regexp.MustCompile("\\s{2,}")
	regexpAnyHtmlAttr   = regexp.MustCompile("<([/]?[\\S]+)[^>]*([/]?)>")
)

// clear all attributes
func ClearHtmlAttr(src string) string {
	src = regexpAnyHtmlAttr.ReplaceAllString(src, "<$1$2>")
	return src
}

// strip tags in html string
func StripTags(src string) string {
	//将HTML标签全转换成小写
	//src = regexpAnyHtmlTag.ReplaceAllStringFunc(src, strings.ToLower)

	//remove tag <style>
	src = regexpStyleHtmlTag.ReplaceAllString(src, "")

	//remove tag <script>
	src = regexpScriptHtmlTag.ReplaceAllString(src, "")

	//replace all html tag into \n
	src = regexpAnyHtmlTag.ReplaceAllString(src, "\n")

	//trim all spaces(2+) into \n
	src = regexpMoreSpace.ReplaceAllString(src, "\n")

	return strings.TrimSpace(src)
}

// change \n to <br/>
func Nl2br(str string) string {
	return strings.Replace(str, "\n", "<br />", -1)
}
