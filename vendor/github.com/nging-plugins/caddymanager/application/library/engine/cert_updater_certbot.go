package engine

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

func init() {
	CertUpdaters.Add(`certbot`, `Certbot`, echo.KVOptX(CertUpdater{
		MakeCommand: MakeCertbotCommand,
		Update:      RenewCertByCertbot,
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
// documention: https://eff-certbot.readthedocs.io/en/latest/using.html#certbot-commands
//
// certbot certonly \
// --dns-cloudflare \
// --dns-cloudflare-credentials ~/.secrets/certbot/cloudflare.ini \
// -d example.com \
// -d www.example.com
func MakeCertbotCommand(data RequestCertUpdate) (command string, args []string, env []string) {
	command = `certbot`
	args = []string{command}
	if len(data.CertSaveDir) > 0 {
		args = append(args, `--config-dir`, data.CertSaveDir)
	}
	if data.Obtain {
		args = append(args, `certonly`, `--webroot`, `--email`, data.Email)
		for _, domain := range data.Domains {
			args = append(args, `-d`, domain)
		}
		if len(data.DNSProvider) > 0 {
			args = append(args,
				`--dns-`+data.DNSProvider,
				`--dns-`+data.DNSProvider+`-credentials`, data.DNSCredentials,
				`--dns-`+data.DNSProvider+`-propagation-seconds`, param.AsString(data.DNSWaitSeconds),
			)
		} else {
			args = append(args, `-w`, data.CertVerifyDir)
		}
		args = append(args,
			`-n`,
			`--agree-tos`,
			`--force-renewal`,
		)
	} else {
		args = append(args, `renew`)
		for _, domain := range data.Domains {
			args = append(args, `--cert-name`, domain)
		}
		//args = append(args, `--force-renewal`)
	}
	if len(data.CmdPathPrefix) > 0 {
		rootArgs := com.ParseArgs(data.CmdPathPrefix)
		if len(rootArgs) > 1 {
			command = rootArgs[0]
			args = append(rootArgs[1:], args...)
		}
	}
	env = data.Env
	return
}

func RenewCertByCertbot(ctx context.Context, data RequestCertUpdate) error {
	if len(data.Domains) == 0 {
		return nil
	}
	command, args, env := MakeCertbotCommand(data)
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, env...)
	result, err := cmd.CombinedOutput()
	//log.Okay(cmd.String())
	if err != nil {
		err = fmt.Errorf(`%s: %w`, result, err)
	}
	return err
}
