package engine

import (
	"context"

	"github.com/webx-top/echo"
)

type DomainSanitizer func(string) (string, error)

type CertExecuteor func(ctx context.Context, data RequestCertUpdate) error

type CertUpdater struct {
	Update          CertExecuteor
	MakeCommand     func(RequestCertUpdate) (command string, args []string, env []string)
	PathFormat      CertPathFormat
	DomainSanitizer DomainSanitizer
}

type RequestCertUpdate struct {
	CmdPathPrefix  string
	Domains        []string
	Email          string
	CertSaveDir    string // 保存证书的目录
	CertVerifyDir  string // 保存证书验证素材的目录
	Obtain         bool
	DNSProvider    string
	DNSCredentials string
	DNSWaitSeconds int
	Env            []string
}

var CertUpdaters = echo.NewKVData()
