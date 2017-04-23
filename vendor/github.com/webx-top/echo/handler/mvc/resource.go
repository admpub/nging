package mvc

import "github.com/webx-top/echo/handler/mvc/static/resource"

// Resource 静态资源
type Resource struct {
	*resource.Static
	Dir string //静态资源所在文件夹
}
