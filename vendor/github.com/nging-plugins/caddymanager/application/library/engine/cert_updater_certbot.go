package engine

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func init() {
	CertUpdaters.Add(`certbot`, `Certbot`, echo.KVOptX(CertUpdater{
		Update: RenewCertByCertbot,
		PathFormat: CertPathFormat{
			Cert:  `/etc/letsencrypt/live/{domain}/fullchain.pem`,
			Key:   `/etc/letsencrypt/live/{domain}/privkey.pem`,
			Trust: `/etc/letsencrypt/live/{domain}/chain.pem`,
		},
		DomainSanitizer: LegoSanitizedDomain,
	}))
}

// CertbotSanitizedDomain Make sure no funny chars are in the cert names (like wildcards ;)).
func CertbotSanitizedDomain(domain string) (string, error) {
	return strings.TrimPrefix(domain, `*`), nil
}

// http://coscms.com/.well-known/acme-challenge/Ito***l4-Fh7O5FpaAA*************LI3vTPo
// 申请：
// certbot certonly --webroot -d example.com --email info@example.com -w /var/www/_letsencrypt -n --agree-tos --force-renewal
// 更新
// certbot renew 更新所有
// certbot renew --cert-name example.com --force-renewal
func RenewCertByCertbot(ctx context.Context, cmdPathPrefix string, domains []string, email string, certDir string, isObtain bool) error {
	if len(domains) == 0 {
		return nil
	}
	command := `certbot`
	var args = []string{command}
	if sv, ok := ctx.Value(CtxCertDir).(string); ok && len(sv) > 0 {
		args = append(args, `--config-dir`, sv)
	}
	if isObtain {
		args = append(args, `certonly`, `--webroot`)
		for _, domain := range domains {
			args = append(args, `-d`, domain)
		}
		args = append(args,
			`--email`, email,
			`-w`, certDir,
			`-n`,
			`--agree-tos`,
			`--force-renewal`,
		)
	} else {
		args = append(args, `renew`)
		for _, domain := range domains {
			args = append(args, `--cert-name`, domain)
		}
		//args = append(args, `--force-renewal`)
	}
	if len(cmdPathPrefix) > 0 {
		rootArgs := com.ParseArgs(cmdPathPrefix)
		if len(rootArgs) > 1 {
			command = rootArgs[0]
			args = append(rootArgs[1:], args...)
		}
	}
	cmd := exec.CommandContext(ctx, command, args...)
	result, err := cmd.CombinedOutput()
	//log.Okay(cmd.String())
	if err != nil {
		err = fmt.Errorf(`%s: %w`, result, err)
	}
	return err
}
