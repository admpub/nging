package ninetail

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"strings"

	"github.com/admpub/tail"
	runewidth "github.com/mattn/go-runewidth"
)

var (
	// red, green, yellow, magenta, cyan
	ansiColorCodes    = [...]int{31, 32, 33, 35, 36}
	seekInfoOnStart   = &tail.SeekInfo{Offset: 0, Whence: os.SEEK_END}
	TailDefaultConfig = &tail.Config{
		Follow:   true,
		Location: seekInfoOnStart,
		Logger:   tail.DiscardingLogger,
	}
)

//Tailer contains watches tailed files and contains per-file output parameters
type Tailer struct {
	*tail.Tail
	colorCode int
	padding   string
}

//NewTailers creates slice of Tailers from file names.
//Colors of file names are cycled through the list.
//maxWidth is a maximum widht of passed file names, for nice alignment
func NewTailers(filenames []string, tailConfigs ...*tail.Config) ([]*Tailer, error) {
	maxLength := maximumNameLength(filenames)
	ts := make([]*Tailer, len(filenames))
	for i, filename := range filenames {
		t, err := newTailer(filename, getColorCode(i), maxLength, tailConfigs...)
		if err != nil {
			return nil, err
		}

		ts[i] = t
	}

	return ts, nil
}

func newTailer(filename string, colorCode int, maxWidth int, tailConfigs ...*tail.Config) (*Tailer, error) {
	var tailConfig *tail.Config
	if len(tailConfigs) > 0 {
		tailConfig = tailConfigs[0]
	}
	if tailConfig == nil {
		tailConfig = TailDefaultConfig
	} else {
		if tailConfig.Location == nil {
			tailConfig.Location = TailDefaultConfig.Location
		}
		if tailConfig.Logger == nil {
			tailConfig.Logger = TailDefaultConfig.Logger
		}
	}
	t, err := tail.TailFile(filename, *tailConfig)

	if err != nil {
		return nil, err
	}

	dispNameLength := displayFilenameLength(filename)

	return &Tailer{
		Tail:      t,
		colorCode: colorCode,
		padding:   strings.Repeat(" ", maxWidth-dispNameLength),
	}, nil
}

//Do formats, colors and writes to stdout appended lines when they happen, exiting on write error
func (t Tailer) Do(output io.Writer) {
	for line := range t.Lines {
		_, err := fmt.Fprintf(
			output,
			"\x1b[%dm%s%s\x1b[0m: %s\n",
			t.colorCode,
			t.padding,
			t.name(),
			line.Text,
		)
		if err != nil {
			return
		}
	}
}

func (t Tailer) name() string {
	return filepath.Base(t.Filename)
}

func getColorCode(index int) int {
	return ansiColorCodes[index%len(ansiColorCodes)]
}

func maximumNameLength(filenames []string) int {
	max := 0
	for _, name := range filenames {
		if current := displayFilenameLength(name); current > max {
			max = current
		}
	}
	return max
}

func displayFilename(filename string) string {
	return filepath.Base(filename)
}

func displayFilenameLength(filename string) int {
	return runewidth.StringWidth(displayFilename(filename))
}
