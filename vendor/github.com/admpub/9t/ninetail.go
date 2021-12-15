package ninetail

import (
	"io"
	"os"
	"sync"

	"github.com/admpub/tail"

	colorable "github.com/mattn/go-colorable"
)

type NineTail struct {
	output  io.Writer
	tailers []*Tailer
}

type Config struct {
	Colorize bool
	Writer   io.Writer
	*tail.Config
}

func Runner(filenames []string, config Config) (*NineTail, error) {
	var output io.Writer
	if config.Writer == nil {
		if config.Colorize {
			output = colorable.NewColorableStdout()
		} else {
			output = colorable.NewNonColorable(os.Stdout)
		}
	} else {
		output = config.Writer
	}

	tailers, err := NewTailers(filenames, config.Config)
	if err != nil {
		return nil, err
	}

	return &NineTail{
		output:  output,
		tailers: tailers,
	}, nil
}

func (n *NineTail) Run() {
	var wg sync.WaitGroup

	for _, t := range n.tailers {
		wg.Add(1)
		go func(t *Tailer) {
			t.Do(n.output)
			wg.Done()
		}(t)
	}

	wg.Wait()
}

func (n *NineTail) Stop() {
	var wg sync.WaitGroup

	for _, t := range n.tailers {
		wg.Add(1)
		go func(t *Tailer) {
			t.Stop()
			wg.Done()
		}(t)
	}

	wg.Wait()
}
