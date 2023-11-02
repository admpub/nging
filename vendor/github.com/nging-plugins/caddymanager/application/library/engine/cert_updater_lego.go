package engine

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"golang.org/x/net/idna"
)

func init() {
	CertUpdaters.Add(`lego`, `Lego`, echo.KVOptX(CertUpdater{
		MakeCommand: MakeLegoCommand,
		Update:      RenewCertByLego,
		PathFormat: CertPathFormat{
			Cert:  `/etc/letsencrypt/.lego/certificates/{domain}.crt`,
			Key:   `/etc/letsencrypt/.lego/certificates/{domain}.key`,
			Trust: ``,
		},
		DomainSanitizer: LegoSanitizedDomain,
	}))
}

// 申请：
// lego --accept-tos --email you@example.com --http --http.webroot /path/to/webroot --domains example.com run
// https://go-acme.github.io/lego/usage/cli/obtain-a-certificate/
// 更新：
// lego --email="you@example.com" --domains="example.com" --http renew
// https://go-acme.github.io/lego/usage/cli/renew-a-certificate/
//
// CLOUDFLARE_EMAIL="you@example.com" \
// CLOUDFLARE_API_KEY="yourprivatecloudflareapikey" \
// lego --email "you@example.com" --dns cloudflare --domains "example.org" run
func MakeLegoCommand(data RequestCertUpdate) (command string, args []string, env []string) {
	command = `lego`
	args = []string{command}
	saveDir := `/etc/letsencrypt`
	if len(data.CertSaveDir) > 0 {
		saveDir = data.CertSaveDir
	}
	args = append(args, `--path`, saveDir)
	args = append(args, `--email`, data.Email)
	if len(data.DNSProvider) > 0 {
		args = append(args, `--dns`, data.DNSProvider)
	} else {
		args = append(args, `--http`)
	}
	for _, domain := range data.Domains {
		args = append(args, `--domains`, domain)
	}
	if data.Obtain {
		args = append(args,
			`--http.webroot`, data.CertVerifyDir,
			`--agree-tos`,
			`run`,
		)
	} else {
		args = append(args, `renew`)
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

func RenewCertByLego(ctx context.Context, data RequestCertUpdate) error {
	if len(data.Domains) == 0 {
		return nil
	}
	command, args, env := MakeLegoCommand(data)
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

var (
	domainReplacer = strings.NewReplacer(
		":", "-",
		"*", "_",
	)
)

// LegoSanitizedDomain Make sure no funny chars are in the cert names (like wildcards ;)).
func LegoSanitizedDomain(domain string) (string, error) {
	return idna.ToASCII(domainReplacer.Replace(domain))
}
