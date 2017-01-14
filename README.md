# nging
基于 caddy 的网站服务程序，带图形化管理界面。

caddy 是由国外开发者开发的一套类似于nginx或apache的网站服务软件。
caddy的配置文件比nginx更简洁易用。但我相信事情还可以变得更简单，所以nging应运而生。

nging不仅仅包含了caddy的在线可视化配置，还包含了ftp服务的管理，下一步即将增加：

- [x] 文件在线管理
- [x] 数据库管理
- [ ] 支持更多caddy指令的在线配置

# 先睹为快
[![](https://github.com/admpub/nging/blob/master/preview/preview_login.png?raw=true)](https://github.com/admpub/nging/blob/master/preview/preview_login.png)
[![](https://github.com/admpub/nging/blob/master/preview/preview_sysinfo.png?raw=true)](https://github.com/admpub/nging/blob/master/preview/preview_sysinfo.png)
[![](https://github.com/admpub/nging/blob/master/preview/preview_listtable.png?raw=true)](https://github.com/admpub/nging/blob/master/preview/preview_listtable.png)

# 开发环境下的启动方式

- 第一步： 安装GO环境(建议1.7版以上)，配置GOPATH、GOROOT环境变量，并将`%GOROOT%/bin`和`%GOPATH%/bin`加入到PATH环境变量中
- 第二步： 执行命令`go get github.com/admpub/nging`
- 第三步： 进入`%GOPATH%/src/github.com/admpub/nging/`目录中启动`run_first_time.bat`(linux系统启动`run_first_time.sh`)
- 第四步： 打开浏览器，访问网址`http://localhost:8080/setup`，在页面中配置数据库账号信息进行安装
- 第五步： 安装成功后会自动跳转到登录页面，使用初始管理账号进行登录(初始管理账号和密码分别为`admin`和`admin`)

# 备注
本项目目前正处于开发阶段，一些功能尚未完成，也会存在一些bug，所以请不要用在生产环境。