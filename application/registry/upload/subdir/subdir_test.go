package subdir_test

import (
	"testing"

	"github.com/webx-top/echo"
	myTesting "github.com/webx-top/echo/testing"
	"github.com/webx-top/echo/testing/test"

	_ "github.com/admpub/nging/application/listener/upload"
	"github.com/admpub/nging/application/registry/upload/checker"
	uploadSubdir "github.com/admpub/nging/application/registry/upload/subdir"
	"github.com/admpub/nging/application/registry/upload/table"
)

func TestChecker(t *testing.T) {
	uploadSubdir.Register((&uploadSubdir.SubdirInfo{
		Allowed: true,
		Key:     "test_user",
		Name:    "测试",
	}).SetTable(`test_table`, ``))
	uploadSubdir.CheckerRegister(`test_user`, func(ctx echo.Context, tis table.TableInfoStorer) (subdir string, name string, err error) {
		tis.SetTableID(`test-user-id`)
		tis.SetTableName(`test-user`)
		tis.SetFieldName(`test`)
		return
	}, ``)
	tbl := &table.TableInfo{}
	checkerFn := uploadSubdir.CheckerGet(`test_user`)
	req, resp := myTesting.NewRequestAndResponse(`GET`, `/`)
	e := echo.New()
	ctx := e.NewContext(req, resp)
	subdir, name, err := checkerFn(ctx, tbl)
	println(subdir)
	println(name)
	test.Eq(t, `test-user-id`, tbl.TableID())
	test.Eq(t, `test-user`, tbl.TableName())
	test.Eq(t, `test`, tbl.FieldName())
	if err != nil {
		t.Fatal(err)
	}
}

func TestChecker2(t *testing.T) {
	uploadSubdir.Register((&uploadSubdir.SubdirInfo{
		Allowed: true,
		Key:     "test_user2",
		Name:    "测试2",
	}).SetTable(`test_table2`, `field2`))
	uploadSubdir.CheckerRegister(`test_user2`, func(ctx echo.Context, tis table.TableInfoStorer) (subdir string, name string, err error) {
		tis.SetTableID(`test-user-id2`)
		tis.SetTableName(`test-user2`)
		tis.SetFieldName(`test2`)
		return
	}, `field2`)
	tbl := &table.TableInfo{}
	checkerFn := uploadSubdir.CheckerGet(`test_user2`)
	req, resp := myTesting.NewRequestAndResponse(`GET`, `/`)
	e := echo.New()
	ctx := e.NewContext(req, resp)
	tbl.SetFieldName(`field2`)
	subdir, name, err := checkerFn(ctx, tbl)
	println(subdir)
	println(name)
	test.Eq(t, `test-user-id2`, tbl.TableID())
	test.Eq(t, `test-user2`, tbl.TableName())
	test.Eq(t, `test2`, tbl.FieldName())
	if err != nil {
		t.Fatal(err)
	}
}

func TestChecker3(t *testing.T) {
	tbl := &table.TableInfo{}
	checkerFn := uploadSubdir.CheckerGet(`nging_user.avatar`)
	req, resp := myTesting.NewRequestAndResponse(`GET`, `/`+checker.URLParam(`nging_user.avatar`, `refid`, `0`))
	e := echo.New()
	ctx := e.NewContext(req, resp)
	subdir, name, err := checkerFn(ctx, tbl)
	println(subdir)
	println(name)
	/*
		test.Eq(t, `0`, tbl.TableID())
		test.Eq(t, `user`, tbl.TableName())
		test.Eq(t, `avatar`, tbl.FieldName())
	// */
	if err != nil {
		t.Log(err)
	}
}

func TestThumbSize(t *testing.T) {
	c := uploadSubdir.ThumbSize{
		Width:  400,
		Height: 650,
	}
	test.Eq(t, `400x650`, c.String())
}
