/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package seaweedfs

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

//step1. weed master
//step2. weed volume -port=9001 -dir=./_test
//step3. weed filer -collection=test -port=8888 -port.readonly=8989 -master=localhost:9333
// or weed filer -collection=test -port=8888 -port.readonly=8989 -master=localhost:9333,localhost:9334
//============================================================================
// 映射为本地虚拟目录(此时可以使用filesystem驱动像操作本地文件一样操作seaweedfs中的文件)
//============================================================================
//挂载：sudo weed mount -filer=localhost:8888 -dir=./public/upload/test -filer.path=/ -collection=test
//取消挂载：关闭服务后执行 sudo umount -l ./public/upload/test

// SeaweedFS对图片有很好的支持，可以指定图片显示的长度、宽度、模式，如：
// http://localhost:8989/3/01637037d6.jpg?height=200&width=200
// http://localhost:8989/3/01637037d6.jpg?height=200&width=200&mode=fit
// http://localhost:8989/3/01637037d6.jpg?height=200&width=200&mode=fill

// 其它fuse方案：
// https://github.com/kahing/goofys
// https://github.com/ncw/rclone :
// 1. rclone配置。执行 rclone config 开始配置rclone
// 2. rclone挂载网盘。命令示例：rclone mount config2:image-dir /root/image-dir --allow-non-empty，此命令会把本地目录/root/image-dir ，挂载到网盘（config2配置文件名）目录image-dir 下。如果你的目标目录中有文件，可以使用参数：--allow-non-empty，但是注意它会清空原目录中的文件。使用rclone mount 命令后，如果不守护进程的话，会话结束后rclone mount 进程就会终止。我们可以用screen来让rclone保持在后台运行。在执行rclone挂载命令前，选择执行：screen –S config2，挂载完成后再用快捷键CTRL-a d 来暂时断开当前会话。最后用screen -r <screen_pid>重新连接上。

func TestSeaweedfs(t *testing.T) {
	return
	r := NewSeaweedfs(nil, `test`)
	f, err := os.Open(`./config.go`)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}
	_, purl, err := r.Put(`/config.go`, f, fi.Size())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(purl)
	return
	err = r.Delete(path.Base(purl))
	if err != nil {
		t.Fatal(err)
	}
	var html string

	assert.Equal(t, "<h2>安装 Go 第三方包 go-sqlite3</h2>", html)
	// 成功取得HTML内容进行后续处理
	fmt.Println(html)
}
