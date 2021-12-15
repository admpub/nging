# Mori's Goforever Fork

See the original at [https://github.com/gwoo/goforever]()

Config based process manager. Goforever could be used in place of supervisor, runit, node-forever, etc.
Goforever will start an http server on the specified port.

	Usage of ./goforever:
	  -conf="goforever.toml": Path to config file.

## 功能描述
1. 主进程退出不影响子进程，当主进程恢复服务后继续根据pid文件进行监控进程状态
2. 进程被外部程序结束后自动启动（respawn > 0）
3. 有延迟启动功能（delay）
4. 心跳时间可配置（ping）

## 配置文件模板
	ip = "127.0.0.1"	# API 监听地址
	port = "2224" 		# API 监听端口
	username = "go" 	# API 用户名
	password = "forever" 	# API 密码
	pidfile = "goforever.pid"
	logfile = "goforever.log"
	errfile = "goforever.log"

	[[process]]
	name = "example-panic"
	command = "./example/example-panic"
	pidfile = "example/example-panic.pid"
	logfile = "example/logs/example-panic.debug.log"
	errfile = "example/logs/example-panic.errors.log"
	respawn = 1 	# 进程被结束时自动重启
	delay = "1m" 	# 延迟一分钟启动
	ping = "30s"	# 30s检测一次进状态

	[[process]]
	name = "example"
	dir = "/tmp/"			# 切换到此目录下运行 command
	env = ["aaa=aaa","bbb=bbb"] 	# 程序指定的环境变量
	command = "./example/example"
	args = ["-name=foo"] 		# 命令行参数
	pidfile = "example/example.pid"
	logfile = "example/logs/example.debug.log"
	errfile = "example/logs/example.errors.log"
	respawn = 1

## Running
Help.

	./goforever -h

Daemonize main process.

	./goforever start

Run main process and output to current session.

	./goforever

## CLI
	list				List processes.
	show [process]	    Show a main proccess or named process.
	start [process]		Start a main proccess or named process.
	stop [process]		Stop a main proccess or named process.
	restart [process]	Restart a main proccess or named process.

## HTTP API

Return a list of managed processes

	GET host:port/

Start the process

	POST host:port/:name

Restart the process

	PUT host:port/:name

Stop the process

	DELETE host:port/:name
