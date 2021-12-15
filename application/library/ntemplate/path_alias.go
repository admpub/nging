package ntemplate

import (
	"path/filepath"
	"strings"
)

type PathAliases map[string]string

func (p PathAliases) Add(alias, absPath string) PathAliases {
	p[alias] = absPath
	return p
}

func (p PathAliases) ParsePrefix(withAliasPrefixPath string) string {
	if len(withAliasPrefixPath) < 3 || withAliasPrefixPath[0] == '/' || withAliasPrefixPath[0] == '.' {
		return withAliasPrefixPath
	}
	parts := strings.SplitN(withAliasPrefixPath, `/`, 2)
	if len(parts) != 2 {
		return withAliasPrefixPath
	}
	alias := parts[0]
	if opath, ok := p[alias]; ok {
		return filepath.Join(opath, withAliasPrefixPath)
	}
	return withAliasPrefixPath
}

func (p PathAliases) RestorePrefix(fullpath string) string {
	for _, absPath := range p {
		if strings.HasPrefix(fullpath, absPath) {
			return fullpath[len(absPath):]
		}
	}
	return fullpath
}

func (p PathAliases) Parse(withAliasTagPath string) string {
	if len(withAliasTagPath) < 3 || withAliasTagPath[0] != '[' {
		return withAliasTagPath
	}
	withAliasTagPath = withAliasTagPath[1:]
	parts := strings.SplitN(withAliasTagPath, `]`, 2)
	if len(parts) != 2 {
		return withAliasTagPath
	}
	alias := parts[0]
	rpath := parts[1]
	if opath, ok := p[alias]; ok {
		rpath = filepath.Join(opath, rpath)
	}
	return rpath
}

func (p PathAliases) Restore(fullpath string) string {
	for alias, absPath := range p {
		if strings.HasPrefix(fullpath, absPath) {
			return `[` + alias + `]` + fullpath[len(absPath):]
		}
	}
	return fullpath
}
