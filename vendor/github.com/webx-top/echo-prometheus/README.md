# Echo Prometheus
Middleware for echo to instrument all handlers as metrics


## Example of usage

### With default config
```go
package main

import (
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine/standard"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	echoPrometheus "github.com/webx-top/echo-prometheus"
)

func main() {
	e := echo.New()

	e.Use(echoPrometheus.MetricsMiddleware())
	e.Get("/metrics", echo.WrapHandler(promhttp.Handler()))

	e.Get("/", func(c echo.Context) error {
		return c.String("Hello, World!")
	})
	e.Logger().Fatal(e.Run(standard.New(":1323")))
}
```

### With custom config
```go
package main

import (
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine/standard"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	echoPrometheus "github.com/webx-top/echo-prometheus"
)

func main() {
	e := echo.New()

	var configMetrics = echoPrometheus.NewConfig()
		configMetrics.Namespace = "namespace"
		configMetrics.Buckets = []float64{
			0.0005, // 0.5ms
			0.001,  // 1ms
			0.005,  // 5ms
			0.01,   // 10ms
			0.05,   // 50ms
			0.1,    // 100ms
			0.5,    // 500ms
			1,      // 1s
			2,      // 2s
	}

	e.Use(echoPrometheus.MetricsMiddlewareWithConfig(configMetrics))
	e.Get("/metrics", echo.WrapHandler(promhttp.Handler()))

	e.Get("/", func(c echo.Context) error {
		return c.String("Hello, World!")
	})
	e.Logger().Fatal(e.Run(standard.New(":1323")))
}
```

### Gzip middleware

A middleware skipper can be passed to avoid gzip `/metrics` URL:

```go
e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
  Skipper: func(c echo.Context) bool {
    return c.Path() == "/metrics"
  },
}))
```

## Example output for metric route

```
# HELP echo_http_request_duration_seconds Spend time by processing a route
# TYPE echo_http_request_duration_seconds histogram
echo_http_request_duration_seconds_bucket{handler="/",method="GET",le="0.0005"} 7
echo_http_request_duration_seconds_bucket{handler="/",method="GET",le="0.001"} 7
echo_http_request_duration_seconds_bucket{handler="/",method="GET",le="0.002"} 7
echo_http_request_duration_seconds_bucket{handler="/",method="GET",le="0.005"} 7
echo_http_request_duration_seconds_bucket{handler="/",method="GET",le="0.01"} 7
echo_http_request_duration_seconds_bucket{handler="/",method="GET",le="0.02"} 7
echo_http_request_duration_seconds_bucket{handler="/",method="GET",le="0.05"} 7
echo_http_request_duration_seconds_bucket{handler="/",method="GET",le="0.1"} 7
echo_http_request_duration_seconds_bucket{handler="/",method="GET",le="0.2"} 7
echo_http_request_duration_seconds_bucket{handler="/",method="GET",le="0.5"} 7
echo_http_request_duration_seconds_bucket{handler="/",method="GET",le="1"} 7
echo_http_request_duration_seconds_bucket{handler="/",method="GET",le="2"} 7
echo_http_request_duration_seconds_bucket{handler="/",method="GET",le="5"} 7
echo_http_request_duration_seconds_bucket{handler="/",method="GET",le="+Inf"} 7
echo_http_request_duration_seconds_sum{handler="/",method="GET"} 7.645099999999999e-05
echo_http_request_duration_seconds_count{handler="/",method="GET"} 7
# HELP echo_http_requests_total Number of HTTP operations
# TYPE echo_http_requests_total counter
echo_http_requests_total{handler="/",method="GET",status="2xx"} 7
```

### View metrics via Grafana

We built a grafana dashboard for these metrics, lookup at

**1. Echo Framework Process Dashboard:**

本地位置：./grafana/tsuru-echo-framework-processes_rev1.json

在线位置：[https://grafana.com/grafana/dashboards/10913](https://grafana.com/grafana/dashboards/10913)

**2. Go Process Dashboard:**

本地位置：./grafana/go-processes_rev2.json

在线位置：[https://grafana.com/grafana/dashboards/6671](https://grafana.com/grafana/dashboards/6671)

## 使用方式

### 第一步：下载和启动 prometheus

1. 下载 prometheus

下载地址：[https://prometheus.io/download/](https://prometheus.io/download/)

下载命令举例：

```bash
wget https://github.com/prometheus/prometheus/releases/download/v2.20.0/prometheus-2.20.0.linux-amd64.tar.gz
```

解压：

```bash
tar -zxvf prometheus-2.20.0.linux-amd64.tar.gz
```

2. 启动 prometheus

如果可执行文件prometheus没有加入到环境变量PATH中，则执行：

```bash
cd prometheus-2.20.0.linux-amd64
./prometheus --config.file=$GOPATH/src/github.com/webx-top/echo-prometheus/prometheus/prometheus.yml
```

如果可执行文件prometheus已经加入到环境变量PATH中，则执行：

```bash
cd $GOPATH/src/github.com/webx-top/echo-prometheus/prometheus
prometheus
```


### 第二步：启动 grafana

```bash
grafana-server --config=/usr/local/etc/grafana/grafana.ini --homepath /usr/local/share/grafana cfg:default.paths.logs=/usr/local/var/log/grafana cfg:default.paths.data=/usr/local/var/lib/grafana cfg:default.paths.plugins=/usr/local/var/lib/grafana/plugins
```

### 第三步：配置 grafana

#### 一、配置数据源

点击路径：`Configuration` -> `Data Sources` -> `Add data source` -> 选择“`Prometheus`”：

1. 配置：Settings

    **HTTP**:

    `URL` -> 填写：`http://localhost:9090`  
    `Access` -> 选择：`Server(default)`  
    点击“`Save & Test`“按钮

2. 安装：Dashboards

   点击“`Dashboards`”选项卡，安装所有项目

#### 二、导入模版

点击路径：`Create` -> `Import` -> 点击“`Upload .json file`”按钮

分别选择本目录下grafana文件夹中的json文件进行上传，并在数据源中选择刚刚配置的数据源，然后点击“`Import`”按钮进行导入。

\- END -
