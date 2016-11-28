// Copyright 2015 Qiang Xue. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"errors"
	"fmt"
	"io"
	"os"

	ct "github.com/admpub/go-colortext"
)

type colorSetting struct {
	Color  ct.Color
	Bright bool
}

var colorBrushes = map[Level]colorSetting{
	LevelDebug: colorSetting{ct.Cyan, true},    // cyan
	LevelInfo:  colorSetting{ct.Green, true},   // green
	LevelWarn:  colorSetting{ct.Yellow, true},  // yellow
	LevelError: colorSetting{ct.Red, true},     // red
	LevelFatal: colorSetting{ct.Magenta, true}, // magenta
}

// ConsoleTarget writes filtered log messages to console window.
type ConsoleTarget struct {
	*Filter
	ColorMode bool      // whether to use colors to differentiate log levels
	Writer    io.Writer // the writer to write log messages
	close     chan bool
}

// NewConsoleTarget creates a ConsoleTarget.
// The new ConsoleTarget takes these default options:
// MaxLevel: LevelDebug, ColorMode: true, Writer: os.Stdout
func NewConsoleTarget() *ConsoleTarget {
	return &ConsoleTarget{
		Filter:    &Filter{MaxLevel: LevelDebug},
		ColorMode: true,
		Writer:    os.Stdout,
		close:     make(chan bool, 0),
	}
}

// Open prepares ConsoleTarget for processing log messages.
func (t *ConsoleTarget) Open(io.Writer) error {
	t.Filter.Init()
	if t.Writer == nil {
		return errors.New("ConsoleTarget.Writer cannot be nil")
	}
	return nil
}

// Process writes a log message using Writer.
func (t *ConsoleTarget) Process(e *Entry) {
	if e == nil {
		t.close <- true
		return
	}
	if !t.Allow(e) {
		return
	}
	msg := e.String()
	if t.ColorMode {
		if t.Colorize(e.Level) {
			defer ct.ResetColor()
		}
	}
	fmt.Fprintln(t.Writer, msg)
}

func (t *ConsoleTarget) Colorize(level Level) bool {
	cs, ok := colorBrushes[level]
	if ok {
		ct.Foreground(cs.Color, cs.Bright)
	}
	return ok
}

// Close closes the console target.
func (t *ConsoleTarget) Close() {
	<-t.close
}
