/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package bindata

import (
	"errors"
	"net/http"
	"os"
	"strings"

	assetfs "github.com/admpub/go-bindata-assetfs"
)

var ErrUnsupported = errors.New(`unsupported bindata`)

var (
	Asset = func(name string) ([]byte, error) {
		return nil, ErrUnsupported
	}

	AssetDir = func(name string) ([]string, error) {
		return nil, ErrUnsupported
	}

	AssetInfo = func(name string) (os.FileInfo, error) {
		return nil, ErrUnsupported
	}
)

type staticAsset struct {
	prefix string
	*assetfs.AssetFS
}

func (s *staticAsset) Open(name string) (http.File, error) {
	name = strings.TrimPrefix(name, s.prefix)
	return s.AssetFS.Open(name)
}

func NewStaticAssetFS(prefix string, afs *assetfs.AssetFS) http.FileSystem {
	if len(prefix) == 0 {
		return afs
	}
	return &staticAsset{
		prefix:  prefix,
		AssetFS: afs,
	}
}
