package engine

import "errors"

var (
	ErrNotSetCertContainerDir          = errors.New(`CertContainerDir value is not set`)
	ErrNotSetCertLocalDir              = errors.New(`CertLocalDir value is not set`)
	ErrNotSetEngineConfigLocalFile     = errors.New(`EngineConfigLocalFile value is not set`)
	ErrNotSetEngineConfigContainerFile = errors.New(`EngineConfigContainerFile value is not set`)
	ErrNotSetVhostConfigLocalDir       = errors.New(`VhostConfigLocalDir value is not set`)
	ErrNotSetVhostConfigContainerDir   = errors.New(`VhostConfigContainerDir value is not set`)
)
