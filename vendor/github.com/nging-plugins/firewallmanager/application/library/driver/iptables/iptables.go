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

package iptables

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/admpub/go-iptables/iptables"
	"github.com/admpub/packer"
	"github.com/nging-plugins/firewallmanager/application/library/cmdutils"
	"github.com/nging-plugins/firewallmanager/application/library/driver"
	"github.com/nging-plugins/firewallmanager/application/library/enums"
)

var _ driver.Driver = (*IPTables)(nil)

func New(proto driver.Protocol, autoInstall bool) (*IPTables, error) {
	t := &IPTables{
		IPProtocol: proto,
		base:       &Base{},
	}
	var family iptables.Protocol
	if t.IPProtocol == driver.ProtocolIPv4 {
		family = iptables.ProtocolIPv4
	} else {
		family = iptables.ProtocolIPv6
	}
	var err error
	t.base.IPTables, err = iptables.New(iptables.IPFamily(family))
	if err != nil && autoInstall && errors.Is(err, exec.ErrNotFound) {
		err = packer.Install(`iptables`)
		if err == nil {
			t.base.IPTables, err = iptables.New(iptables.IPFamily(family))
		}
	}
	if err == nil {
		err = t.init()
	}
	return t, err
}

type IPTables struct {
	IPProtocol driver.Protocol
	base       *Base
}

