# 开发说明

本软件项目不仅仅实现了一些网站服务工具，本身还是一个具有很好扩展性的通用网站后台管理系统，通过本项目，您可以很轻松的构建一个全新的网站项目，省去从头构建项目的麻烦，减少重复性劳动。

## Ⅰ、模块说明

当您基于本项目来构建新软件的时候，您可以根据需要来选用本系统的网站服务工具：
```go
import (
	"github.com/admpub/nging/v5/application/library/module"

	// module
	"github.com/admpub/nging/v5/application/handler/cloud"
	"github.com/admpub/nging/v5/application/handler/task"
	"github.com/nging-plugins/caddymanager"
	"github.com/nging-plugins/collector"
	"github.com/nging-plugins/dbmanager"
	"github.com/nging-plugins/ddnsmanager"
	"github.com/nging-plugins/dlmanager"
	"github.com/nging-plugins/frpmanager"
	"github.com/nging-plugins/ftpmanager"
	"github.com/nging-plugins/servermanager"
	"github.com/nging-plugins/sshmanager"
)
```
并注册功能模块
```go
func main(){
    initModule()
}

func initModule() {
	module.Register(
		&caddymanager.Module,
		&servermanager.Module,
		&ftpmanager.Module,
		&collector.Module,
		&task.Module,
		&dlmanager.Module,
		&cloud.Module,
		&dbmanager.Module,
		&frpmanager.Module,
		&sshmanager.Module,
		&ddnsmanager.Module,
	)
}
```

## Ⅱ、开发环境下的启动方式

- 第一步： 安装GO环境(必须1.12.1版以上)，配置GOPATH、GOROOT环境变量，并将`%GOROOT%/bin`和`%GOPATH%/bin`加入到PATH环境变量中
- 第二步： 执行命令`go get github.com/admpub/nging`
- 第三步： 进入`%GOPATH%/src/github.com/admpub/nging/`目录中启动`run_first_time.bat`(linux系统启动`run_first_time.sh`)
- 第四步： 打开浏览器，访问网址`http://localhost:8080/setup`，在页面中配置数据库账号和管理员账号信息进行安装
- 第五步： 安装成功后会自动跳转到登录页面，使用安装时设置的管理员账号进行登录
