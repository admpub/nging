package local

import (
	"testing"
	"context"

	"github.com/stretchr/testify/assert"
)

func TestOne(t *testing.T) {
	ctx := context.Background()
	f := NewFilesystem(ctx, `test`, `https://img.admpub.com/`)
	viewURL := `https://img.admpub.com`+f.URLDir(`user/1/2020/1/2/a.jpg`)
	expected := `https://img.admpub.com/public/upload/test/user/1/2020/1/2/a.jpg`
	assert.Equal(t, `https://img.admpub.com`, f.BaseURL())
	assert.Equal(t, expected, viewURL)
	assert.Equal(t, `user/1/2020/1/2/a.jpg`, f.URLToFile(viewURL))
	expectedPath := `/public/upload/test/user/1/2020/1/2/a.jpg`
	assert.Equal(t, expectedPath, f.URLToPath(viewURL))
	assert.Equal(t, `user/1/2020/1/2/a.jpg`, f.URLToFile(expectedPath))
}