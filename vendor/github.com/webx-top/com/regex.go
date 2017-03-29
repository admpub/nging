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
	regex_email_pattern        = `(?i)[A-Z0-9._%+-]+@(?:[A-Z0-9-]+\.)+[A-Z]{2,6}`
	regex_strict_email_pattern = `(?i)[A-Z0-9!#$%&'*+/=?^_{|}~-]+` +
		`(?:\.[A-Z0-9!#$%&'*+/=?^_{|}~-]+)*` +
		`@(?:[A-Z0-9](?:[A-Z0-9-]*[A-Z0-9])?\.)+` +
		`[A-Z0-9](?:[A-Z0-9-]*[A-Z0-9])?`
	regex_url_pattern      = `(ftp|http|https):\/\/(\w+:{0,1}\w*@)?(\S+)(:[0-9]+)?(\/|\/([\w#!:.?+=&%@!\-\/]))?`
	regex_username_pattern = `^[\w\p{Han}]+$`
	regex_eol_pattern      = "[\r\n]+"
)

var (
	regex_email        *regexp.Regexp
	regex_strict_email *regexp.Regexp
	regex_url          *regexp.Regexp
	regex_username     *regexp.Regexp
	regex_eol          *regexp.Regexp
)

func init() {
	regex_email = regexp.MustCompile(regex_email_pattern)
	regex_strict_email = regexp.MustCompile(regex_strict_email_pattern)
	regex_url = regexp.MustCompile(regex_url_pattern)
	regex_username = regexp.MustCompile(regex_username_pattern)
	regex_eol = regexp.MustCompile(regex_eol_pattern)
}

// IsEmail validate string is an email address, if not return false
// basically validation can match 99% cases
func IsEmail(email string) bool {
	return regex_email.MatchString(email)
}

// IsEmailRFC validate string is an email address, if not return false
// this validation omits RFC 2822
func IsEmailRFC(email string) bool {
	return regex_strict_email.MatchString(email)
}

// IsURL validate string is a url link, if not return false
// simple validation can match 99% cases
func IsURL(url string) bool {
	return regex_url.MatchString(url)
}

// IsUsername validate string is a available username
func IsUsername(username string) bool {
	return regex_username.MatchString(username)
}

// IsSingleLineText validate string is a single-line text
func IsSingleLineText(text string) bool {
	return !regex_eol.MatchString(text)
}

// IsMultiLineText validate string is a multi-line text
func IsMultiLineText(text string) bool {
	return regex_eol.MatchString(text)
}

// RemoveEOL remove \r and \n
func RemoveEOL(text string) string {
	return regex_eol.ReplaceAllString(text, ``)
}
