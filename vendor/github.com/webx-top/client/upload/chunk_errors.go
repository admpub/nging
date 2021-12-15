package upload

import "errors"

var (
	// Success
	ErrFileUploadCompleted  = errors.New("文件已经上传完成")
	ErrChunkUploadCompleted = errors.New("文件分片已经上传完成")

	// Support
	ErrChunkUnsupported = errors.New("不支持分片上传")

	// Failure
	ErrChunkHistoryOpenFailed     = errors.New("打开历史分片文件失败")
	ErrChunkMergeFileCreateFailed = errors.New("创建分片合并文件失败")
	ErrChunkFileOpenFailed        = errors.New("分片文件打开失败")
	ErrChunkFileMergeFailed       = errors.New("分片文件合并失败")
	ErrChunkFileDeleteFailed      = errors.New("分片文件删除失败")
)
