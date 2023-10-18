package iptables

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/admpub/go-iptables/iptables"
	"github.com/nging-plugins/firewallmanager/application/library/cmdutils"
	"github.com/nging-plugins/firewallmanager/application/library/enums"
)

type Base struct {
	blackListSetName string
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
		a.GetExeclutor(),
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

func (a *Base) AttachSet(set string, action string) error {
	rulespec := []string{"-j", action, "-m", "set", "--match-set", set, "src"}
	cmt := &ModuleComment{
		Comment: `Nging` + set,
	}
	rulespec = append(rulespec, cmt.Args()...)
	if err := a.AppendUnique(enums.TableFilter, FilterChainInput, rulespec...); err != nil {
		return fmt.Errorf(`failed to create %s entry for set "%s": %w`, a.GetExeclutor(), set, err)
	}
	return nil
}

func (a *Base) CreateSet(ctx context.Context, set string) error {
	err := cmdutils.RunCmd(ctx, "ipset", []string{`list`, set}, os.Stdout)
	if err == nil {
		return nil
	}
	args := []string{"create", set, "hash:ip"}
	if a.Proto() == iptables.ProtocolIPv6 {
		args = append(args, "family", "inet6")
	}
	args = append(args, "timeout", "0")
	return cmdutils.RunCmd(ctx, "ipset", args, os.Stdout)
}

func (a *Base) CreateBlackListSet() error {
	err := a.CreateSet(context.Background(), a.blackListSetName)
	if err != nil {
		return err
	}
	return a.AttachBlackListSet()
}

func (a *Base) AttachBlackListSet() error {
	return a.AttachSet(a.blackListSetName, enums.TargetDrop)
}

func (a *Base) RemoveBlackListSet() error {
	return a.RemoveSet(a.blackListSetName, enums.TargetDrop)
}

func (a *Base) RemoveSet(set string, action string) error {
	err := a.DeleteIfExists(enums.TableFilter, FilterChainInput, `-j`, action, "-m", "set", "--match-set", set, "src")
	if err != nil {
		return err
	}
	ctx := context.Background()
	args := []string{"destroy", set}
	return cmdutils.RunCmd(ctx, "ipset", args, os.Stdout)
}

func (a *Base) AddToBlacklistSet(ips []net.IP, d time.Duration) error {
	return a.AddToSet(a.blackListSetName, ips, d)
}

func (a *Base) AddToSet(set string, ips []net.IP, d time.Duration) error {
	var err error
	ctx := context.Background()
	dur := fmt.Sprint(d.Seconds())
	for _, ip := range ips {
		ipStr := ip.String()
		args := []string{"test", set, ipStr}
		if err = cmdutils.RunCmd(ctx, "ipset", args, os.Stdout); err != nil {
			args = []string{"add", set, ipStr, "timeout", dur}
			err = cmdutils.RunCmd(ctx, "ipset", args, os.Stdout)
		}
		if err != nil {
			return err
		}
	}
	return err
}

func (a *Base) GetExeclutor() string {
	return iptables.GetIptablesCommand(a.Proto())
}
