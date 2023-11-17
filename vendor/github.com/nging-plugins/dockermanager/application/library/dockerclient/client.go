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
