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
	"os"
	"os/user"
	"path/filepath"

	"github.com/admpub/log"
	"github.com/admpub/once"
	"github.com/docker/docker/client"
	"github.com/webx-top/com"
)

var SearchableDockerSockFiles = []string{`/var/run/docker.sock`}
var sockFileCheckOnce = once.Once{}

func sockFileCheck() {
	dockerHost := os.Getenv("DOCKER_HOST")
	if len(dockerHost) != 0 {
		return
	}
	for _, sockFile := range SearchableDockerSockFiles {
		sockFilePath, err := os.Readlink(sockFile)
		if err != nil {
			log.Error(err)
			continue
		}
		if com.FileExists(sockFilePath) {
			os.Setenv("DOCKER_HOST", "unix://"+sockFilePath)
			return
		}
	}
	u, err := user.Current()
	if err == nil {
		userSockFilePath := filepath.Join(u.HomeDir, `.docker/run/docker.sock`)
		sockFilePath, err := os.Readlink(userSockFilePath)
		if err != nil {
			log.Error(err)
		} else {
			if com.FileExists(sockFilePath) {
				os.Setenv("DOCKER_HOST", "unix://"+sockFilePath)
				return
			}
		}
	}
}

func Client() (*client.Client, error) {
	sockFileCheckOnce.Do(sockFileCheck)
	return client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
}
