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
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/admpub/go-iptables/iptables"
	parser "github.com/admpub/iptables_parser"
	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/library/errorslice"
	"github.com/admpub/packer"
	"github.com/nging-plugins/firewallmanager/application/library/cmdutils"
	"github.com/nging-plugins/firewallmanager/application/library/driver"
	"github.com/nging-plugins/firewallmanager/application/library/enums"
	"github.com/webx-top/echo/param"
)

var _ driver.Driver = (*IPTables)(nil)

func New(proto driver.Protocol, autoInstall bool) (*IPTables, error) {
	t := &IPTables{
		IPProtocol: proto,
	}
	var family iptables.Protocol
	if t.IPProtocol == driver.ProtocolIPv4 {
		family = iptables.ProtocolIPv4
	} else {
		family = iptables.ProtocolIPv6
	}
	var err error
	t.IPTables, err = iptables.New(iptables.IPFamily(family))
	if err != nil && autoInstall && errors.Is(err, exec.ErrNotFound) {
		err = packer.Install(`iptables`)
		if err == nil {
			t.IPTables, err = iptables.New(iptables.IPFamily(family))
		}
	}
	return t, err
}

type IPTables struct {
	IPProtocol driver.Protocol
	*iptables.IPTables
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

func (a *IPTables) Reset() error {
	return driver.ErrUnsupported
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
	var existsIndexes map[int]uint64
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
		chain := copyRule.Direction
		err = a.IPTables.InsertUnique(table, chain, int(copyRule.Number), rulespec...)
		if err != nil {
			return
		}
	}
	return err
}

func (a *IPTables) getExistsIndexes(rules []driver.Rule) (map[int]uint64, error) {
	comments := map[string]map[string][]string{}
	commentk := map[string]int{}
	for index, rule := range rules {
		if rule.ID > 0 {
			comment := CommentPrefix + param.AsString(rule.ID)
			commentk[comment] = index
			if _, ok := comments[rule.Type]; !ok {
				comments[rule.Type] = map[string][]string{}
			}
			if _, ok := comments[rule.Type][rule.Direction]; !ok {
				comments[rule.Type][rule.Direction] = []string{}
			}
			comments[rule.Type][rule.Direction] = append(comments[rule.Type][rule.Direction], comment)
		}
	}
	exists := map[int]uint64{}
	if len(comments) > 0 {
		for table, chains := range comments {
			for chain, cmts := range chains {
				nums, err := a.findByComment(table, chain, cmts...)
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
	var existsIndexes map[int]uint64
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
		chain := copyRule.Direction
		err = a.IPTables.AppendUnique(table, chain, rulespec...)
		if err != nil {
			return
		}
	}
	return err
}

func (a *IPTables) AsWhitelist(table, chain string) error {
	return a.IPTables.AppendUnique(table, chain, `-j`, enums.TargetReject)
}

// Update update rulespec in specified table/chain
func (a *IPTables) Update(rule driver.Rule) error {
	rulespec, err := a.ruleFrom(&rule)
	if err != nil {
		return err
	}
	table := rule.Type
	chain := rule.Direction
	args := []string{"-t", table, "-R", chain}
	if rule.Number <= 0 && rule.ID > 0 {
		cmt := CommentPrefix + param.AsString(rule.ID)
		nums, err := a.findByComment(table, chain, cmt)
		if err != nil {
			return err
		}
		rule.Number = nums[cmt]
	}
	args = append(args, strconv.FormatUint(rule.Number, 10))
	cmd := append(args, rulespec...)
	return a.IPTables.Run(cmd...)
}

func (a *IPTables) Delete(rules ...driver.Rule) (err error) {
	for _, rule := range rules {
		copyRule := rule
		var rulespec []string
		if rule.Number > 0 {
			rulespec = append(rulespec, strconv.FormatUint(rule.Number, 10))
		} else {
			rulespec, err = a.ruleFrom(&copyRule)
			if err != nil {
				return
			}
		}
		table := rule.Type
		chain := rule.Direction
		err = a.IPTables.Delete(table, chain, rulespec...)
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
	chain := rule.Direction
	return a.IPTables.Exists(table, chain, rulespec...)
}

func (a *IPTables) findByComment(table, chain string, findComments ...string) (map[string]uint64, error) {
	result := map[string]uint64{}
	if len(findComments) == 0 {
		return result, nil
	}
	rows, _, err := cmdutils.RecvCmdOutputs(1, uint(len(findComments)+1),
		iptables.GetIptablesCommand(a.Proto()),
		[]string{
			`-t`, table,
			`-L`, chain,
			`--line-number`,
		}, LineCommentParser(findComments))
	if err != nil {
		return result, nil
	}
	for _, row := range rows {
		result[row.Row] = row.GetHandleID()
	}
	return result, nil
}

func (a *IPTables) Stats(table, chain string) ([]map[string]string, error) {
	return a.IPTables.StatsWithLineNumber(table, chain)
}

func (a *IPTables) List(table, chain string) ([]*driver.Rule, error) {
	rows, err := a.IPTables.List(table, chain)
	if err != nil {
		return nil, err
	}
	errs := errorslice.New()
	var rules []*driver.Rule
	var ipVersion string
	switch a.IPProtocol {
	case driver.ProtocolIPv6:
		ipVersion = `6`
	case driver.ProtocolIPv4:
		fallthrough
	default:
		ipVersion = `4`
	}
	for _, row := range rows {
		tr, err := parser.NewFromString(row)
		if err != nil {
			err = fmt.Errorf("[iptables] failed to parse rule: %s: %v", row, err)
			errs.Add(err)
			continue
		}
		//pp.Println(tr)
		rule := &driver.Rule{Type: table, Direction: chain, IPVersion: ipVersion}
		switch r := tr.(type) {
		case parser.Rule:
			log.Debugf("[iptables] rule parsed: %v", r)
			rule.Direction = r.Chain
			if r.Source != nil {
				rule.RemoteIP = r.Source.Value.String()
				if r.Source.Not {
					rule.RemoteIP = `!` + rule.RemoteIP
				}
			}
			if r.Destination != nil {
				rule.LocalIP = r.Destination.Value.String()
				if r.Destination.Not {
					rule.LocalIP = `!` + rule.LocalIP
				}
			}
			if r.Protocol != nil {
				rule.Protocol = r.Protocol.Value
				if r.Protocol.Not {
					rule.Protocol = `!` + rule.Protocol
				}
			}
			if r.Jump != nil {
				rule.Action = r.Jump.Name
			}
			for _, match := range r.Matches {
				for flagKey, flagValue := range match.Flags {
					switch flagKey {
					case `destination-port`:
						rule.LocalPort = strings.Join(flagValue.Values, ` `)
					case `source-port`:
						rule.RemotePort = strings.Join(flagValue.Values, ` `)
					}
				}
			}
		case parser.Policy:
			log.Debugf("[iptables] policy parsed: %v", r)
			// if r.UserDefined == nil || !*r.UserDefined {
			// 	continue
			// }
			rule.Action = r.Action
			rule.Direction = r.Chain
		// case parser.Comment:
		// case parser.Header:
		default:
			log.Debugf("[iptables] something else happend: %v", r)
		}
		rules = append(rules, rule)
	}
	return rules, errs.ToError()
}
