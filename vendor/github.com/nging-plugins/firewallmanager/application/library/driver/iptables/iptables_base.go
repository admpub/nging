package iptables

import (
	"strconv"

	"github.com/admpub/go-iptables/iptables"
	"github.com/nging-plugins/firewallmanager/application/library/cmdutils"
	"github.com/nging-plugins/firewallmanager/application/library/enums"
)

type Base struct {
	*iptables.IPTables
}

func (a *Base) AsWhitelist(table, chain string) error {
	return a.IPTables.AppendUnique(table, chain, `-j`, enums.TargetReject)
}

func (a *Base) DeleteByPosition(table, chain string, pos uint64) (err error) {
	err = a.IPTables.Delete(table, chain, strconv.FormatUint(pos, 10))
	return
}

func (a *Base) findByComment(table, chain string, findComments ...string) (map[string]uint, error) {
	result := map[string]uint{}
	if len(findComments) == 0 {
		return result, nil
	}
	rows, _, _, err := cmdutils.RecvCmdOutputs(0, uint(len(findComments)),
		iptables.GetIptablesCommand(a.Proto()),
		[]string{
			`-t`, table,
			`-L`, chain,
			`--line-number`,
		}, LineCommentParser(findComments))
	if err != nil {
		return result, err
	}
	for _, row := range rows {
		result[row.Row] = row.GetHandleID()
	}
	return result, nil
}

func (a *Base) Stats(table, chain string) ([]map[string]string, error) {
	return a.IPTables.StatsWithLineNumber(table, chain)
}

func (a *Base) FindPositionByID(table, chain string, id uint) (uint, error) {
	var position uint
	comment := CommentPrefix + strconv.FormatUint(uint64(id), 10)
	nums, err := a.findByComment(table, chain, comment)
	if err == nil {
		position = nums[comment]
	}
	return position, err
}
