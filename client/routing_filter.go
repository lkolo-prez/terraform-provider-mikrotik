package client

import (
	"fmt"

	"github.com/go-routeros/routeros/v3"
)

// RoutingFilterRule represents a routing filter rule in RouterOS v7
// Path: /routing/filter/rule
//
// RouterOS v7 completely redesigned the routing filter system with a new rule-based syntax.
// Rules use an if-then structure with powerful matching and action capabilities.
//
// Example rule syntax:
//   if (dst == 0.0.0.0/0) { reject }
//   if (dst in 10.0.0.0/8 && bgp-communities includes 65001:100) { accept }
//   if (bgp-communities includes 65001:200) { set bgp-local-pref 200; accept }
type RoutingFilterRule struct {
	Id       string `mikrotik:".id" codegen:"id,mikrotikID,read"`
	Chain    string `mikrotik:"chain" codegen:"chain,required"`
	Rule     string `mikrotik:"rule" codegen:"rule,required"`
	Disabled bool   `mikrotik:"disabled" codegen:"disabled"`
	Comment  string `mikrotik:"comment" codegen:"comment"`
	Invalid  bool   `mikrotik:"invalid" codegen:"invalid,computed"`
	Dynamic  bool   `mikrotik:"dynamic" codegen:"dynamic,computed"`
}

// Implement Resource interface for RoutingFilterRule
var _ Resource = (*RoutingFilterRule)(nil)

func (r *RoutingFilterRule) IDField() string {
	return ".id"
}

func (r *RoutingFilterRule) ID() string {
	return r.Id
}

func (r *RoutingFilterRule) SetID(id string) {
	r.Id = id
}

func (r *RoutingFilterRule) AfterAddHook(reply *routeros.Reply) {
	r.Id = reply.Done.Map["ret"]
}

func (r *RoutingFilterRule) FindField() string {
	return ".id"
}

func (r *RoutingFilterRule) FindFieldValue() string {
	return r.Id
}

func (r *RoutingFilterRule) DeleteField() string {
	return ".id"
}

func (r *RoutingFilterRule) DeleteFieldValue() string {
	return r.Id
}

func (r *RoutingFilterRule) ActionToCommand(action Action) string {
	return map[Action]string{
		Add:    "/routing/filter/rule/add",
		Find:   "/routing/filter/rule/print",
		Update: "/routing/filter/rule/set",
		Delete: "/routing/filter/rule/remove",
	}[action]
}

// RoutingFilterChain represents a routing filter chain
// Path: /routing/filter/chain
//
// Chains group multiple filter rules together and can be referenced by
// BGP connections, OSPF instances, and other routing protocols.
type RoutingFilterChain struct {
	Id            string `mikrotik:".id" codegen:"id,mikrotikID,read"`
	Name          string `mikrotik:"name" codegen:"name,required,unique_key"`
	DynamicChain  bool   `mikrotik:"dynamic" codegen:"dynamic"`
	Disabled      bool   `mikrotik:"disabled" codegen:"disabled"`
	Comment       string `mikrotik:"comment" codegen:"comment"`
}

// Implement Resource interface for RoutingFilterChain
var _ Resource = (*RoutingFilterChain)(nil)

func (r *RoutingFilterChain) IDField() string {
	return ".id"
}

func (r *RoutingFilterChain) ID() string {
	return r.Id
}

func (r *RoutingFilterChain) SetID(id string) {
	r.Id = id
}

func (r *RoutingFilterChain) AfterAddHook(reply *routeros.Reply) {
	r.Id = reply.Done.Map["ret"]
}

func (r *RoutingFilterChain) FindField() string {
	return "name"
}

func (r *RoutingFilterChain) FindFieldValue() string {
	return r.Name
}

func (r *RoutingFilterChain) DeleteField() string {
	return "name"
}

func (r *RoutingFilterChain) DeleteFieldValue() string {
	return r.Name
}

func (r *RoutingFilterChain) ActionToCommand(action Action) string {
	return map[Action]string{
		Add:    "/routing/filter/chain/add",
		Find:   "/routing/filter/chain/print",
		Update: "/routing/filter/chain/set",
		Delete: "/routing/filter/chain/remove",
	}[action]
}

// RoutingFilterSelectChain represents a routing filter select chain
// Path: /routing/filter/select-chain
//
// Select chains are used to choose which chain to execute based on conditions.
// This allows for more complex filtering logic.
type RoutingFilterSelectChain struct {
	Id       string `mikrotik:".id" codegen:"id,mikrotikID,read"`
	Name     string `mikrotik:"name" codegen:"name,required,unique_key"`
	Chain    string `mikrotik:"chain" codegen:"chain"`
	Disabled bool   `mikrotik:"disabled" codegen:"disabled"`
	Comment  string `mikrotik:"comment" codegen:"comment"`
}

// Implement Resource interface for RoutingFilterSelectChain
var _ Resource = (*RoutingFilterSelectChain)(nil)

func (r *RoutingFilterSelectChain) IDField() string {
	return ".id"
}

func (r *RoutingFilterSelectChain) ID() string {
	return r.Id
}

func (r *RoutingFilterSelectChain) SetID(id string) {
	r.Id = id
}

func (r *RoutingFilterSelectChain) AfterAddHook(reply *routeros.Reply) {
	r.Id = reply.Done.Map["ret"]
}

