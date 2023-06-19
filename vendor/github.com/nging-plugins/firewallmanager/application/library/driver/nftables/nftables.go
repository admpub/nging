/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package nftables

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/admpub/nftablesutils"
	"github.com/admpub/nftablesutils/biz"
	ruleutils "github.com/admpub/nftablesutils/rule"
	"github.com/google/nftables"
	"github.com/nging-plugins/firewallmanager/application/library/cmdutils"
	"github.com/nging-plugins/firewallmanager/application/library/driver"
	"github.com/nging-plugins/firewallmanager/application/library/enums"
)

var _ driver.Driver = (*NFTables)(nil)

func New(proto driver.Protocol) (*NFTables, error) {
	var family nftables.TableFamily
	if proto == driver.ProtocolIPv4 {
		family = nftables.TableFamilyIPv4
	} else {
		family = nftables.TableFamilyIPv6
	}
	cfg := biz.Config{
		NetworkNamespace: ``,
		Enabled:          true,
		DefaultPolicy:    `accept`,
		TablePrefix:      `nging_`,
		TrustPorts:       []uint16{},
	}
	t := &NFTables{
		TableFamily: family,
		cfg:         &cfg,
		NFTables:    biz.New(family, cfg, nil),
	}
	err := t.Init()
	if err == nil {
		err = t.Do(t.initTableOnly)
		//err = t.ApplyDefault()
	}
	t.bin, err = exec.LookPath(`nft`)
	return t, err
}

type NFTables struct {
	TableFamily nftables.TableFamily
	cfg         *biz.Config
	bin         string
	*biz.NFTables
}

var notNumberRegexp = regexp.MustCompile(`[^\d]+`)

func (a *NFTables) isIPv4() bool {
	return a.TableFamily == nftables.TableFamilyIPv4
}

func (a *NFTables) initTableOnly(conn *nftables.Conn) error {
	if err := a.ApplyBase(conn); err != nil {
		return err
	}
	return conn.Flush()
}

func (a *NFTables) fullTableName(table string) string {
	if len(a.cfg.TablePrefix) == 0 || strings.HasPrefix(table, a.cfg.TablePrefix) {
		return table
	}
	return a.cfg.TablePrefix + table
}

func (a *NFTables) ruleFrom(c *nftables.Conn, rule *driver.Rule) (args nftablesutils.Exprs, err error) {
	if len(rule.Type) == 0 {
		rule.Type = enums.TableFilter // table
	}
	if len(rule.Direction) == 0 {
		rule.Direction = enums.ChainInput // chain
	} else {
		rule.Direction = rule.Direction
	}
	if rule.Type == enums.TableNAT {
		return a.ruleNATFrom(c, rule)
	}
	return a.ruleFilterFrom(c, rule)
}

func (a *NFTables) Enabled(on bool) error {
	return driver.ErrUnsupported
}

func (a *NFTables) Reset() error {
	return a.NFTables.Do(func(conn *nftables.Conn) error {
		conn.FlushTable(a.NFTables.TableFilter())
		conn.FlushTable(a.NFTables.TableNAT())
		return conn.Flush()
	})
}

func (a *NFTables) Import(wfwFile string) error {
	return cmdutils.RunCmd(context.Background(), a.bin, []string{`-f`, wfwFile}, nil)
}

func (a *NFTables) Export(wfwFile string) error {
	os.MkdirAll(filepath.Dir(wfwFile), os.ModePerm)
	f, err := os.Create(wfwFile)
	if err != nil {
		return err
	}
	defer f.Close()
	err = cmdutils.RunCmd(context.Background(), a.bin, []string{`list`, `ruleset`}, f)
	if err != nil {
		return err
	}
	return f.Sync()
}

func (a *NFTables) ListSets(table, set string, page, limit uint) (rows []cmdutils.RowInfo, hasMore bool, err error) {
	//nft --handle list set test_filter trust_ipset
	return cmdutils.RecvCmdOutputs(page, limit, a.bin, []string{`--handle`, `list`, `set`, a.getTableFamilyString(), table, set}, LineParser)
}

