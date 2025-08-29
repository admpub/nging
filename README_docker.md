# 用 Docker 启动 Nging

**基础版**镜像：`docker.io/admpub/nging:latest`  (Docker主页：https://hub.docker.com/r/admpub/nging)  
**先锋版**镜像：`docker.io/admpub/nging-dockermgr:latest`  (Docker主页：https://hub.docker.com/r/admpub/nging-dockermgr)  

复制 ./docker-compose-mysql.yml 或 ./docker-compose-sqlite.yml 文件到自己的文件夹，按需修改文件内的相关参数后再用 `docker-compose` 或 `docker compose` 启动。

## Nging + MySQL

```sh
docker-compose -f ./docker-compose-mysql.yml up --build -d
```

在安装页面 `/setup` ，数据库信息设置如下：

* 数据库选择 `MySQL`
* 主机地址输入 `mysql:3306`
* 用户名输入 `root`
* 密码输入 `root`

### Nging + SQLite

```sh
docker-compose -f ./docker-compose-sqlite.yml up --build -d
```

在安装页面 `/setup` ，数据库信息设置如下：

* 数据库选择 `SQLite`
* 数据库名称输入 `myconfig/nging.db`

