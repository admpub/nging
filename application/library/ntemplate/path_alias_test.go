package ntemplate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathAliases(t *testing.T) {
	pathAliases := PathAliases{}
	pathAliases.Add(`test`, `/a/b/c/d/`)
	assert.Equal(t, `/a/b/c/d/`, pathAliases.aliases[`test`][0])

	withAliasTagPath := `[test]user/index`
	fullpath := pathAliases.Parse(withAliasTagPath)
	assert.Equal(t, `/a/b/c/d/user/index`, fullpath)
	assert.Equal(t, withAliasTagPath, pathAliases.Restore(fullpath))

	withAliasPrefixPath := `test/user/index`
	fullpath = pathAliases.ParsePrefix(withAliasPrefixPath)
	assert.Equal(t, `/a/b/c/d/test/user/index`, fullpath)
	assert.Equal(t, withAliasPrefixPath, pathAliases.RestorePrefix(fullpath))

}