func (r *RoutingFilterSelectChain) FindField() string {
	return "name"
}

func (r *RoutingFilterSelectChain) FindFieldValue() string {
	return r.Name
}

func (r *RoutingFilterSelectChain) DeleteField() string {
	return "name"
}

func (r *RoutingFilterSelectChain) DeleteFieldValue() string {
	return r.Name
}

func (r *RoutingFilterSelectChain) ActionToCommand(action Action) string {
	return map[Action]string{
		Add:    "/routing/filter/select-chain/add",
		Find:   "/routing/filter/select-chain/print",
		Update: "/routing/filter/select-chain/set",
		Delete: "/routing/filter/select-chain/remove",
	}[action]
}

// Wrapper functions for RoutingFilterRule using generic CRUD

func (c *Mikrotik) AddRoutingFilterRule(rule *RoutingFilterRule) (*RoutingFilterRule, error) {
	res, err := c.Add(rule)
	if err != nil {
		return nil, err
	}
	return res.(*RoutingFilterRule), nil
}

func (c *Mikrotik) FindRoutingFilterRule(chain string) (*RoutingFilterRule, error) {
	// This function finds by chain field, not by ID
	// We need to use direct API call since generic Find() uses FindField
	cmd := []string{"/routing/filter/rule/print", "?chain=" + chain}
	res, err := c.connection.Run(cmd...)
	if err != nil {
		return nil, fmt.Errorf("failed to find routing filter rule: %w", err)
	}
	
	if len(res.Re) == 0 {
		return nil, fmt.Errorf("routing filter rule not found in chain: %s", chain)
	}
	
	var rule RoutingFilterRule
	if err := Unmarshal(*res, &rule); err != nil {
		return nil, fmt.Errorf("failed to unmarshal routing filter rule: %w", err)
	}
	
	return &rule, nil
}

func (c *Mikrotik) FindRoutingFilterRuleById(id string) (*RoutingFilterRule, error) {
	rule := &RoutingFilterRule{Id: id}
	cmd := []string{rule.ActionToCommand(Find), "?.id=" + id}
	res, err := c.connection.Run(cmd...)
	if err != nil {
		return nil, fmt.Errorf("failed to find routing filter rule by id: %w", err)
	}
	
	if len(res.Re) == 0 {
		return nil, fmt.Errorf("routing filter rule not found with id: %s", id)
	}
	
	if err := Unmarshal(*res, rule); err != nil {
		return nil, fmt.Errorf("failed to unmarshal routing filter rule: %w", err)
	}
	
	return rule, nil
}

func (c *Mikrotik) UpdateRoutingFilterRule(rule *RoutingFilterRule) (*RoutingFilterRule, error) {
	res, err := c.Update(rule)
	if err != nil {
		return nil, err
	}
	return res.(*RoutingFilterRule), nil
}

func (c *Mikrotik) DeleteRoutingFilterRule(id string) error {
	rule := &RoutingFilterRule{Id: id}
	return c.Delete(rule)
}

// Wrapper functions for RoutingFilterChain using generic CRUD

func (c *Mikrotik) AddRoutingFilterChain(chain *RoutingFilterChain) (*RoutingFilterChain, error) {
	res, err := c.Add(chain)
	if err != nil {
		return nil, err
	}
	return res.(*RoutingFilterChain), nil
}

func (c *Mikrotik) FindRoutingFilterChain(name string) (*RoutingFilterChain, error) {
	chain := &RoutingFilterChain{Name: name}
	res, err := c.Find(chain)
	if err != nil {
		return nil, err
	}
	return res.(*RoutingFilterChain), nil
}

func (c *Mikrotik) UpdateRoutingFilterChain(chain *RoutingFilterChain) (*RoutingFilterChain, error) {
	res, err := c.Update(chain)
	if err != nil {
		return nil, err
	}
	return res.(*RoutingFilterChain), nil
}

func (c *Mikrotik) DeleteRoutingFilterChain(name string) error {
	chain := &RoutingFilterChain{Name: name}
	return c.Delete(chain)
}

// Wrapper functions for RoutingFilterSelectChain using generic CRUD

func (c *Mikrotik) AddRoutingFilterSelectChain(selectChain *RoutingFilterSelectChain) (*RoutingFilterSelectChain, error) {
	res, err := c.Add(selectChain)
	if err != nil {
		return nil, err
	}
	return res.(*RoutingFilterSelectChain), nil
}

func (c *Mikrotik) FindRoutingFilterSelectChain(name string) (*RoutingFilterSelectChain, error) {
	selectChain := &RoutingFilterSelectChain{Name: name}
	res, err := c.Find(selectChain)
	if err != nil {
		return nil, err
	}
	return res.(*RoutingFilterSelectChain), nil
}

func (c *Mikrotik) UpdateRoutingFilterSelectChain(selectChain *RoutingFilterSelectChain) (*RoutingFilterSelectChain, error) {
	res, err := c.Update(selectChain)
	if err != nil {
		return nil, err
	}
	return res.(*RoutingFilterSelectChain), nil
}

func (c *Mikrotik) DeleteRoutingFilterSelectChain(name string) error {
	selectChain := &RoutingFilterSelectChain{Name: name}
	return c.Delete(selectChain)
}
