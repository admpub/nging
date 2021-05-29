package filter

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var paramReplacementPattern = regexp.MustCompile("\\{[a-zA-Z0-9_\\-.]+}")

type ruleReplaceAction struct {
	request        *http.Request
	responseHeader *http.Header
	searchPattern  *regexp.Regexp
	replacement    []byte
}

func (instance *ruleReplaceAction) replacer(input []byte) []byte {
	pattern := instance.searchPattern
	if pattern == nil {
		return input
	}
	rawReplacement := instance.replacement
	if len(rawReplacement) <= 0 {
		return []byte{}
	}
	groups := pattern.FindSubmatch(input)
	replacement := paramReplacementPattern.ReplaceAllFunc(rawReplacement, func(input2 []byte) []byte {
		return instance.paramReplacer(input2, groups)
	})
	return replacement
}

func (instance *ruleReplaceAction) paramReplacer(input []byte, groups [][]byte) []byte {
	if len(input) < 3 {
		return input
	}
	name := string(input[1 : len(input)-1])
	if index, err := strconv.Atoi(name); err == nil {
		if index >= 0 && index < len(groups) {
			return groups[index]
		}
		return input
	}

	if value, ok := instance.contextValueBy(name); ok {
		return []byte(value)
	}
	return input
}

func (instance *ruleReplaceAction) contextValueBy(name string) (string, bool) {
	if strings.HasPrefix(name, "request_") {
		return instance.contextRequestValueBy(name[8:])
	}
	if strings.HasPrefix(name, "response_") {
		return instance.contextResponseValueBy(name[9:])
	}
	if strings.HasPrefix(name, "env_") {
		return instance.contextEnvironmentValueBy(name[4:])
	}
	if name == "now" {
		return instance.contextNowValueBy("")
	}
	if strings.HasPrefix(name, "now:") {
		return instance.contextNowValueBy(name[4:])
	}
	return "", false
}

func (instance *ruleReplaceAction) contextRequestValueBy(name string) (string, bool) {
	request := instance.request
	if strings.HasPrefix(name, "header_") {
		return request.Header.Get(name[7:]), true
	}
	switch name {
	case "url":
		return request.URL.String(), true
	case "path":
		return request.URL.Path, true
	case "method":
		return request.Method, true
	case "host":
		return request.Host, true
	case "proto":
		return request.Proto, true
	case "remoteAddress":
		return request.RemoteAddr, true
	}
	return "", false
}

func (instance *ruleReplaceAction) contextResponseValueBy(name string) (string, bool) {
	if name == "header_last_modified" || name == "header_last-modified" {
		return instance.contextLastModifiedValueBy("")
	}
	if strings.HasPrefix(name, "header_last_modified:") || strings.HasPrefix(name, "header_last-modified:") {
		return instance.contextLastModifiedValueBy(name[21:])
	}
	if strings.HasPrefix(name, "header_") {
		return (*instance.responseHeader).Get(name[7:]), true
	}
	return "", false
}

func (instance *ruleReplaceAction) contextEnvironmentValueBy(name string) (string, bool) {
	return os.Getenv(name), true
}

func (instance *ruleReplaceAction) contextNowValueBy(pattern string) (string, bool) {
	return instance.formatTimeBy(time.Now(), pattern), true
}

func (instance *ruleReplaceAction) contextLastModifiedValueBy(pattern string) (string, bool) {
	plain := instance.responseHeader.Get("last-Modified")
	if plain == "" {
		// Fallback to now
		return instance.contextNowValueBy(pattern)
	}
	t, err := time.Parse(time.RFC1123, plain)
	if err != nil {
		log.Printf("[WARN] Serving illegal 'Last-Modified' header value '%v' for '%v': Got: %v", plain, instance.request.URL, err)
		// Fallback to now
		return instance.contextNowValueBy(pattern)
	}
	return instance.formatTimeBy(t, pattern), true
}

func (instance *ruleReplaceAction) formatTimeBy(t time.Time, pattern string) string {
	if pattern == "" || pattern == "RFC" || pattern == "RFC3339" {
		return t.Format(time.RFC3339)
	}
	if pattern == "unix" {
		return fmt.Sprintf("%d", t.Unix())
	}
	if pattern == "timestamp" {
		stamp := t.Unix() * 1000
		stamp += int64(t.Nanosecond()) / int64(1000000)
		return fmt.Sprintf("%d", stamp)
	}
	return t.Format(pattern)
}
