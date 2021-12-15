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
	rpath, _ := p.ParsePrefixOk(withAliasPrefixPath)
	return rpath
}

func (p PathAliases) ParsePrefixOk(withAliasPrefixPath string) (string, bool) {
	if len(withAliasPrefixPath) < 3 || withAliasPrefixPath[0] == '/' || withAliasPrefixPath[0] == '.' {
		return withAliasPrefixPath, false
	}
	parts := strings.SplitN(withAliasPrefixPath, `/`, 2)
	if len(parts) != 2 {
		return withAliasPrefixPath, false
	}
	alias := parts[0]
	if opath, ok := p[alias]; ok {
		return filepath.Join(opath, withAliasPrefixPath), true
	}
	return withAliasPrefixPath, false
}

func (p PathAliases) RestorePrefix(fullpath string) string {
	rpath, _ := p.RestorePrefixOk(fullpath)
	return rpath
}

func (p PathAliases) RestorePrefixOk(fullpath string) (string, bool) {
	for _, absPath := range p {
		if strings.HasPrefix(fullpath, absPath) {
			return fullpath[len(absPath):], true
		}
	}
	return fullpath, false
}

func (p PathAliases) Parse(withAliasTagPath string) string {
	rpath, _ := p.ParseOk(withAliasTagPath)
	return rpath
}

func (p PathAliases) ParseOk(withAliasTagPath string) (string, bool) {
	if len(withAliasTagPath) < 3 || withAliasTagPath[0] != '[' {
		return withAliasTagPath, false
	}
	withAliasTagPath = withAliasTagPath[1:]
	parts := strings.SplitN(withAliasTagPath, `]`, 2)
	if len(parts) != 2 {
		return withAliasTagPath, false
	}
	alias := parts[0]
	rpath := parts[1]
	if opath, ok := p[alias]; ok {
		rpath = filepath.Join(opath, rpath)
		return rpath, true
	}
	return rpath, false
}

func (p PathAliases) Restore(fullpath string) string {
	rpath, _ := p.RestoreOk(fullpath)
	return rpath
}

func (p PathAliases) RestoreOk(fullpath string) (string, bool) {
	for alias, absPath := range p {
		if strings.HasPrefix(fullpath, absPath) {
			return `[` + alias + `]` + fullpath[len(absPath):], true
		}
	}
	return fullpath, false
}
