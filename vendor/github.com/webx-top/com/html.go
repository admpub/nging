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

// HTML2JS converts []byte type of HTML content into JS format.
func HTML2JS(data []byte) []byte {
	s := string(data)
	s = strings.Replace(s, `\`, `\\`, -1)
	s = strings.Replace(s, "\n", `\n`, -1)
	s = strings.Replace(s, "\r", "", -1)
	s = strings.Replace(s, "\"", `\"`, -1)
	s = strings.Replace(s, "<table>", "&lt;table>", -1)
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

var (
	regexpAnyHTMLTag    = regexp.MustCompile("\\<[\\S\\s]+?\\>")
	regexpStyleHTMLTag  = regexp.MustCompile("\\<(?i:style)[\\S\\s]+?\\</(?i:style)\\>")
	regexpScriptHTMLTag = regexp.MustCompile("\\<(?i:script)[\\S\\s]+?\\</(?i:script)\\>")
	regexpMoreSpace     = regexp.MustCompile("\\s{2,}")
	regexpAnyHTMLAttr   = regexp.MustCompile("<([/]?[\\S]+)[^>]*([/]?)>")
)

// ClearHTMLAttr clear all attributes
func ClearHTMLAttr(src string) string {
	src = regexpAnyHTMLAttr.ReplaceAllString(src, "<$1$2>")
	return src
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

	//trim all spaces(2+) into \n
	src = regexpMoreSpace.ReplaceAllString(src, "\n")

	return strings.TrimSpace(src)
}

// Nl2br change \n to <br/>
func Nl2br(str string) string {
	return strings.Replace(str, "\n", "<br />", -1)
}