func (a *NFTables) getTableFamilyString() string {
	var family string
	if a.isIPv4() {
		family = `ip`
	} else {
		family = `ip6`
	}
	return family
}

func (a *NFTables) ListChainRules(table, chain string, page, limit uint) (rows []cmdutils.RowInfo, hasMore bool, err error) {
	//nft --handle list chain test_filter input
	return cmdutils.RecvCmdOutputs(page, limit, a.bin, []string{`--handle`, `list`, `chain`, a.getTableFamilyString(), table, chain}, LineParser)
}

func (a *NFTables) NewRuleTarget(table, chain string) (ruleutils.RuleTarget, error) {
	var t *nftables.Table
	var c *nftables.Chain
	if len(a.cfg.TablePrefix) > 0 {
		table = strings.TrimPrefix(table, a.cfg.TablePrefix)
	}
	switch table {
	case `filter`:
		t = a.NFTables.TableFilter()
		switch chain {
		case `INPUT`, `input`:
			c = a.NFTables.ChainInput()
		case `FORWARD`, `forward`:
			c = a.NFTables.ChainForward()
		case `OUTPUT`, `output`:
			c = a.NFTables.ChainOutput()
		default:
			return ruleutils.RuleTarget{}, fmt.Errorf(`%w: %s (table=%v)`, driver.ErrUnsupportedChain, chain, table)
		}
	case `nat`:
		t = a.NFTables.TableNAT()
		switch chain {
		case `PREROUTING`, `prerouting`:
			c = a.NFTables.ChainPrerouting()
		case `POSTROUTING`, `postrouting`:
			c = a.NFTables.ChainPostrouting()
		default:
			return ruleutils.RuleTarget{}, fmt.Errorf(`%w: %s (table=%v)`, driver.ErrUnsupportedChain, chain, table)
		}
	default:
		return ruleutils.RuleTarget{}, fmt.Errorf(`%w: %s`, driver.ErrUnsupportedTable, table)
	}
	return ruleutils.New(t, c), nil
}

func (a *NFTables) Insert(rules ...driver.Rule) (err error) {
	return a.NFTables.Do(func(conn *nftables.Conn) error {
		for _, rule := range rules {
			copyRule := rule
			exprs, err := a.ruleFrom(conn, &copyRule)
			if err != nil {
				return err
			}
			ruleTarget, err := a.NewRuleTarget(copyRule.Type, copyRule.Direction)
			if err != nil {
				return err
			}
			id := rule.IDBytes()
			ruleData := ruleutils.NewData(id, exprs)
			_, err = ruleTarget.Insert(conn, ruleData)
			if err != nil {
				return err
			}
		}
		return conn.Flush()
	})
}

func (a *NFTables) Append(rules ...driver.Rule) (err error) {
	return a.NFTables.Do(func(conn *nftables.Conn) error {
		for _, rule := range rules {
			copyRule := rule
			exprs, err := a.ruleFrom(conn, &copyRule)
			if err != nil {
				return err
			}
			ruleTarget, err := a.NewRuleTarget(copyRule.Type, copyRule.Direction)
			if err != nil {
				return err
			}
			id := rule.IDBytes()
			ruleData := ruleutils.NewData(id, exprs)
			_, err = ruleTarget.Add(conn, ruleData)
			if err != nil {
				return err
			}
		}
		return conn.Flush()
	})
}

func (a *NFTables) AsWhitelist(tableName, chainName string) error {
	// a.cfg.DefaultPolicy = `drop`
	// return a.NFTables.Do(func(conn *nftables.Conn) error {
	// 	conn.FlushTable(a.NFTables.TableFilter())
	// 	// reapply
	// 	return conn.Flush()
	// })
	return driver.ErrUnsupported
}

// Update update rulespec in specified table/chain
func (a *NFTables) Update(rule driver.Rule) error {
	return a.NFTables.Do(func(conn *nftables.Conn) error {
		exprs, err := a.ruleFrom(conn, &rule)
		if err != nil {
			return err
		}
		ruleTarget, err := a.NewRuleTarget(rule.Type, rule.Direction)
		if err != nil {
			return err
		}
		id := rule.IDBytes()
		ruleData := ruleutils.NewData(id, exprs)
		_, err = ruleTarget.Update(conn, ruleData)
		if err != nil {
			return err
		}
		return conn.Flush()
	})
}

