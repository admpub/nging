package engine

import (
	"context"

	"github.com/webx-top/echo"
)

type DomainSanitizer func(string) (string, error)

type CertExecuteor func(ctx context.Context, cmdPathPrefix string, domains []string, email string, certDir string, isObtain bool) error

type CertUpdater struct {
	Update          CertExecuteor
	PathFormat      CertPathFormat
	DomainSanitizer DomainSanitizer
}

var CertUpdaters = echo.NewKVData()

type CtxKeyCertDir string

var CtxCertDir = `certDir`
