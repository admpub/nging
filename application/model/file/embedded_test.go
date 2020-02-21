package file_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
	myTesting "github.com/webx-top/echo/testing"
	"github.com/webx-top/echo/testing/test"

	"github.com/admpub/log"
	"github.com/admpub/nging/application/dbschema"
	_ "github.com/admpub/nging/application/handler/manager/file"
	"github.com/admpub/nging/application/library/config"
	mw "github.com/admpub/nging/application/middleware"
	modelFile "github.com/admpub/nging/application/model/file"
)

// TestUpdateEmbedded 测试图片修改
func TestUpdateEmbedded(t *testing.T) {
	ownerID := uint64(1)
	log.Sync()
	config.DefaultCLIConfig.Conf = filepath.Join(os.Getenv("GOPATH"), `src`, `github.com/admpub/nging/config/config.yaml`)
	if err := config.ParseConfig(); err != nil {
		panic(err)
	}
	config.DefaultConfig.SetDebug(true)
	config.DefaultConfig.ConnectedDB()
	e := echo.New()
	e.Use(mw.Tansaction())
	req, resp := myTesting.NewRequestAndResponse(`GET`, `/`)
	ctx := e.NewContext(req, resp)
	ctx.SetTransaction(factory.NewParam())
	tables := []string{}
	for table := range dbschema.DBI.Events {
		tables = append(tables, table)
	}
	echo.Dump(tables)
	test.Contains(t, tables, `user`)
	userM := &dbschema.NgingUser{}
	userM.SetContext(ctx)
	userM.Get(nil, `id`, ownerID)
	if len(userM.Avatar) > 0 {
		userM.Avatar = ``
	} else {
		userM.Avatar = `/public/upload/user/` + fmt.Sprint(ownerID) + `/avatar.jpg`
	}
	if err := userM.Edit(nil, `id`, ownerID); err != nil {
		panic(err)
	}
	m := modelFile.NewFile(ctx)
	m.Get(nil, db.And(
		db.Cond{`table_id`: ownerID},
		db.Cond{`table_name`: `user`},
		db.Cond{`field_name`: `avatar`},
	))
	em := modelFile.NewEmbedded(ctx, m)
	num, err := em.Count(nil, db.And(
		db.Cond{`table_id`: ownerID},
		db.Cond{`table_name`: `user`},
		db.Cond{`field_name`: `avatar`},
	))
	if err != nil && err != db.ErrNoMoreRows {
		panic(err)
	}
	if len(userM.Avatar) == 0 { // 头像删除
		fmt.Println(`头像删除`)
		test.Eq(t, uint(0), m.UsedTimes)
		test.Eq(t, int64(0), num)
	} else { // 添加头像
		fmt.Println(`添加头像`)
		test.Eq(t, uint(1), m.UsedTimes)
		test.Eq(t, int64(1), num)
	}
}
