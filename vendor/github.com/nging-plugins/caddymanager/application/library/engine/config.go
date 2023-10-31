package engine

import (
	"context"
	"os"
)

type Configer interface {
	GetVhostConfigLocalDirAbs() (string, error)
	GetTemplateFile() string
	GetIdent() string
	GetEngine() string
	GetEnviron() string
	GetCertLocalDir() string
	GetCertContainerDir() string
	GetVhostConfigLocalDir() string
	GetVhostConfigContainerDir() string
	GetEngineConfigLocalFile() string
	GetEngineConfigContainerFile() string
}

type CertPathFormat struct {
	Cert  string
	Key   string
	Trust string
}

type CertPathFormatGetter interface {
	GetCertPathFormat() CertPathFormat
}

type CertRenewaler interface {
	RenewalCert(ctx context.Context, id uint, domains []string, email string) error
}

type EngineConfigFileFixer interface {
	FixEngineConfigFile(deleteMode ...bool) (bool, error)
}

type VhostConfigRemover interface {
	RemoveVhostConfig(id uint) error
}

type CertFileRemover interface {
	RemoveCertFile(id uint) error
}

func FixEngineConfigFile(cfg Configer, deleteMode ...bool) (bool, error) {
	if fx, ok := cfg.(EngineConfigFileFixer); ok {
		hasUpdate, err := fx.FixEngineConfigFile(deleteMode...)
		if err != nil && os.IsNotExist(err) {
			return hasUpdate, nil
		}
		return hasUpdate, err
	}
	return false, nil
}

// 如果 id 为 0 代表删除此配置下的所有其它文件
func RemoveVhostConfigFile(cfg Configer, id uint) error {
	if rm, ok := cfg.(VhostConfigRemover); ok {
		err := rm.RemoveVhostConfig(id)
		if err != nil && os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return nil
}

func RemoveCertFile(cfg Configer, id uint) error {
	if rm, ok := cfg.(CertFileRemover); ok {
		err := rm.RemoveCertFile(id)
		if err != nil && os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return nil
}

func RenewalCert(cfg Configer, ctx context.Context, id uint, domains []string, email string) error {
	if rm, ok := cfg.(CertRenewaler); ok {
		err := rm.RenewalCert(ctx, id, domains, email)
		if err != nil && os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return nil
}
