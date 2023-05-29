package fileperm

import (
	"path/filepath"
	"regexp"
	"strings"
)

func NewRule(path string, readable bool, writeable *bool, regexpress bool) *Rule {
	r := &Rule{
		Regexpress: regexpress,
		Writeable:  writeable,
		Readable:   readable,
	}
	r.SetPath(path)
	return r
}

func NewRuleAndInit(path string, readable bool, writeable *bool, regexpress bool) (*Rule, error) {
	r := &Rule{
		Regexpress: regexpress,
		Writeable:  writeable,
		Readable:   readable,
	}
	r.SetPath(path)
	return r, r.Init()
}

type Rule struct {
	Regexpress bool
	Writeable  *bool
	Readable   bool
	Path       string
	regexp     *regexp.Regexp
}

func (r *Rule) Init() (err error) {
	if r.Regexpress && r.regexp == nil {
		r.regexp, err = regexp.Compile(r.Path)
	}
	return
}

func (r *Rule) SetPath(path string) *Rule {
	r.Path = filepath.ToSlash(path)
	return r
}

func (r *Rule) detectRuleType() (err error) {
	if strings.Contains(r.Path, `*`) {
		r.Regexpress = true
		r.regexp, err = regexp.Compile(strings.ReplaceAll(r.Path, `*`, `(.*)`))
	} else if strings.Contains(r.Path, `|`) {
		r.Regexpress = true
		r.regexp, err = regexp.Compile(r.Path)
	}
	return
}

func (r *Rule) SetWriteable(on bool) *Rule {
	r.Writeable = &on
	return r
}

func (r *Rule) FixedPathPrefixAndSuffix(hasPrefix, hasSuffix bool) string {
	rulePath := r.Path
	if hasPrefix {
		if !strings.HasPrefix(rulePath, `/`) {
			rulePath = `/` + rulePath
		}
	} else {
		rulePath = strings.TrimPrefix(rulePath, `/`)
	}
	if hasSuffix {
		if !strings.HasSuffix(rulePath, `/`) {
			rulePath += `/`
		}
	} else {
		rulePath = strings.TrimSuffix(rulePath, `/`)
	}
	return rulePath
}

type User struct {
	Writeable bool   // 是否默认可写
	RootDir   string // 跟路径
	Rules     Rules
}

// Allowed checks if the user has permission to access a directory/file
func (u User) Allowed(path string, modification bool) bool {
	var rule *Rule
	i := len(u.Rules) - 1

	path = filepath.ToSlash(path)
	hasSuffix := strings.HasSuffix(path, `/`)
	hasPrefix := strings.HasPrefix(path, `/`)

	for i >= 0 {
		rule = u.Rules[i]

		isAllowed := rule.Readable
		if modification {
			if rule.Writeable != nil {
				isAllowed = *rule.Writeable
			} else {
				isAllowed = u.Writeable
			}
		}
		if rule.Regexpress {
			if rule.regexp.MatchString(path) {
				return isAllowed
			}
		} else {
			rulePath := rule.FixedPathPrefixAndSuffix(hasPrefix, hasSuffix)
			if strings.HasPrefix(path, rulePath) {
				return isAllowed
			}
		}

		i--
	}

	return !modification || u.Writeable
}
