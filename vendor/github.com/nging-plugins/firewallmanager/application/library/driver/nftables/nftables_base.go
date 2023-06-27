package nftables

import (
	"context"
	"fmt"
	"strconv"

	"github.com/admpub/nftablesutils/biz"
	ruleutils "github.com/admpub/nftablesutils/rule"
	"github.com/google/nftables"
	"github.com/nging-plugins/firewallmanager/application/library/cmdutils"
	"github.com/nging-plugins/firewallmanager/application/library/driver"
)

type Base struct {
	TableFamily nftables.TableFamily
	cfg         *biz.Config
	bin         string
	*biz.NFTables
}

func (a *Base) isIPv4() bool {
	return a.TableFamily == nftables.TableFamilyIPv4
}

func (a *Base) ListSets(table, set string, page, limit uint) (rows []cmdutils.RowInfo, hasMore bool, err error) {
	//nft --handle list set test_filter trust_ipset
	return cmdutils.RecvCmdOutputs(page, limit, a.bin, []string{`--handle`, `list`, `set`, a.getTableFamilyString(), table, set}, LineParser)
}

func (a *Base) getTableFamilyString() string {
	var family string
	if a.isIPv4() {
		family = `ip`
	} else {
		family = `ip6`
	}
	return family
}

func (a *Base) ListChainRules(table, chain string, page, limit uint) (rows []cmdutils.RowInfo, hasMore bool, err error) {
	//nft --handle list chain test_filter input
	return cmdutils.RecvCmdOutputs(page, limit, a.bin, []string{`--handle`, `list`, `chain`, a.getTableFamilyString(), table, chain}, LineParser)
}

func (a *Base) DeleteElementInSet(table, set, element string) (err error) {
	//nft delete element global ipv4_ad { 192.168.1.5 }
	//element = com.AddCSlashes(element, ';')
	err = cmdutils.RunCmd(context.Background(), a.bin, []string{
		`delete`, `element`, a.getTableFamilyString(), table, set, `{ ` + element + ` }`,
	}, nil)
	return
}

func (a *Base) DeleteElementInSetByHandleID(table, set string, handleID uint64) (err error) {
	err = cmdutils.RunCmd(context.Background(), a.bin, []string{
		`delete`, `element`, a.getTableFamilyString(), table, set,
		`handle`, strconv.FormatUint(handleID, 10),
	}, nil)
	return
}

func (a *Base) DeleteSet(table, set string) (err error) {
	//nft delete set global myset
	err = cmdutils.RunCmd(context.Background(), a.bin, []string{
		`delete`, `set`, a.getTableFamilyString(), table, set,
	}, nil)
	return
}

func (a *Base) DeleteRuleByHandleID(table, chain string, handleID uint64) (err error) {
	//nft delete rule filter output handle 10
	err = cmdutils.RunCmd(context.Background(), a.bin, []string{
		`delete`, `rule`, a.getTableFamilyString(), table, chain,
		`handle`, strconv.FormatUint(handleID, 10),
	}, nil)
	return
}

func (a *Base) NewRuleTarget(table, chain string) (ruleutils.RuleTarget, error) {
	var t *nftables.Table
	var c *nftables.Chain
	switch table {
	case `filter`:
		t = a.TableFilter()
		switch chain {
		case `INPUT`, `input`:
			c = a.ChainInput()
		case `FORWARD`, `forward`:
			c = a.ChainForward()
		case `OUTPUT`, `output`:
			c = a.ChainOutput()
		default:
			return ruleutils.RuleTarget{}, fmt.Errorf(`%w: %s (table=%v)`, driver.ErrUnsupportedChain, chain, table)
		}
	case `nat`:
		t = a.TableNAT()
		switch chain {
		case `PREROUTING`, `prerouting`:
			c = a.ChainPrerouting()
		case `POSTROUTING`, `postrouting`:
			c = a.ChainPostrouting()
		default:
			return ruleutils.RuleTarget{}, fmt.Errorf(`%w: %s (table=%v)`, driver.ErrUnsupportedChain, chain, table)
		}
	default:
		return ruleutils.RuleTarget{}, fmt.Errorf(`%w: %s`, driver.ErrUnsupportedTable, table)
	}
	return ruleutils.New(t, c), nil
}

func (a *Base) FindPositionByID(table, chain string, id uint) (uint64, error) {
	var position uint64
	err := a.NFTables.Do(func(conn *nftables.Conn) (err error) {
		ruleTarget, err := a.NewRuleTarget(table, chain)
		if err != nil {
			return err
		}
		s := strconv.FormatUint(uint64(id), 10)
		ruleData := ruleutils.NewData([]byte(s), nil, 0)
		rule, err := ruleTarget.FindRuleByID(conn, ruleData)
		if err != nil {
			return err
		}
		position = rule.Position
		return nil
	})
	return position, err
}
