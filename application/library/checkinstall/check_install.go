package checkinstall

import (
	"os/exec"

	"github.com/admpub/once"
)

func New(name string) *CheckInstall {
	return &CheckInstall{name: name}
}

type CheckInstall struct {
	name      string
	supported bool
	checkonce once.Once
}

func (c *CheckInstall) check() {
	_, err := exec.LookPath(c.name)
	c.supported = err == nil
}

func (c *CheckInstall) IsSupported() bool {
	c.checkonce.Do(c.check)
	return c.supported
}

func (c *CheckInstall) ResetCheck() {
	c.checkonce.Reset()
}
