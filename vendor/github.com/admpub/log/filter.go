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
	hasCategory bool

	MaxLevel   Leveler          // the maximum severity level that is allowed
	Levels     map[Leveler]bool // 此属性被设置时，MaxLevel 无效
	Categories []string         // the allowed message categories. Categories can use "*" as a suffix for wildcard matching.
}

// Init initializes the filter.
// Init must be called before Allow is called.
func (t *Filter) Init() {
	t.catNames = make(map[string]bool)
	t.catPrefixes = make([]string, 0)
	for _, cat := range t.Categories {
		if strings.HasSuffix(cat, "*") {
			t.catPrefixes = append(t.catPrefixes, cat[:len(cat)-1])
		} else {
			t.catNames[cat] = true
		}
	}
	if t.Levels != nil {
		t.MaxLevel = Level(-1)
	}
	t.hasCategory = len(t.catNames) > 0 || len(t.catPrefixes) > 0
}

// Allow checks if a message meets the severity level and category requirements.
func (t *Filter) Allow(e *Entry) bool {
	if e == nil {
		return true
	}
	if t.MaxLevel.Int() > -1 {
		if e.Level.Int() > t.MaxLevel.Int() {
			return false
		}
	} else {
		if t.Levels == nil || !t.Levels[e.Level] {
			return false
		}
	}
	if !t.hasCategory {
		return true
	}
	if t.catNames[e.Category] {
		return true
	}
	for _, cat := range t.catPrefixes {
		if strings.HasPrefix(e.Category, cat) {
			return true
		}
	}
	return false
}

func (t *Filter) SetLevel(level interface{}) {
	if name, ok := level.(string); ok {
		if le, ok := GetLevel(name); ok {
			t.MaxLevel = le
		}
	} else if id, ok := level.(Leveler); ok {
		t.MaxLevel = id
	}
}

func (t *Filter) SetLevels(levels ...Leveler) {
	t.Levels = map[Leveler]bool{}
	for _, level := range levels {
		t.Levels[level] = true
	}
}
