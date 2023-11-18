/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

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

package dockerclient

import (
	"bytes"
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/dbschema"
)

func PullImage(ctx echo.Context, user *dbschema.NgingUser, ref string, c *client.Client, options *types.ImagePullOptions) error {
	return StartBackgroundRun(ctx, user.Username, `dockerImagePull`, ref, func(ctx context.Context) (io.ReadCloser, error) {
		var err error
		if c == nil {
			c, err = Client()
			if err != nil {
				return nil, err
			}
		}
		if options == nil {
			options = &types.ImagePullOptions{}
		}
		return c.ImagePull(ctx, ref, *options)
	})
}

func BuildImage(ctx echo.Context, user *dbschema.NgingUser, dockerfileContent []byte, c *client.Client, options types.ImageBuildOptions) error {
	optionJSONb, _ := com.JSONEncode(options)
	keyBytes := append([]byte{}, dockerfileContent...)
	keyBytes = append(keyBytes, optionJSONb...)
	bgKey := com.ByteMd5(keyBytes)
	return StartBackgroundRun(ctx, user.Username, `dockerImageBuild`, bgKey, func(ctx context.Context) (io.ReadCloser, error) {
		var err error
		if c == nil {
			c, err = Client()
			if err != nil {
				return nil, err
			}
		}
		var r io.Reader
		if len(dockerfileContent) > 0 {
			r = bytes.NewReader(dockerfileContent)
		}
		result, err := c.ImageBuild(ctx, r, options)
		if err != nil {
			return nil, err
		}
		return result.Body, err
	})
}
