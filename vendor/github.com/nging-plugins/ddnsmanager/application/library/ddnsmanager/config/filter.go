package config

import (
	"regexp"
	"strings"
)

type Filter struct {
	Include       string
	includeRegexp *regexp.Regexp
	Exclude       string
	excludeRegexp *regexp.Regexp
}

const RegexpPrefix = `regexp:`

func (f *Filter) Init() (err error) {
	if len(f.Include) > 0 && strings.HasPrefix(f.Include, RegexpPrefix) {
		f.includeRegexp, err = regexp.Compile(strings.TrimPrefix(f.Include, RegexpPrefix))
		if err != nil {
			return
		}
	}
	if len(f.Exclude) > 0 && strings.HasPrefix(f.Exclude, RegexpPrefix) {
		f.excludeRegexp, err = regexp.Compile(strings.TrimPrefix(f.Exclude, RegexpPrefix))
	}
	return
}

func (f *Filter) Match(ip string) bool {
	if len(f.Include) > 0 {
		if f.includeRegexp != nil {
			return f.includeRegexp.MatchString(ip)
		}
		return strings.Contains(ip, f.Include)
	}
	if len(f.Exclude) > 0 {
		if f.excludeRegexp != nil {
			return !f.excludeRegexp.MatchString(ip)
		}
		return !strings.Contains(ip, f.Exclude)
	}
	return true
}
