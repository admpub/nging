package request

import (
	"github.com/docker/docker/api/types"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

type ImageAdd struct {
	Ref          string `validate:"required"`
	All          bool
	RegistryAuth string // RegistryAuth is the base64 encoded credentials for the registry
	User         string
	Password     string
	Platform     string
}

func (c *ImageAdd) AfterValidate(ctx echo.Context) error {
	if len(c.RegistryAuth) == 0 && len(c.User) > 0 {
		c.RegistryAuth = com.Base64Encode(c.User + `:` + c.Password)
	}
	return nil
}

type ImageTag struct {
	Source string `validate:"required"`
	Target string `validate:"required"`
}

type ImageBuild struct {
	Tags        []string
	PullParent  bool
	Dockerfile  string
	BuildArgs   map[string]*string
	AuthConfigs map[string]types.AuthConfig
	Target      string
	Version     types.BuilderVersion
	Platform    string
}

func (c *ImageBuild) AfterValidate(ctx echo.Context) error {
	if len(c.Dockerfile) == 0 {
		c.Dockerfile = `Dockerfile`
	}
	return nil
}

type ImageImport struct {
	Ref        string `validate:"required"`
	SourceName string
	types.ImageImportOptions
}

func (c *ImageImport) AfterValidate(ctx echo.Context) error {
	return nil
}