func (a *IPTables) init() error {
	for _, chain := range FilterChains {
		err := a.base.NewChain(enums.TableFilter, chain)
		if err != nil && !IsExist(err) {
			return err
		}
		refChain := RefFilterChains[chain]
		err = a.base.AppendUnique(enums.TableFilter, refChain, `-j`, chain)
		if err != nil {
			return err
		}
	}
	for _, chain := range NATChains {
		err := a.base.NewChain(enums.TableNAT, chain)
		if err != nil && !IsExist(err) {
			return err
		}
		refChain := RefNATChains[chain]
		err = a.base.AppendUnique(enums.TableNAT, refChain, `-j`, chain)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *IPTables) ruleFrom(rule *driver.Rule) ([]string, error) {
	if len(rule.Type) == 0 {
		rule.Type = enums.TableFilter
	}
	if len(rule.Protocol) == 0 {
		rule.Protocol = enums.ProtocolTCP
	}
	if len(rule.Direction) == 0 {
		rule.Direction = enums.ChainInput
	}
	if rule.Type == enums.TableNAT {
		return a.ruleNATFrom(rule)
	}
	return a.ruleFilterFrom(rule)
}

func (a *IPTables) Enabled(on bool) error {
	return driver.ErrUnsupported
}

func (a *IPTables) Clear() error {
	for _, chain := range FilterChains {
		err := a.base.ClearChain(enums.TableFilter, chain)
		if err != nil {
			return err
		}
	}
	for _, chain := range NATChains {
		err := a.base.ClearChain(enums.TableNAT, chain)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *IPTables) Reset() error {
	var err error
	for _, chain := range FilterChains {
		refChain := RefFilterChains[chain]
		err = a.base.DeleteIfExists(enums.TableFilter, refChain, `-j`, chain)
		if err != nil {
			return err
		}
		err = a.base.ClearAndDeleteChain(enums.TableFilter, chain)
		if err != nil {
			return err
		}
	}
	for _, chain := range NATChains {
		refChain := RefNATChains[chain]
		err = a.base.DeleteIfExists(enums.TableNAT, refChain, `-j`, chain)
		if err != nil {
			return err
		}
		err = a.base.ClearAndDeleteChain(enums.TableNAT, chain)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *IPTables) Import(wfwFile string) error {
	var restoreBin string
	switch a.IPProtocol {
	case driver.ProtocolIPv6:
		restoreBin = `ip6tables-restore`
	case driver.ProtocolIPv4:
		fallthrough
	default:
		restoreBin = `iptables-restore`
	}
	f, err := os.Open(wfwFile)
	if err != nil {
		return err
	}
	defer f.Close()
	return cmdutils.RunCmd(context.Background(), restoreBin, nil, nil, f)
}

func (a *IPTables) Export(wfwFile string) error {
	var saveBin string
	switch a.IPProtocol {
	case driver.ProtocolIPv6:
		saveBin = `ip6tables-save`
	case driver.ProtocolIPv4:
		fallthrough
	default:
		saveBin = `iptables-save`
	}
	os.MkdirAll(filepath.Dir(wfwFile), os.ModePerm)
	f, err := os.Create(wfwFile)
	if err != nil {
		return err
	}
	defer f.Close()
	err = cmdutils.RunCmd(context.Background(), saveBin, nil, f)
	if err != nil {
		return err
	}
	return f.Sync()
}

func (a *IPTables) Insert(rules ...driver.Rule) (err error) {
	var existsIndexes map[int]uint
	existsIndexes, err = a.getExistsIndexes(rules)
	if err != nil {
		return
	}
	for index, rule := range rules {
		if _, ok := existsIndexes[index]; ok {
			continue
		}
		copyRule := rule
		var rulespec []string
		rulespec, err = a.ruleFrom(&copyRule)
		if err != nil {
			return
		}
		table := copyRule.Type
		chain := getNgingChain(table, copyRule.Direction)
		err = a.base.Insert(table, chain, int(copyRule.Number), rulespec...)
		if err != nil {
			return
		}
	}
	return err
}

func (a *IPTables) getExistsIndexes(rules []driver.Rule) (map[int]uint, error) {
	comments := map[string]map[string][]string{}
	commentk := map[string]int{}
	exists := map[int]uint{}
	for index, rule := range rules {
		idStr := rule.IDString()
		if len(idStr) > 0 {
			comment := CommentPrefix + idStr
			commentk[comment] = index
			if _, ok := comments[rule.Type]; !ok {
				comments[rule.Type] = map[string][]string{}
			}
			chain := getNgingChain(rule.Type, rule.Direction)
			if _, ok := comments[rule.Type][chain]; !ok {
				comments[rule.Type][chain] = []string{}
			}
			comments[rule.Type][chain] = append(comments[rule.Type][chain], comment)
		}
	}
	if len(comments) > 0 {
		for table, chains := range comments {
			for chain, cmts := range chains {
				nums, err := a.base.findByComment(table, chain, cmts...)
				if err != nil {
					return exists, err
				}
				for comment, num := range nums {
					exists[commentk[comment]] = num
				}
			}
		}
	}
	return exists, nil
}

func (a *IPTables) Append(rules ...driver.Rule) (err error) {
	var existsIndexes map[int]uint
	existsIndexes, err = a.getExistsIndexes(rules)
	if err != nil {
		return
	}
	for index, rule := range rules {
		if _, ok := existsIndexes[index]; ok {
			continue
		}
		copyRule := rule
		var rulespec []string
		rulespec, err = a.ruleFrom(&copyRule)
		if err != nil {
			return
		}
		table := copyRule.Type
		chain := getNgingChain(rule.Type, copyRule.Direction)
		err = a.base.Append(table, chain, rulespec...)
		if err != nil {
			return
		}
	}
	return err
}

// Update update rulespec in specified table/chain
func (a *IPTables) Update(rule driver.Rule) error {
	rulespec, err := a.ruleFrom(&rule)
	if err != nil {
		return err
	}
	table := rule.Type
	chain := getNgingChain(rule.Type, rule.Direction)
	args := []string{"-t", table, "-R", chain}
	idStr := rule.IDString()
	if rule.Number <= 0 && len(idStr) > 0 {
		cmt := CommentPrefix + idStr
		nums, err := a.base.findByComment(table, chain, cmt)
		if err != nil {
			return err
		}
		rule.Number = nums[cmt]
	}
	args = append(args, strconv.FormatUint(uint64(rule.Number), 10))
	cmd := append(args, rulespec...)
	return a.base.Run(cmd...)
}

func (a *IPTables) Delete(rules ...driver.Rule) (err error) {
	for _, rule := range rules {
		copyRule := rule
		var rulespec []string
		if rule.Number > 0 {
			rulespec = append(rulespec, strconv.FormatUint(uint64(rule.Number), 10))
		} else {
			rulespec, err = a.ruleFrom(&copyRule)
			if err != nil {
				return
			}
		}
		table := rule.Type
		chain := getNgingChain(rule.Type, rule.Direction)
		err = a.base.Delete(table, chain, rulespec...)
		if err != nil {
			return
		}
	}
	return err
}

func (a *IPTables) Exists(rule driver.Rule) (bool, error) {
	rulespec, err := a.ruleFrom(&rule)
	if err != nil {
		return false, err
	}
	table := rule.Type
	chain := getNgingChain(rule.Type, rule.Direction)
	return a.base.Exists(table, chain, rulespec...)
}

func (a *IPTables) AsWhitelist(table, chain string) error {
	chain = getNgingChain(table, chain)
	return a.base.AppendUnique(table, chain, `-j`, enums.TargetReject)
}

func (a *IPTables) FindPositionByID(table, chain string, id uint) (uint, error) {
	chain = getNgingChain(table, chain)
	return a.base.FindPositionByID(table, chain, id)
}

func (a *IPTables) Base() *Base {
	return a.base
}
