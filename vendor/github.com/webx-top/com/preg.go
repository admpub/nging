package com

import (
	"regexp"
)

func PregReplace(expr string, repl string, src string) string {
	re := regexp.MustCompile(expr)
	return re.ReplaceAllString(src, repl)
}

func PregReplaceByte(expr string, repl []byte, src []byte) []byte {
	re := regexp.MustCompile(expr)
	return re.ReplaceAll(src, repl)
}

func PregReplaceCallback(expr string, repl func(string) string, src string) string {
	re := regexp.MustCompile(expr)
	return re.ReplaceAllStringFunc(src, repl)
}

func PregReplaceByteCallback(expr string, repl func([]byte) []byte, src []byte) []byte {
	re := regexp.MustCompile(expr)
	return re.ReplaceAllFunc(src, repl)
}

func PregSplit(expr string, src string, n int) []string {
	re := regexp.MustCompile(expr)
	return re.Split(src, n)
}

func PregMatch(expr string, src string) string {
	re := regexp.MustCompile(expr)
	return re.FindString(src)
}

func PregIsMatch(expr string, src []byte) (hasMatched bool) {
	hasMatched, _ = regexp.Match(expr, src)
	return
}
func PregIsMatchString(expr string, src string) (hasMatched bool) {
	hasMatched, _ = regexp.MatchString(expr, src)
	return
}

func PregMatchAll(expr string, src string, n int) [][]string {
	re := regexp.MustCompile(expr)
	return re.FindAllStringSubmatch(src, n)
}

func PregMatchAll2(expr string, src string, n int) []string {
	re := regexp.MustCompile(expr)
	return re.FindAllString(src, n)
}
