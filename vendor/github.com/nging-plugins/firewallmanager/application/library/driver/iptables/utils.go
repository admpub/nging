package iptables

import (
	"strconv"
	"strings"

	"github.com/admpub/go-iptables/iptables"
	"github.com/nging-plugins/firewallmanager/application/library/cmdutils"
	"github.com/nging-plugins/firewallmanager/application/library/enums"
	"github.com/webx-top/com"
)

func LineCommentParser(findComments []string) func(i uint, t string) (rowInfo cmdutils.RowInfo, err error) {
	return func(i uint, t string) (rowInfo cmdutils.RowInfo, err error) {
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
		var handleID uint64
		handleID, err = strconv.ParseUint(strings.SplitN(t, ` `, 2)[0], 10, 0)
		if err != nil {
			return
		}
		rowInfo = cmdutils.RowInfo{
			RowNo: i,
			Row:   comment,
		}
		rowInfo.Handle.SetValid(uint(handleID))
		return
	}
}

func getNgingChain(table string, originalChain string) string {
	var m map[string]string
	if table == enums.TableFilter {
		m = RefFilterChains
	} else {
		m = RefNATChains
	}
	for chain, oriChain := range m {
		if originalChain == oriChain {
			return chain
		}
	}
	return ``
}

func IsExist(err error) bool {
	e, y := err.(*iptables.Error)
	if !y {
		return false
	}
	return !e.IsNotExist()
}