func (a *NFTables) DeleteElementInSet(table, set, element string) (err error) {
	//nft delete element global ipv4_ad { 192.168.1.5 }
	//element = com.AddCSlashes(element, ';')
	err = cmdutils.RunCmd(context.Background(), a.bin, []string{
		`delete`, `element`, a.getTableFamilyString(), table, set, `{ ` + element + ` }`,
	}, nil)
	return
}

func (a *NFTables) DeleteElementInSetByHandleID(table, set string, handleID uint64) (err error) {
	err = cmdutils.RunCmd(context.Background(), a.bin, []string{
		`delete`, `element`, a.getTableFamilyString(), table, set,
		`handle`, strconv.FormatUint(handleID, 10),
	}, nil)
	return
}

func (a *NFTables) DeleteSet(table, set string) (err error) {
	//nft delete set global myset
	err = cmdutils.RunCmd(context.Background(), a.bin, []string{
		`delete`, `set`, a.getTableFamilyString(), table, set,
	}, nil)
	return
}

func (a *NFTables) DeleteRuleByHandleID(table, chain string, handleID uint64) (err error) {
	//nft delete rule filter output handle 10
	err = cmdutils.RunCmd(context.Background(), a.bin, []string{
		`delete`, `rule`, a.getTableFamilyString(), table, chain,
		`handle`, strconv.FormatUint(handleID, 10),
	}, nil)
	return
}

func (a *NFTables) DeleteByHandleID(rules ...driver.Rule) (err error) {
	//nft delete rule filter output handle 10
	ctx := context.Background()
	for _, rule := range rules {
		err = cmdutils.RunCmd(ctx, a.bin, []string{
			`delete`, `rule`, a.getTableFamilyString(), a.fullTableName(rule.Type), rule.Direction,
			`handle`, strconv.FormatUint(rule.Number, 10),
		}, nil)
		if err != nil {
			return
		}
	}
	return
}

func (a *NFTables) Delete(rules ...driver.Rule) (err error) {
	return a.NFTables.Do(func(conn *nftables.Conn) error {
		for _, rule := range rules {
			ruleTarget, err := a.NewRuleTarget(rule.Type, rule.Direction)
			if err != nil {
				return err
			}
			id := rule.IDBytes()
			ruleData := ruleutils.NewData(id, nil, rule.Number)
			ok, err := ruleTarget.Delete(conn, ruleData)
			if err != nil {
				return err
			}
			//fmt.Printf("deleted: %s ====================> %b\n", id, ok)
			_ = ok
		}
		return conn.Flush()
	})
}

func (a *NFTables) Exists(rule driver.Rule) (bool, error) {
	var exists bool
	err := a.NFTables.Do(func(conn *nftables.Conn) (err error) {
		exprs, err := a.ruleFrom(conn, &rule)
		if err != nil {
			return err
		}
		ruleTarget, err := a.NewRuleTarget(rule.Type, rule.Direction)
		if err != nil {
			return err
		}
		id := rule.IDBytes()
		ruleData := ruleutils.NewData(id, exprs, rule.Number)
		exists, err = ruleTarget.Exists(conn, ruleData)
		return
	})
	return exists, err
}

func (a *NFTables) Stats(tableName, chainName string) ([]map[string]string, error) {
	var result []map[string]string
	// ruleTarget := a.NewFilterRuleTarget()
	// err := a.NFTables.Do(func(conn *nftables.Conn) error {
	// 	rules, err := ruleTarget.List(conn)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	for _, rule := range rules {
	// 		result = append(result, map[string]string{
	// 			`id`: string(rule.UserData),
	// 		})
	// 	}
	// 	return err
	// })
	return result, driver.ErrUnsupported
}

func (a *NFTables) List(tableName, chainName string) ([]*driver.Rule, error) {
	return nil, driver.ErrUnsupported
}
