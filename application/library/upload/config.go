package upload

import (
	"mime"
	"strings"

	"github.com/admpub/nging/v4/application/library/config"
	"github.com/admpub/nging/v4/application/library/config/extend"
	uploadClient "github.com/webx-top/client/upload"
	"github.com/webx-top/com"
	"github.com/webx-top/echo/middleware/bytes"
)

const Name = `upload`

const (
	defaultUploadURLPath = `/public/upload/`
	defaultUploadDir     = `./public/upload`
)

var (
	// UploadURLPath 上传文件网址访问路径
	UploadURLPath = defaultUploadURLPath

	// UploadDir 定义上传目录（首尾必须带“/”）
	UploadDir = defaultUploadDir

	// AllowedExtensions 被允许上传的文件的扩展名
	AllowedExtensions = []string{
		`.jpeg`, `.jpg`, `.gif`, `.png`, `.mp4`,
	}
)

func init() {
	extend.Register(Name, func() interface{} {
		return NewConfig()
	})
}

func Get() *Config {
	cfg, _ := config.MustGetConfig().Extend.Get(Name).(*Config)
	if cfg == nil {
		cfg = NewConfig()
	}
	return cfg
}

type FileType struct {
	Icon         string   `json:"icon"`
	Extensions   []string `json:"extensions"`
	MIMEs        []string `json:"mimes"`
	MIMEKeywords []string `json:"mimeKeywords"`
	MaxSize      string   `json:"maxSize"`
	Description  string   `json:"description"`
	Disabled     bool     `json:"disabled"`
	maxSizeBytes int
}

func (c *FileType) Init() {
	max, _ := bytes.Parse(c.MaxSize)
	c.maxSizeBytes = int(max)
}

func (c *FileType) MaxSizeBytes() int {
	return c.maxSizeBytes
}

type Config struct {
	FileTypes         map[string]*FileType `json:"fileTypes"`
	MaxSize           string               `json:"maxSize"`
	Icon              string               `json:"icon"`
	AllowedExtensions []string             `json:"allowedExtensions"`
	maxSizeBytes      int
	fileTypes         map[string]string
}

func NewConfig() *Config {
	return &Config{
		FileTypes: map[string]*FileType{
			`image`: {
				Icon:         uploadClient.FileTypeIcons[uploadClient.TypeImage.String()],
				Extensions:   uploadClient.FileTypeExts[uploadClient.TypeImage],
				MIMEKeywords: uploadClient.FileTypeMimeKeywords[uploadClient.TypeImage.String()],
			},
			`video`: {
				Icon:         uploadClient.FileTypeIcons[uploadClient.TypeVideo.String()],
				Extensions:   uploadClient.FileTypeExts[uploadClient.TypeVideo],
				MIMEKeywords: uploadClient.FileTypeMimeKeywords[uploadClient.TypeVideo.String()],
				MaxSize:      `200M`,
			},
			`audio`: {
				Icon:         uploadClient.FileTypeIcons[uploadClient.TypeAudio.String()],
				Extensions:   uploadClient.FileTypeExts[uploadClient.TypeAudio],
				MIMEKeywords: uploadClient.FileTypeMimeKeywords[uploadClient.TypeAudio.String()],
			},
			`archive`: {
				Icon:         uploadClient.FileTypeIcons[uploadClient.TypeArchive.String()],
				Extensions:   uploadClient.FileTypeExts[uploadClient.TypeArchive],
				MIMEKeywords: uploadClient.FileTypeMimeKeywords[uploadClient.TypeArchive.String()],
			},
			`pdf`: {
				Icon:         uploadClient.FileTypeIcons[uploadClient.TypePDF.String()],
				Extensions:   uploadClient.FileTypeExts[uploadClient.TypePDF],
				MIMEKeywords: uploadClient.FileTypeMimeKeywords[uploadClient.TypePDF.String()],
			},
			`xls`: {
				Icon:         uploadClient.FileTypeIcons[uploadClient.TypeXLS.String()],
				Extensions:   uploadClient.FileTypeExts[uploadClient.TypeXLS],
				MIMEKeywords: uploadClient.FileTypeMimeKeywords[uploadClient.TypeXLS.String()],
			},
			`ppt`: {
				Icon:         uploadClient.FileTypeIcons[uploadClient.TypePPT.String()],
				Extensions:   uploadClient.FileTypeExts[uploadClient.TypePPT],
				MIMEKeywords: uploadClient.FileTypeMimeKeywords[uploadClient.TypePPT.String()],
			},
			`doc`: {
				Icon:         uploadClient.FileTypeIcons[uploadClient.TypeDOC.String()],
				Extensions:   uploadClient.FileTypeExts[uploadClient.TypeDOC],
				MIMEKeywords: uploadClient.FileTypeMimeKeywords[uploadClient.TypeDOC.String()],
			},
			`bt`: {
				Icon:         uploadClient.FileTypeIcons[uploadClient.TypeBT.String()],
				Extensions:   uploadClient.FileTypeExts[uploadClient.TypeBT],
				MIMEKeywords: uploadClient.FileTypeMimeKeywords[uploadClient.TypeBT.String()],
			},
			`photoshop`: {
				Icon:         uploadClient.FileTypeIcons[uploadClient.TypePhotoshop.String()],
				Extensions:   uploadClient.FileTypeExts[uploadClient.TypePhotoshop],
				MIMEKeywords: uploadClient.FileTypeMimeKeywords[uploadClient.TypePhotoshop.String()],
			},
		},
		MaxSize:           `2M`,
		Icon:              `file-o`,
		AllowedExtensions: AllowedExtensions,
		fileTypes:         map[string]string{},
	}
}

