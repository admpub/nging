package driver

import (
	"context"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/url"
	"os"
	"path"
	"sort"

	uploadClient "github.com/webx-top/client/upload"
	"github.com/webx-top/client/upload/watermark"
	"github.com/webx-top/echo"

	"github.com/admpub/checksum"
	"github.com/admpub/color"
	"github.com/admpub/log"
	"github.com/admpub/nging/application/registry/upload/table"
	"github.com/admpub/nging/application/model/file/storer"
)

var (
	// ErrExistsFile 文件不存在
	ErrExistsFile = table.ErrExistsFile
)

// BatchUpload 批量上传
func BatchUpload(
	ctx echo.Context,
	fieldName string,
	dstNamer func(*uploadClient.Result) (dst string, err error),
	storer Storer,
	callback func(*uploadClient.Result, multipart.File) error,
	markFile string,
) (results uploadClient.Results, err error) {
	req := ctx.Request()
	if req == nil {
		err = ctx.E(`Invalid upload content`)
		return
	}
	m := req.MultipartForm()
	if m == nil || m.File == nil {
		err = ctx.E(`Invalid upload content`)
		return
	}
	files, ok := m.File[fieldName]
	if !ok {
		err = echo.ErrNotFoundFileInput
		return
	}
	var dstFile string
	for _, fileHdr := range files {
		//for each fileheader, get a handle to the actual file
		var file multipart.File
		file, err = fileHdr.Open()
		if err != nil {
			if file != nil {
				file.Close()
			}
			return
		}
		result := &uploadClient.Result{
			FileName: fileHdr.Filename,
			FileSize: fileHdr.Size,
		}
		result.Md5, err = checksum.MD5sumReader(file)
		if err != nil {
			file.Close()
			return
		}

		dstFile, err = dstNamer(result)
		if err != nil {
			file.Close()
			if err == ErrExistsFile {
				results.Add(result)
				err = nil
				continue
			}
			return
		}
		if len(dstFile) == 0 {
			file.Close()
			continue
		}
		if len(result.SavePath) > 0 {
			file.Close()
			results.Add(result)
			continue
		}
		file.Seek(0, 0)
		if len(markFile) > 0 && result.FileType.String() == `image` {
			var b []byte
			b, err = ioutil.ReadAll(file)
			if err != nil {
				file.Close()
				return
			}
			b, err = watermark.Bytes(b, path.Ext(result.FileName), markFile)
			if err != nil {
				file.Close()
				return
			}
			file = watermark.Bytes2file(b)
			result.SavePath, result.FileURL, err = storer.Put(dstFile, file, int64(len(b)))
		} else {
			result.SavePath, result.FileURL, err = storer.Put(dstFile, file, fileHdr.Size)
		}
		if err != nil {
			file.Close()
			return
		}
		file.Seek(0, 0)
		if err = callback(result, file); err != nil {
			file.Close()
			return
		}
		file.Close()
		results.Add(result)
	}
	return
}

// Sizer 尺寸接口
type Sizer interface {
	Size() int64
}

// Storer 文件存储引擎接口
type Storer interface {
	// 引擎名
	Name() string

	// FileDir 文件夹物理路径
	FileDir(subpath string) string

	// URLDir 文件夹网址路径
	URLDir(subpath string) string

	// Put 保存文件
	Put(dst string, src io.Reader, size int64) (savePath string, viewURL string, err error)

	// Get 获取文件
	Get(file string) (io.ReadCloser, error)

	// Exists 文件是否存在
	Exists(file string) (bool, error)

	// FileInfo 文件信息
	FileInfo(file string) (os.FileInfo, error)

	// SendFile 输出文件到浏览器
	SendFile(ctx echo.Context, file string) error

	// Delete 删除文件
	Delete(file string) error

	// DeleteDir 删除目录
	DeleteDir(dir string) error

	// Move 移动文件
	Move(src, dst string) error

	// PublicURL 文件物理路径转网址
	PublicURL(dst string) string

	// URLToFile 网址转文件存储路径(非完整路径)
	URLToFile(viewURL string) string

	// URLToPath 网址转文件路径(完整路径)
	URLToPath(viewURL string) string

	// 根网址(末尾不含"/")
	SetBaseURL(baseURL string) string
	BaseURL() string

	// FixURL 修正网址
	FixURL(content string, embedded ...bool) string

	// FixURLWithParams 修正网址并增加网址参数
	FixURLWithParams(content string, values url.Values, embedded ...bool) string

	// Close 关闭连接
	Close() error
}

// Constructor 存储引擎构造函数
type Constructor func(ctx context.Context, typ string) (Storer, error)

var storers = map[string]Constructor{}

// DefaultConstructor 默认构造器
var DefaultConstructor Constructor

// Register 存储引擎注册
func Register(engine string, constructor Constructor) {
	log.Info(color.CyanString(`storer.register:`), engine)
	storers[engine] = constructor
}

// Get 获取存储引擎构造器
func Get(engine string) Constructor {
	constructor, ok := storers[engine]
	if !ok {
		return DefaultConstructor
	}
	return constructor
}

// GetBySettings 获取存储引擎构造器
func GetBySettings() Constructor {
	engine := `local`
	storerConfig, ok := storer.GetOk()
	if ok {
		engine = storerConfig.Name
	}
	return Get(engine)
}

// All 存储引擎集合
func All() map[string]Constructor {
	return storers
}

// AllNames 存储引擎集合
func AllNames() []string {
	var names []string
	for name := range storers {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
