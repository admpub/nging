// Copyright 2015 Qiang Xue. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"strings"
)

// Filter checks if a log message meets the level and category requirements.
type Filter struct {
	catNames    map[string]bool
	catPrefixes []string

	MaxLevel   Level          // the maximum severity level that is allowed
	Levels     map[Level]bool // 此属性被设置时，MaxLevel 无效
	Categories []string       // the allowed message categories. Categories can use "*" as a suffix for wildcard matching.
}

// Init initializes the filter.
// Init must be called before Allow is called.
func (t *Filter) Init() {
	t.catNames = make(map[string]bool, 0)
	t.catPrefixes = make([]string, 0)
	for _, cat := range t.Categories {
		if strings.HasSuffix(cat, "*") {
			t.catPrefixes = append(t.catPrefixes, cat[:len(cat)-1])
		} else {
			t.catNames[cat] = true
		}
	}
	if t.Levels != nil {
		t.MaxLevel = -1
	}
}

// Allow checks if a message meets the severity level and category requirements.
func (t *Filter) Allow(e *Entry) bool {
	if e == nil {
		return true
	}
	if t.MaxLevel > -1 {
		if e.Level > t.MaxLevel {
			return false
		}
	} else {
		if t.Levels == nil || !t.Levels[e.Level] {
			return false
		}
	}
	if t.catNames[e.Category] {
		return true
	}
	for _, cat := range t.catPrefixes {
		if strings.HasPrefix(e.Category, cat) {
			return true
		}
	}
	return len(t.catNames) == 0 && len(t.catPrefixes) == 0
}

func (t *Filter) SetLevel(level interface{}) {
	if name, ok := level.(string); ok {
		if le, ok := GetLevel(name); ok {
			t.MaxLevel = le
		}
	} else if id, ok := level.(Level); ok {
		t.MaxLevel = id
	}
}

func (t *Filter) SetLevels(levels ...Level) {
	t.Levels = map[Level]bool{}
	for _, level := range levels {
		t.Levels[level] = true
	}
}
