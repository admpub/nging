package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlob(t *testing.T) {
	files, err := filepath.Glob(filepath.Join(os.Getenv(`GOPATH`), `src/github.com/admpub/webx/config/install.*.sql`))
	if err != nil {
		panic(err)
	}

	assert.True(t, len(files) > 0)
	//assert.Equal(t, []string{}, files)
}
