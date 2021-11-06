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

import "regexp"

const (
	regexEmailPattern       = `(?i)[A-Z0-9._%+-]+@(?:[A-Z0-9-]+\.)+[A-Z]{2,6}`
	regexStrictEmailPattern = `(?i)[A-Z0-9!#$%&'*+/=?^_{|}~-]+` +
		`(?:\.[A-Z0-9!#$%&'*+/=?^_{|}~-]+)*` +
		`@(?:[A-Z0-9](?:[A-Z0-9-]*[A-Z0-9])?\.)+` +
		`[A-Z0-9](?:[A-Z0-9-]*[A-Z0-9])?`
	regexURLPattern             = `(ftp|http|https):\/\/(\w+:{0,1}\w*@)?(\S+)(:[0-9]+)?(\/|\/([\w#!:.?+=&%@!\-\/]))?`
	regexUsernamePattern        = `^[\w\p{Han}]+$`
	regexChinesePattern         = `^[\p{Han}]+$`
	regexChineseFirstPattern    = `^[\p{Han}]`
	regexContainsChinesePattern = `[\p{Han}]+`
	regexEOLPattern             = "[\r\n]+"
)

var (
	regexEmail           = regexp.MustCompile(regexEmailPattern)
	regexStrictEmail     = regexp.MustCompile(regexStrictEmailPattern)
	regexURL             = regexp.MustCompile(regexURLPattern)
	regexUsername        = regexp.MustCompile(regexUsernamePattern)
	regexChinese         = regexp.MustCompile(regexChinesePattern)
	regexContainsChinese = regexp.MustCompile(regexContainsChinesePattern)
	regexChinesePrefix   = regexp.MustCompile(regexChineseFirstPattern)
	regexEOL             = regexp.MustCompile(regexEOLPattern)
	regexFloat           = regexp.MustCompile(`^[-]?[\d]+\.[\d]+$`)
	regexInteger         = regexp.MustCompile(`^[-]?[\d]+$`)
	regexUnsignedInteger = regexp.MustCompile(`^[\d]+$`)
)

// IsEmail validate string is an email address, if not return false
// basically validation can match 99% cases
func IsEmail(email string) bool {
	return regexEmail.MatchString(email)
}

// IsEmailRFC validate string is an email address, if not return false
// this validation omits RFC 2822
func IsEmailRFC(email string) bool {
	return regexStrictEmail.MatchString(email)
}

// IsURL validate string is a url link, if not return false
// simple validation can match 99% cases
func IsURL(url string) bool {
	return regexURL.MatchString(url)
}

// IsUsername validate string is a available username
func IsUsername(username string) bool {
	return regexUsername.MatchString(username)
}

// IsChinese validate string is Chinese
func IsChinese(str string) bool {
	return regexChinese.MatchString(str)
}

// HasChinese contains Chinese
func HasChinese(str string) bool {
	return regexContainsChinese.MatchString(str)
}

func HasChineseFirst(str string) bool {
	return regexChinesePrefix.MatchString(str)
}

// IsSingleLineText validate string is a single-line text
func IsSingleLineText(text string) bool {
	return !regexEOL.MatchString(text)
}

// IsMultiLineText validate string is a multi-line text
func IsMultiLineText(text string) bool {
	return regexEOL.MatchString(text)
}

// IsFloat validate string is a float number
func IsFloat(val string) bool {
	return regexFloat.MatchString(val)
}

// IsInteger validate string is a integer
func IsInteger(val string) bool {
	return regexInteger.MatchString(val)
}

// IsUnsignedInteger validate string is a unsigned-integer
func IsUnsignedInteger(val string) bool {
	return regexUnsignedInteger.MatchString(val)
}

// RemoveEOL remove \r and \n
func RemoveEOL(text string) string {
	return regexEOL.ReplaceAllString(text, ``)
}
