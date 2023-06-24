package gerberos

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/BurntSushi/toml"
)

type Configuration struct {
	Verbose      bool
	Backend      string
	SaveFilePath string
	Rules        map[string]*Rule
}

func (c *Configuration) ReadFile(path string) error {
	cf, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open configuration file: %w", err)
	}
	defer cf.Close()

	return c.Read(cf)
}

func (c *Configuration) Read(r io.Reader) error {
	if _, err := toml.NewDecoder(r).Decode(&c); err != nil {
		var terr toml.ParseError
		if errors.As(err, &terr) {
			return fmt.Errorf("failed to decode configuration file: %s", terr.ErrorWithUsage())
		}
		return fmt.Errorf("failed to decode configuration file: %w", err)
	}

	return nil
}
