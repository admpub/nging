package cmder

import (
	"io"
)

type Cmder interface {
	Boot() error // 服务自身的启动逻辑
	Control      // 控制操作
}

type Control interface {
	StopHistory(...string) error       // 停止已启动服务
	Start(writer ...io.Writer) error   // 启动服务
	Stop() error                       // 停止服务
	Reload() error                     // 重载服务
	Restart(writer ...io.Writer) error // 重启服务
}

type RestartBy interface {
	RestartBy(id string, writer ...io.Writer) error // 重启“服务中指定 ID 所指向的项“
}

type StopBy interface {
	StopBy(id string) error // 停止”服务中指定 ID 所指向的项“
}

type StartBy interface {
	StartBy(id string, writer ...io.Writer) error // 启动”服务中指定 ID 所指向的项“
}
