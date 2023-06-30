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
		DisableInitSet:   true,
		DefaultPolicy:    `accept`,
		TablePrefix:      `nging_`,
		TrustPorts:       []uint16{},
	}
	t := &NFTables{
		base: &Base{
			TableFamily: family,
			cfg:         &cfg,
			NFTables:    biz.New(family, cfg, nil),
		},
	}
	err := t.base.Init()
	if err == nil {
		err = t.base.Do(t.initTableOnly)
		//err = t.ApplyDefault()
	}
	t.base.bin, err = exec.LookPath(`nft`)
	return t, err
}

type NFTables struct {
	base *Base
}

var notNumberRegexp = regexp.MustCompile(`[^\d]+`)

func (a *NFTables) initTableOnly(conn *nftables.Conn) error {
	if err := a.base.ApplyBase(conn); err != nil {
		return err
	}
	return conn.Flush()
}

func (a *NFTables) fullTableName(table string) string {
	if len(a.base.cfg.TablePrefix) == 0 || strings.HasPrefix(table, a.base.cfg.TablePrefix) {
		return table
	}
	return a.base.cfg.TablePrefix + table
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

func (a *NFTables) Clear() error {
	return a.base.Do(func(conn *nftables.Conn) error {
		conn.FlushTable(a.base.TableFilter())
		conn.FlushTable(a.base.TableNAT())
		return conn.Flush()
	})
}

func (a *NFTables) Reset() error {
	return a.base.Do(func(conn *nftables.Conn) error {
		conn.DelTable(a.base.TableFilter())
		conn.DelTable(a.base.TableNAT())
		return conn.Flush()
	})
}

func (a *NFTables) Import(wfwFile string) error {
	return cmdutils.RunCmd(context.Background(), a.base.bin, []string{`-f`, wfwFile}, nil)
}

func (a *NFTables) Export(wfwFile string) error {
	os.MkdirAll(filepath.Dir(wfwFile), os.ModePerm)
	f, err := os.Create(wfwFile)
	if err != nil {
		return err
	}
	defer f.Close()
	err = cmdutils.RunCmd(context.Background(), a.base.bin, []string{`list`, `ruleset`}, f)
	if err != nil {
		return err
	}
	return f.Sync()
}

func (a *NFTables) Insert(rules ...driver.Rule) (err error) {
	return a.base.Do(func(conn *nftables.Conn) error {
		for _, rule := range rules {
			copyRule := rule
			exprs, err := a.ruleFrom(conn, &copyRule)
			if err != nil {
				return err
			}
			ruleTarget, err := a.base.NewRuleTarget(copyRule.Type, copyRule.Direction)
			if err != nil {
				return err
			}
			id := rule.IDBytes()
			ruleData := ruleutils.NewData(id, exprs, 0, uint64(copyRule.Number))
			_, err = ruleTarget.Insert(conn, ruleData)
			if err != nil {
				return err
			}
		}
		return conn.Flush()
	})
}

func (a *NFTables) Append(rules ...driver.Rule) (err error) {
	return a.base.Do(func(conn *nftables.Conn) error {
		for _, rule := range rules {
			copyRule := rule
			exprs, err := a.ruleFrom(conn, &copyRule)
			if err != nil {
				return err
			}
			ruleTarget, err := a.base.NewRuleTarget(copyRule.Type, copyRule.Direction)
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
	return a.base.Do(func(conn *nftables.Conn) error {
		exprs, err := a.ruleFrom(conn, &rule)
		if err != nil {
			return err
		}
		ruleTarget, err := a.base.NewRuleTarget(rule.Type, rule.Direction)
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

func (a *NFTables) DeleteByHandleID(rules ...driver.Rule) (err error) {
	//nft delete rule filter output handle 10
	ctx := context.Background()
	for _, rule := range rules {
		err = cmdutils.RunCmd(ctx, a.base.bin, []string{
			`delete`, `rule`, a.base.getTableFamilyString(), a.fullTableName(rule.Type), rule.Direction,
			`handle`, strconv.FormatUint(uint64(rule.Number), 10),
		}, nil)
		if err != nil {
			return
		}
	}
	return
}

func (a *NFTables) Delete(rules ...driver.Rule) (err error) {
	return a.base.Do(func(conn *nftables.Conn) error {
		for _, rule := range rules {
			ruleTarget, err := a.base.NewRuleTarget(rule.Type, rule.Direction)
			if err != nil {
				return err
			}
			id := rule.IDBytes()
			ruleData := ruleutils.NewData(id, nil, uint64(rule.Number))
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
	err := a.base.Do(func(conn *nftables.Conn) (err error) {
		ruleTarget, err := a.base.NewRuleTarget(rule.Type, rule.Direction)
		if err != nil {
			return err
		}
		id := rule.IDBytes()
		ruleData := ruleutils.NewData(id, nil, uint64(rule.Number))
		exists, err = ruleTarget.Exists(conn, ruleData)
		return
	})
	return exists, err
}

func (a *NFTables) FindPositionByID(table, chain string, id uint) (uint, error) {
	return a.base.FindPositionByID(table, chain, id)
}

func (a *NFTables) Base() *Base {
	return a.base
}
