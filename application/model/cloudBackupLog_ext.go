package model

import "github.com/webx-top/echo"

const (
	// type
	CloudBackupTypeFull   = `full`
	CloudBackupTypeChange = `change`

	// status
	CloudBackupStatusSuccess = `success`
	CloudBackupStatusFailure = `failure`
)

var CloudBackupTypes = echo.NewKVData().Add(CloudBackupTypeFull, `全量备份`).Add(CloudBackupTypeChange, `监控备份`)
var CloudBackupStatuses = echo.NewKVData().Add(CloudBackupStatusSuccess, `成功`).Add(CloudBackupStatusFailure, `失败`)