func (c *Config) Reload() error {
	c.Init()
	return nil
}

func (c *Config) SetDefaults() {
	c.Init()
}

func (c *Config) Init() {
	for typeName, ft := range c.FileTypes {
		ft.Init()
		for _, extension := range ft.Extensions {
			c.fileTypes[extension] = typeName
		}
	}
	max, _ := bytes.Parse(c.MaxSize)
	c.maxSizeBytes = int(max)
}

func (c *Config) MaxSizeBytes(typ string) int {
	if len(typ) == 0 {
		return c.maxSizeBytes
	}
	if ft, ok := c.FileTypes[typ]; ok && ft.MaxSizeBytes() > 0 {
		return ft.MaxSizeBytes()
	}
	return c.maxSizeBytes
}

func (c *Config) FileIcon(typ string) string {
	if ft, ok := c.FileTypes[typ]; ok {
		return ft.Icon
	}
	return c.Icon
}

// Extensions 文件类型文件扩展名
func (c *Config) Extensions(typ string) (r []string) {
	if v, ok := c.FileTypes[typ]; ok {
		r = v.Extensions
	}
	return
}

// CheckTypeExtension 检查类型扩展名
func (c *Config) CheckTypeExtension(typ string, extension string) bool {
	extension = strings.TrimPrefix(extension, `.`)
	return com.InSlice(extension, c.Extensions(typ))
}

// DetectType 根据扩展名判断类型
func (c *Config) DetectType(extension string) string {
	extension = strings.TrimPrefix(extension, `.`)
	extension = strings.ToLower(extension)
	if v, ok := c.fileTypes[extension]; ok {
		return v
	}
	mimeType := mime.TypeByExtension(`.` + extension)
	mimeType = strings.SplitN(mimeType, ";", 2)[0]
	for typeK, ft := range c.FileTypes {
		for _, words := range ft.MIMEKeywords {
			if strings.Contains(mimeType, words) {
				return typeK
			}
		}
	}
	return `file`
}

/*
func (c *Config) Register() {
	uploadClient.FileTypeExts = map[uploadClient.FileType][]string{}
	uploadClient.FileTypeMimeKeywords = map[string][]string{}
	for typeName, ft := range c.FileTypes {
		uploadClient.TypeRegister(uploadClient.FileType(typeName), ft.Extensions...)
		if _, ok := uploadClient.FileTypeMimeKeywords[typeName]; !ok {
			uploadClient.FileTypeMimeKeywords[typeName] = ft.MIMEKeywords
		} else {
			for _, kw := range ft.MIMEKeywords {
				if com.InSlice(kw, uploadClient.FileTypeMimeKeywords[typeName]) {
					continue
				}
				uploadClient.FileTypeMimeKeywords[typeName] = append(uploadClient.FileTypeMimeKeywords[typeName], kw)
			}
		}
		if len(ft.Icon) > 0 {
			uploadClient.FileTypeIcons[typeName] = ft.Icon
		}
	}
}
*/
