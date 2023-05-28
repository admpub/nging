package ftp

type PathType int

// 路径类型
const (
	PathTypeDir  PathType = 0 // 文件夹
	PathTypeFile PathType = 1 // 文件
	PathTypeBoth PathType = 2 // 不限
)

type Operate int

// 操作类型
const (
	OperateRead   Operate = 0 // 读
	OperateModify Operate = 1 // 修改
	OperateCreate Operate = 2 // 新建
)
