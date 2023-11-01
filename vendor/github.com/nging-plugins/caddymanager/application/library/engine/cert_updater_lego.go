package engine

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"golang.org/x/net/idna"
)

func init() {
	CertUpdaters.Add(`lego`, `Lego`, echo.KVOptX(CertUpdater{
		Update: RenewCertByLego,
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
func RenewCertByLego(ctx context.Context, cmdPathPrefix string, domains []string, email string, certDir string, isObtain bool) error {
	if len(domains) == 0 {
		return nil
	}
	command := `lego`
	var args = []string{command}
	saveDir := `/etc/letsencrypt`
	if sv, ok := ctx.Value(CtxCertDir).(string); ok && len(sv) > 0 {
		saveDir = sv
	}
	args = append(args, `--path`, saveDir)
	args = append(args, `--email`, email, `--http`)
	for _, domain := range domains {
		args = append(args, `--domains`, domain)
	}
	if isObtain {
		args = append(args,
			`--http.webroot`, certDir,
			`--agree-tos`,
			`run`,
		)
	} else {
		args = append(args, `renew`)
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
