package engine

import (
	"strings"

	"github.com/admpub/nging/v5/application/library/checkinstall"
	"github.com/nging-plugins/caddymanager/application/dbschema"
	"github.com/webx-top/echo"
)

type CertPathFormat struct {
	Cert  string
	Key   string
	Trust string
}

type CertPathFormatWithUpdater struct {
	CertPathFormat
	LocalUpdater     string // certbot / lego
	ContainerUpdater string // certbot / lego
	SaveDir          string
}

func (c *CertPathFormatWithUpdater) CopyFrom(m *dbschema.NgingVhostServer) {
	c.CertPathFormat.Cert = m.CertPathFormatCert
	c.CertPathFormat.Key = m.CertPathFormatKey
	c.CertPathFormat.Trust = m.CertPathFormatTrust
	if len(c.CertPathFormat.Cert) > 0 {
		parts := strings.SplitN(c.CertPathFormat.Cert, `.lego`, 2)
		if len(parts) == 2 {
			c.SaveDir = parts[0]
			if !strings.HasPrefix(m.CertContainerDir, `{lego}`) && !strings.HasPrefix(m.CertContainerDir, `{lego:`) {
				m.CertContainerDir = `{lego}` + m.CertContainerDir
			}
			if !strings.HasPrefix(m.CertLocalDir, `{lego}`) && !strings.HasPrefix(m.CertLocalDir, `{lego:`) {
				m.CertContainerDir = `{lego}` + m.CertLocalDir
			}
		} else {
			parts = strings.SplitN(c.CertPathFormat.Cert, `live`, 2)
			if len(parts) == 2 {
				c.SaveDir = parts[0]
				if !strings.HasPrefix(m.CertContainerDir, `{certbot}`) && !strings.HasPrefix(m.CertContainerDir, `{certbot:`) {
					m.CertContainerDir = `{certbot}` + m.CertContainerDir
				}
				if !strings.HasPrefix(m.CertLocalDir, `{certbot}`) && !strings.HasPrefix(m.CertLocalDir, `{certbot:`) {
					m.CertContainerDir = `{certbot}` + m.CertLocalDir
				}
			}
		}
	}
	if len(m.CertLocalDir) > 2 { // {certbot}/etc/nginx/... or {certbot:cloudflare}/etc/nginx/...
		if strings.HasPrefix(m.CertLocalDir, `{`) {
			parts := strings.SplitN(m.CertLocalDir[1:], `}`, 2)
			if len(parts) == 2 {
				c.LocalUpdater = strings.TrimSpace(parts[0])
				m.CertLocalDir = strings.TrimSpace(parts[1])
			}
		}
	}
	if len(m.CertContainerDir) > 2 { // {certbot}/etc/nginx/... or {certbot:cloudflare}/etc/nginx/...
		if strings.HasPrefix(m.CertContainerDir, `{`) {
			parts := strings.SplitN(m.CertContainerDir[1:], `}`, 2)
			if len(parts) == 2 {
				c.ContainerUpdater = strings.TrimSpace(parts[0])
				m.CertContainerDir = strings.TrimSpace(parts[1])
			}
		}
	}
}

func (c *CertPathFormatWithUpdater) AutoDetect(ctx echo.Context) {
	// 此函数任务：需要确认可以使用的更新程序和证书路径模板
	if len(c.CertLocalUpdater()) == 0 ||
		(len(c.Cert) == 0 && len(c.Key) == 0 && len(c.Trust) == 0) { // 如果用户没有设置模板
		for _, v := range CertUpdaters.Slice() {
			installed, ok := ctx.Internal().Get(`installed.` + v.K).(bool)
			if !ok {
				installed = checkinstall.DefaultChecker(v.K)
				ctx.Internal().Set(`installed.`+v.K, installed)
			}
			if installed {
				c.CertPathFormat = v.X.(CertPathFormat)
				if len(c.CertLocalUpdater()) > 0 {
					parts := strings.SplitN(c.CertLocalUpdater(), `:`, 2)
					if len(parts) == 2 {
						c.SetCertLocalUpdater(v.K + `:` + parts[1])
						break
					}
				}
				c.SetCertLocalUpdater(v.K)
				break
			}
		}
	}
}

func (c *CertPathFormatWithUpdater) CertLocalUpdater() string {
	return c.LocalUpdater
}

func (c *CertPathFormatWithUpdater) SetCertLocalUpdater(name string) {
	c.LocalUpdater = name
}

func (c *CertPathFormatWithUpdater) CertContainerUpdater() string {
	return c.ContainerUpdater
}
