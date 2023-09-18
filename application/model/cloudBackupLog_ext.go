package model

import "github.com/webx-top/echo"

const (
	// type
	CloudBackupTypeFull   = `full`
	CloudBackupTypeChange = `change`

	// status
	CloudBackupStatusSuccess = `success`
	CloudBackupStatusFailure = `failure`

	// operation
	CloudBackupOperationCreate = `create`
	CloudBackupOperationUpdate = `update`
	CloudBackupOperationDelete = `delete`
	CloudBackupOperationNone   = `none`
)

var CloudBackupTypes = echo.NewKVData().Add(CloudBackupTypeFull, `全量`).Add(CloudBackupTypeChange, `监控`)
var CloudBackupStatuses = echo.NewKVData().Add(CloudBackupStatusSuccess, `成功`).Add(CloudBackupStatusFailure, `失败`)
var CloudBackupOperations = echo.NewKVData().Add(CloudBackupOperationCreate, `创建`).
	Add(CloudBackupOperationUpdate, `更新`).
	Add(CloudBackupOperationDelete, `删除`).
	Add(CloudBackupOperationNone, `未知`)
