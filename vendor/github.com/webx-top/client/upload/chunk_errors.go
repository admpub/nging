package upload

import "errors"

var (
	// Success
	ErrFileUploadCompleted  = errors.New("文件已经上传完成")
	ErrChunkUploadCompleted = errors.New("文件分片已经上传完成")

	// Support
	ErrChunkUnsupported        = errors.New("不支持分片上传")
	ErrFileSizeExceedsLimit    = errors.New("文件尺寸超出限制")
	ErrRequestBodyExceedsLimit = errors.New("提交内容尺寸超出限制")
	ErrIncorrectSize           = errors.New("文件尺寸不正确")

	// Failure
	ErrChunkHistoryOpenFailed     = errors.New("打开历史分片文件失败")
	ErrChunkMergeFileCreateFailed = errors.New("创建分片合并文件失败")
	ErrChunkFileOpenFailed        = errors.New("分片文件打开失败")
	ErrChunkFileMergeFailed       = errors.New("分片文件合并失败")
	ErrChunkFileDeleteFailed      = errors.New("分片文件删除失败")
)
