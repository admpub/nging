// Copyright 2015 Qiang Xue. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/admpub/color"
)

var colorBrushes = map[Leveler]*color.Color{
	LevelDebug:    color.New(color.FgHiCyan),    // cyan
	LevelProgress: color.New(color.FgHiBlack),   // black
	LevelInfo:     color.New(color.FgHiWhite),   // white
	LevelOkay:     color.New(color.FgHiGreen),   // green
	LevelWarn:     color.New(color.FgHiYellow),  // yellow
	LevelError:    color.New(color.FgHiRed),     // red
	LevelFatal:    color.New(color.FgHiMagenta), // magenta
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
		ColorMode: DefaultConsoleColorize,
		ColorType: ColorFlag,
		Writer:    os.Stdout,
		close:     make(chan bool),
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
				fmt.Fprintln(t.Writer, t.ColorizeFlag(e))
			}
		default:
			t.outputFunc = func(t *ConsoleTarget, e *Entry) {
				fmt.Fprintln(t.Writer, t.ColorizeRow(e))
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

func (t *ConsoleTarget) ColorizeFlag(e *Entry) string {
	s := e.Level.Tag()
	cs := e.Level.Color()
	if cs != nil {
		return cs.SprintFunc()(s) + e.String()
	}
	return s + e.String()
}

func (t *ConsoleTarget) ColorizeRow(e *Entry) string {
	cs := e.Level.Color()
	if cs != nil {
		return cs.SprintFunc()(e.String())
	}
	return e.String()
}

// Close closes the console target.
func (t *ConsoleTarget) Close() {
	<-t.close
}
