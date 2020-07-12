package sqlquery

import (
	"github.com/admpub/nging/application/cmd"
	"github.com/admpub/nging/application/library/common"
	"github.com/spf13/cobra"
	"github.com/webx-top/echo"
)

func init() {
	cmd.TestSuiteRegister(`sqlquery`, GetRow)
}

func GetRow(cmd *cobra.Command, args []string) error {
	row, err := common.SQLQuery().GetRow("SELECT * FROM nging_user WHERE id > 0")
	if err != nil {
		panic(err)
	}
	echo.Dump(row.Timestamp(`created`).Format(`2006-01-02 15:04:05`))
	return err
}
