// Package rule A library for managing nftables rules
package rule

import (
	"bytes"
	"fmt"

	"github.com/google/nftables"
)

// RuleTarget represents a location to manipulate nftables rules
type RuleTarget struct {
	table *nftables.Table
	chain *nftables.Chain
}

// Create a new location to manipulate nftables rules
func New(table *nftables.Table, chain *nftables.Chain) RuleTarget {
	return RuleTarget{
		table: table,
		chain: chain,
	}
}

// Add a rule with a given ID to a specific table and chain, returns true if the rule was added
func (r *RuleTarget) Add(c *nftables.Conn, ruleData RuleData) (bool, error) {
	exists, err := r.Exists(c, ruleData)
	if err != nil || exists {
		return false, err
	}

	add(c, r.table, r.chain, ruleData)
	return true, nil
}

func (r *RuleTarget) Insert(c *nftables.Conn, ruleData RuleData) (bool, error) {
	exists, err := r.Exists(c, ruleData)
	if err != nil || exists {
		return false, err
	}

	insert(c, r.table, r.chain, ruleData)
	return true, nil
}

func add(c *nftables.Conn, table *nftables.Table, chain *nftables.Chain, ruleData RuleData) {
	r := ruleData.ToRule(table, chain)
	c.AddRule(&r)
}

func insert(c *nftables.Conn, table *nftables.Table, chain *nftables.Chain, ruleData RuleData) {
	r := ruleData.ToRule(table, chain)
	c.InsertRule(&r)
}

// Delete a rule with a given ID from a specific table and chain, returns true if the rule was deleted
func (r *RuleTarget) Delete(c *nftables.Conn, ruleData RuleData) (bool, error) {
	rule, err := r.FindRuleByID(c, ruleData)
	if err != nil || rule == nil {
		return false, err
	}

	if err := c.DelRule(rule); err != nil {
		return false, err
	}

	return true, nil
}

func (r *RuleTarget) FindRuleByID(c *nftables.Conn, ruleData RuleData) (*nftables.Rule, error) {
	rules, err := c.GetRules(r.table, r.chain)
	if err != nil {
		return nil, err
	}

	rule := findRuleByID(ruleData.ID, rules, ruleData.Handle)

	if rule.Table == nil {
		// if the rule we get back is empty (the final return in findRuleByID) we didn't find it
		return nil, nil
	}

	return rule, nil
}

// Determine if a rule with a given ID exists in a specific table and chain
func (r *RuleTarget) Exists(c *nftables.Conn, ruleData RuleData) (bool, error) {
	rule, err := r.FindRuleByID(c, ruleData)
	if err != nil || rule == nil {
		return false, err
	}
	return true, nil
}

func (r *RuleTarget) Update(c *nftables.Conn, ruleData RuleData) (bool, error) {
	rule, err := r.FindRuleByID(c, ruleData)
	if err != nil || rule == nil {
		return false, err
	}

	ruleNew := ruleData.ToRule(rule.Table, rule.Chain)
	ruleNew.Handle = rule.Handle
	c.ReplaceRule(&ruleNew)
	return true, nil
}

// Compare existing and incoming rule IDs adding/removing the difference
//
// First return value is true if the number of rules has changed, false if there were no updates. The second
// and third return values indicate the number of rules added or removed, respectively.
func (r *RuleTarget) UpdateAll(c *nftables.Conn, rules []RuleData) (bool, int, int, error) {
	var modified bool
	existingRules, err := c.GetRules(r.table, r.chain)
	if err != nil {
		return false, 0, 0, fmt.Errorf("error getting existing rules for update: %v", err)
	}

	addRDList, removeRDList := genRuleDelta(existingRules, rules)

	if len(removeRDList) > 0 {
		for _, rule := range removeRDList {
			err := c.DelRule(rule)
			if err != nil {
				return false, 0, 0, err
			}
			modified = true
		}
	}

	if len(addRDList) > 0 {
		for _, rule := range addRDList {
			add(c, r.table, r.chain, rule)
			modified = true
		}
	}

	return modified, len(addRDList), len(removeRDList), nil
}

// Get the nftables table and chain associated with this RuleTarget
func (r *RuleTarget) GetTableAndChain() (*nftables.Table, *nftables.Chain) {
	return r.table, r.chain
}

func (r *RuleTarget) List(c *nftables.Conn) ([]*nftables.Rule, error) {
	return c.GetRules(r.table, r.chain)
}

func genRuleDelta(existingRules []*nftables.Rule, newRules []RuleData) (add []RuleData, remove []*nftables.Rule) {
	existingRuleMap := make(map[string]*nftables.Rule)
	for _, existingRule := range existingRules {
		existingRuleMap[string(existingRule.UserData)] = existingRule
	}

	for _, ruleData := range newRules {
		if _, exists := existingRuleMap[string(ruleData.ID)]; !exists {
			add = append(add, ruleData)
		} else {
			delete(existingRuleMap, string(ruleData.ID))
		}
	}

	for _, v := range existingRuleMap {
		remove = append(remove, v)
	}

	return
}

func findRuleByID(id []byte, rules []*nftables.Rule, handleID uint64) *nftables.Rule {
	for _, rule := range rules {
		if handleID > 0 {
			if rule.Handle == handleID {
				return rule
			}
		} else if bytes.Equal(rule.UserData, id) {
			return rule
		}
	}
	return &nftables.Rule{}
}
