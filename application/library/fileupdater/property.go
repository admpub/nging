package fileupdater

import (
	"database/sql"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
)

// GenUpdater 生成Updater
func GenUpdater(m factory.Model, cond db.Compound) func(event string, content string) error {
	return func(field string, content string) error {
		m.EventOFF()
		err := m.SetField(nil, field, content, cond)
		m.EventON()
		return err
	}
}

// Property 附加属性
type Property struct {
	embedded  sql.NullBool   // 是否为嵌入图片
	seperator sql.NullString // 文件字段中多个文件路径之间的分隔符，空字符串代表为单个文件
	updater   func(event string, content string) error
	exit      bool
}

func NewProperty() *Property {
	return &Property{}
}

// NewPropertyWith 创建并生成带默认Updater的实例
func NewPropertyWith(m factory.Model, cond db.Compound) *Property {
	return NewProperty().GenUpdater(m, cond)
}

func (pro *Property) GenUpdater(m factory.Model, cond db.Compound) *Property {
	pro.updater = GenUpdater(m, cond)
	return pro
}

func (pro *Property) SetUpdater(updater func(field string, content string) error) *Property {
	pro.updater = updater
	return pro
}

func (pro *Property) Updater() func(field string, content string) error {
	return pro.updater
}

func (pro *Property) Embedded() sql.NullBool {
	return pro.embedded
}

func (pro *Property) Seperator() sql.NullString {
	return pro.seperator
}

func (pro *Property) Exit() bool {
	return pro.exit
}

func (pro *Property) SetExit(exit bool) *Property {
	pro.exit = exit
	return pro
}

func (pro *Property) SetEmbedded(on bool) *Property {
	pro.embedded.Bool = on
	pro.embedded.Valid = true
	return pro
}

func (pro *Property) SetSeperator(sep string) *Property {
	pro.seperator.String = sep
	pro.seperator.Valid = true
	return pro
}
