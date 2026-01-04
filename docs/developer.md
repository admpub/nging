# 开发说明

本软件项目不仅仅实现了一些网站服务工具，本身还是一个具有很好扩展性的通用网站后台管理系统，通过本项目，您可以很轻松的构建一个全新的网站项目，省去从头构建项目的麻烦，减少重复性劳动。

## Ⅰ、模块说明

当您基于本项目来构建新软件的时候，您可以根据需要来选用本系统的网站服务工具：
```go
import (
	"github.com/coscms/webcore/library/module"

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

## 注意事项
1. 为了便于css代码清理工具准确运行，在模板和js文件中写css类名时确保完整，不要拆开来写  

  	例如：
  	```html
  	<span class="wd-{{if $v}}100{{else}}200{{end}} text-primary">content</span>
  	```
  	应该写为
  	```html
  	<span class="{{if $v}}wd-100{{else}}wd-200{{end}} text-primary">content</span>
  	```

2. 支持 SQLite3

	默认是不支持 SQLite3 数据库的，如要支持，需要在编译的时候添加 tag `db_sqlite`。
	
	此时使用纯 go 版本的 sqlite 驱动，如欲使用更高性能的 cgo 版本驱动，再次添加 tag `sqlitecgo` 即可

	例如：
	```bash
	go build -tags db_sqlite,sqlitecgo
	```

## 模板引擎

* [模板基础语法](https://admpub.com/blog/post/admin/Golang%E6%A8%A1%E6%9D%BF%E8%AF%AD%E6%B3%95%E7%AE%80%E6%98%8E%E6%95%99%E7%A8%8B%5B%E6%96%B0%5D-2)
* [模板引擎文档](https://github.com/webx-top/echo/blob/master/middleware/render/standard/README.md)
* [模板自定义全局函数文档](https://github.com/webx-top/echo/blob/master/middleware/tplfunc/TplFuncMap_Documentation.md)

## 前端 JS 库文档

### 官方库文档
* [JQuery 3.5.1 文档](https://api.jquery.com/)
* [Bootstrap 3.4.1 文档](https://getbootstrap.com/docs/3.4/getting-started/)
* [Select2 4.0.13 文档](https://select2.org/)
* [Typeahead.js 0.11.1 文档](https://github.com/twitter/typeahead.js/blob/master/doc/jquery_typeahead.md)
* [Gritter 1.7.0 文档](https://github.com/jboesch/Gritter/blob/master/README.md)
* [NProgress 0.2.0 文档](https://ricostacruz.com/nprogress/)
* [daterangepicker 2.1.25 文档](https://www.daterangepicker.com/)
* [NanoScroller 0.9.1 文档](https://jamesflorentino.github.io/nanoScrollerJS/)

### 自定义 JS 库文档

* [general.js 文档](developer/general.js.md)
* [editor.js 文档](developer/editor.js.md)
* [modal.js 文档](developer/modal.js.md)
