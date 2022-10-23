package sqlquery

import (
	"fmt"
	"sort"

	"github.com/admpub/nging/v4/application/cmd"
	"github.com/admpub/nging/v4/application/dbschema"
	"github.com/admpub/nging/v4/application/library/common"
	"github.com/spf13/cobra"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
)

func init() {
	cmd.TestSuiteRegister(`sqlquery`, GetRow)
	cmd.TestSuiteRegister(`dbschemas`, GetDbSchemas)
	cmd.TestSuiteRegister(`args`, GetCmdArgs)
}

func GetRow(cmd *cobra.Command, args []string) error {
	ctx := defaults.NewMockContext()
	row, err := common.NewSQLQuery(ctx).GetRow("SELECT * FROM nging_user WHERE id > 0")
	if err != nil {
		panic(err)
	}
	fmt.Print(`nging_user.created =========> `)
	echo.Dump(row.Timestamp(`created`).Format(`2006-01-02 15:04:05`))

	m, err := common.NewSQLQuery(ctx).GetModel(`NgingUser`, `id >`, 0)
	if err != nil {
		panic(err)
	}
	fmt.Print("\nNgingUser =========> ")
	echo.Dump(m)

	rows, err := common.NewSQLQuery(ctx).Limit(2).GetModels(`NgingUser`, `id >`, 0)
	if err != nil {
		panic(err)
	}
	fmt.Print("\nNgingUser list limit 2 =========> ")
	echo.Dump(rows)

	fmt.Print("\nnging_user GetRows =========> ")
	list, err := common.NewSQLQuery(ctx).Limit(2).GetRows("SELECT * FROM nging_user WHERE id > 0")
	if err != nil {
		panic(err)
	}
	echo.Dump(list)
	return err
}

// GetDbSchemas 打印表结构体名称列表
func GetDbSchemas(cmd *cobra.Command, args []string) error {
	list := make([]string, 0, len(dbschema.DBI.Models))
	for structName := range dbschema.DBI.Models {
		list = append(list, structName)
	}
	sort.Strings(list)
	echo.Dump(list)
	return nil
}

// GetCmdArgs 获取命令行参数（for testing）
func GetCmdArgs(cmd *cobra.Command, args []string) error {
	echo.Dump(args)
	return nil
}
