package client

import (
	"fmt"
	"github.com/go-routeros/routeros/v3"
)

// RoutingTable represents RouterOS v7 routing table (VRF support)
// Reference: https://help.mikrotik.com/docs/display/ROS/Routing+Tables
type RoutingTable struct {
	Id       string `mikrotik:".id" codegen:"id,mikrotikID"`
	Name     string `mikrotik:"name" codegen:"name,required"`
	Fib      string `mikrotik:"fib" codegen:"fib"`      // Push routes to FIB: yes/no
	Disabled string `mikrotik:"disabled" codegen:"disabled"`
	Comment  string `mikrotik:"comment" codegen:"comment"`
}

// RoutingRule represents RouterOS v7 routing rule (policy-based routing)
// Reference: https://help.mikrotik.com/docs/display/ROS/Policy+Routing
type RoutingRule struct {
	Id              string `mikrotik:".id" codegen:"id,mikrotikID"`
	DstAddress      string `mikrotik:"dst-address" codegen:"dst_address"`
	SrcAddress      string `mikrotik:"src-address" codegen:"src_address"`
	Interface       string `mikrotik:"interface" codegen:"interface"`
	RoutingMark     string `mikrotik:"routing-mark" codegen:"routing_mark"`
	Table           string `mikrotik:"table" codegen:"table,required"`
	Action          string `mikrotik:"action" codegen:"action,required"` // lookup, lookup-only-in-table, unreachable, blackhole
	Disabled        string `mikrotik:"disabled" codegen:"disabled"`
	Comment         string `mikrotik:"comment" codegen:"comment"`
	MinPrefix       string `mikrotik:"min-prefix" codegen:"min_prefix"`
	MaxPrefix       string `mikrotik:"max-prefix" codegen:"max_prefix"`
	IngressPriority string `mikrotik:"ingress-priority" codegen:"ingress_priority"`
}

// VRF represents RouterOS v7 VRF (Virtual Routing and Forwarding)
// Reference: https://help.mikrotik.com/docs/display/ROS/VRF
type VRF struct {
	Id        string `mikrotik:".id" codegen:"id,mikrotikID"`
	Name      string `mikrotik:"name" codegen:"name,required"`
	Interfaces string `mikrotik:"interfaces" codegen:"interfaces"` // Comma-separated list
	Disabled  string `mikrotik:"disabled" codegen:"disabled"`
	Comment   string `mikrotik:"comment" codegen:"comment"`
}

// ActionToCommand returns the RouterOS CLI path for RoutingTable
func (rt *RoutingTable) ActionToCommand(action Action) string {
	return "/routing/table" + action.ToCommand()
}

// ActionToCommand returns the RouterOS CLI path for RoutingRule
func (rr *RoutingRule) ActionToCommand(action Action) string {
	return "/routing/rule" + action.ToCommand()
}

// ActionToCommand returns the RouterOS CLI path for VRF
func (vrf *VRF) ActionToCommand(action Action) string {
	return "/ip/vrf" + action.ToCommand()
}

// FindRoutingTable finds a routing table by name
func (client Mikrotik) FindRoutingTable(name string) (*RoutingTable, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := []string{"/routing/table/print", "?name=" + name}
	reply, err := c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}

	if len(reply.Re) == 0 {
		return nil, NewNotFound(fmt.Sprintf("routing table '%s' not found", name))
	}

	table := &RoutingTable{}
	err = Unmarshal(*reply.Re[0], table)
	if err != nil {
		return nil, err
	}

	return table, nil
}

// CreateRoutingTable creates a new routing table (VRF)
func (client Mikrotik) CreateRoutingTable(table *RoutingTable) (*RoutingTable, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := Marshal("/routing/table/add", table)
	reply, err := c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}

	return client.FindRoutingTable(table.Name)
}

// UpdateRoutingTable updates an existing routing table
func (client Mikrotik) UpdateRoutingTable(table *RoutingTable) (*RoutingTable, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := Marshal("/routing/table/set", table)
	_, err = c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}

	return client.FindRoutingTable(table.Name)
}

// DeleteRoutingTable deletes a routing table
func (client Mikrotik) DeleteRoutingTable(name string) error {
	c, err := client.getMikrotikClient()
	if err != nil {
		return err
	}

	table, err := client.FindRoutingTable(name)
	if err != nil {
		return err
	}

	cmd := []string{"/routing/table/remove", "=.id=" + table.Id}
	_, err = c.RunArgs(cmd)
	return err
}

// FindRoutingRule finds a routing rule by ID
func (client Mikrotik) FindRoutingRule(id string) (*RoutingRule, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := []string{"/routing/rule/print", "?.id=" + id}
	reply, err := c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}

	if len(reply.Re) == 0 {
		return nil, NewNotFound(fmt.Sprintf("routing rule '%s' not found", id))
	}

	rule := &RoutingRule{}
	err = Unmarshal(*reply.Re[0], rule)
	if err != nil {
		return nil, err
	}

	return rule, nil
}

// CreateRoutingRule creates a new routing rule
func (client Mikrotik) CreateRoutingRule(rule *RoutingRule) (*RoutingRule, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := Marshal("/routing/rule/add", rule)
	reply, err := c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}

	id := reply.Done.Map["ret"]
	return client.FindRoutingRule(id)
}

// UpdateRoutingRule updates an existing routing rule
func (client Mikrotik) UpdateRoutingRule(rule *RoutingRule) (*RoutingRule, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := Marshal("/routing/rule/set", rule)
	_, err = c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}

	return client.FindRoutingRule(rule.Id)
}

// DeleteRoutingRule deletes a routing rule
func (client Mikrotik) DeleteRoutingRule(id string) error {
	c, err := client.getMikrotikClient()
	if err != nil {
		return err
	}

	cmd := []string{"/routing/rule/remove", "=.id=" + id}
	_, err = c.RunArgs(cmd)
	return err
}

// FindVRF finds a VRF by name
func (client Mikrotik) FindVRF(name string) (*VRF, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := []string{"/ip/vrf/print", "?name=" + name}
	reply, err := c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}

	if len(reply.Re) == 0 {
		return nil, NewNotFound(fmt.Sprintf("VRF '%s' not found", name))
	}

	vrf := &VRF{}
	err = Unmarshal(*reply.Re[0], vrf)
	if err != nil {
		return nil, err
	}

	return vrf, nil
}

// CreateVRF creates a new VRF
func (client Mikrotik) CreateVRF(vrf *VRF) (*VRF, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := Marshal("/ip/vrf/add", vrf)
	reply, err := c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}

	return client.FindVRF(vrf.Name)
}

// UpdateVRF updates an existing VRF
func (client Mikrotik) UpdateVRF(vrf *VRF) (*VRF, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := Marshal("/ip/vrf/set", vrf)
	_, err = c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}

	return client.FindVRF(vrf.Name)
}

// DeleteVRF deletes a VRF
func (client Mikrotik) DeleteVRF(name string) error {
	c, err := client.getMikrotikClient()
	if err != nil {
		return err
	}

	vrf, err := client.FindVRF(name)
	if err != nil {
		return err
	}

	cmd := []string{"/ip/vrf/remove", "=.id=" + vrf.Id}
	_, err = c.RunArgs(cmd)
	return err
}
