package iptables

import (
	"strconv"
	"strings"

	"github.com/nging-plugins/firewallmanager/application/library/cmdutils"
	"github.com/webx-top/com"
)

func LineCommentParser(findComments []string) func(i uint64, t string) (rowInfo *cmdutils.RowInfo, err error) {
	return func(i uint64, t string) (rowInfo *cmdutils.RowInfo, err error) {
		pos := strings.Index(t, `/* `)
		if pos == -1 {
			return
		}
		part := t[pos+3:]
		pos = strings.Index(part, ` */`)
		if pos == -1 {
			return
		}
		comment := part[0:pos]
		if !com.InSlice(comment, findComments) {
			return
		}
		handleID, err := strconv.ParseUint(strings.SplitN(t, ` `, 2)[0], 10, 64)
		if err != nil {
			return
		}
		rowInfo = &cmdutils.RowInfo{
			RowNo: i,
			Row:   comment,
		}
		rowInfo.Handle.SetValid(handleID)
		return
	}
}
