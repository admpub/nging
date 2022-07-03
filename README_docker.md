# 用 Docker 启动 Nging

## Nging + MySQL

```sh
docker-compose -f ./docker-compose-mysql.yml up --build -d
```

在页面 `/setup` 安装页面，数据库信息设置如下：

* 数据库选择 `MySQL`
* 主机地址输入 `mysql:3306`
* 用户名输入 `root`
* 密码输入 `root`

### Nging + SQLite

```sh
docker-compose -f ./docker-compose-sqlite.yml up --build -d
```

在页面 `/setup` 安装页面，数据库信息设置如下：

* 数据库选择 `SQLite`
* 数据库名称输入 `myconfig/nging.db`

