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

const (
	ColorFlag = iota
	ColorRow
)

// ConsoleTarget writes filtered log messages to console window.
type ConsoleTarget struct {
	*Filter
	ColorMode  bool // whether to use colors to differentiate log levels
	ColorType  int
	Writer     io.Writer // the writer to write log messages
	close      chan bool
	outputFunc func(*ConsoleTarget, *Entry)
}

// NewConsoleTarget creates a ConsoleTarget.
// The new ConsoleTarget takes these default options:
// MaxLevel: LevelDebug, ColorMode: true, Writer: os.Stdout
func NewConsoleTarget() *ConsoleTarget {
	return &ConsoleTarget{
		Filter:    &Filter{MaxLevel: LevelDebug},
		ColorMode: true,
		ColorType: ColorFlag,
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
	if t.ColorMode {
		switch t.ColorType {
		case ColorFlag:
			t.outputFunc = func(t *ConsoleTarget, e *Entry) {
				t.ColorizeFlag(e.Level)
				fmt.Fprintln(t.Writer, e.String())
			}
		default:
			t.outputFunc = func(t *ConsoleTarget, e *Entry) {
				t.ColorizeRow(e.Level)
				fmt.Fprintln(t.Writer, e.String())
				ct.ResetColor()
			}
		}
	} else {
		t.outputFunc = func(t *ConsoleTarget, e *Entry) {
			fmt.Fprintln(t.Writer, e.String())
		}
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
	t.outputFunc(t, e)
}

func (t *ConsoleTarget) ColorizeFlag(level Level) bool {
	cs, ok := colorBrushes[level]
	if ok {
		ct.Foreground(cs.Color, cs.Bright)
		fmt.Fprint(t.Writer, `[`+level.String()[0:1]+`]`)
		ct.ResetColor()
	}
	return ok
}

func (t *ConsoleTarget) ColorizeRow(level Level) bool {
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
