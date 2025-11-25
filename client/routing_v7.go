package client

import (
	"fmt"

	"github.com/go-routeros/routeros/v3"
)

// RoutingTable represents RouterOS v7 routing table (VRF support)
// Reference: https://help.mikrotik.com/docs/display/ROS/Routing+Tables
type RoutingTable struct {
	Id       string `mikrotik:".id" codegen:"id,mikrotikID"`
	Name     string `mikrotik:"name" codegen:"name,required,terraformID"`
	Fib      string `mikrotik:"fib" codegen:"fib"`
	Disabled bool   `mikrotik:"disabled" codegen:"disabled"`
	Comment  string `mikrotik:"comment" codegen:"comment"`
}

var _ Resource = (*RoutingTable)(nil)

// IDField returns the field name for MikroTik ID
func (rt *RoutingTable) IDField() string {
	return ".id"
}

// ID returns the resource ID
func (rt *RoutingTable) ID() string {
	return rt.Id
}

// SetID sets the resource ID
func (rt *RoutingTable) SetID(id string) {
	rt.Id = id
}

// AfterAddHook is called after successful Add operation
func (rt *RoutingTable) AfterAddHook(r *routeros.Reply) {
	rt.Id = r.Done.Map["ret"]
}

// FindField returns the field name to use for finding resources
func (rt *RoutingTable) FindField() string {
	return "name"
}

// FindFieldValue returns the value for finding resources
func (rt *RoutingTable) FindFieldValue() string {
	return rt.Name
}

// DeleteField returns the field name to use for deletion
func (rt *RoutingTable) DeleteField() string {
	return "numbers"
}

// DeleteFieldValue returns the value for deletion
func (rt *RoutingTable) DeleteFieldValue() string {
	return rt.Name
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
	return map[Action]string{
		Add:    "/routing/table/add",
		Find:   "/routing/table/print",
		Update: "/routing/table/set",
		Delete: "/routing/table/remove",
	}[action]
}

// ActionToCommand returns the RouterOS CLI path for RoutingRule
func (rr *RoutingRule) ActionToCommand(action Action) string {
	return map[Action]string{
		Add:    "/routing/rule/add",
		Find:   "/routing/rule/print",
		Update: "/routing/rule/set",
		Delete: "/routing/rule/remove",
	}[action]
}

// ActionToCommand returns the RouterOS CLI path for VRF
func (vrf *VRF) ActionToCommand(action Action) string {
	return map[Action]string{
		Add:    "/ip/vrf/add",
		Find:   "/ip/vrf/print",
		Update: "/ip/vrf/set",
		Delete: "/ip/vrf/remove",
	}[action]
}

// FindRoutingTable finds a routing table by name (legacy method)
func (client Mikrotik) FindRoutingTable(name string) (*RoutingTable, error) {
	rt := &RoutingTable{Name: name}
	res, err := client.Find(rt)
	if err != nil {
		return nil, err
	}
	return res.(*RoutingTable), nil
}

// CreateRoutingTable creates a new routing table (VRF) (legacy method)
func (client Mikrotik) CreateRoutingTable(table *RoutingTable) (*RoutingTable, error) {
	res, err := client.Add(table)
	if err != nil {
		return nil, err
	}
	return res.(*RoutingTable), nil
}

// UpdateRoutingTable updates an existing routing table (legacy method)
func (client Mikrotik) UpdateRoutingTable(table *RoutingTable) (*RoutingTable, error) {
	res, err := client.Update(table)
	if err != nil {
		return nil, err
	}
	return res.(*RoutingTable), nil
}

// DeleteRoutingTable deletes a routing table (legacy method)
func (client Mikrotik) DeleteRoutingTable(name string) error {
	rt := &RoutingTable{Name: name}
	return client.Delete(rt)
}

// Typed wrappers for generic CRUD operations
func (c Mikrotik) AddRoutingTable(r *RoutingTable) (*RoutingTable, error) {
	res, err := c.Add(r)
	if err != nil {
		return nil, err
	}
	return res.(*RoutingTable), nil
}

func (c Mikrotik) UpdateRoutingTable2(r *RoutingTable) (*RoutingTable, error) {
	res, err := c.Update(r)
	if err != nil {
		return nil, err
	}
	return res.(*RoutingTable), nil
}

func (c Mikrotik) FindRoutingTable2(name string) (*RoutingTable, error) {
	rt := &RoutingTable{Name: name}
	res, err := c.Find(rt)
	if err != nil {
		return nil, err
	}
	return res.(*RoutingTable), nil
}

func (c Mikrotik) ListRoutingTables() ([]RoutingTable, error) {
	res, err := c.List(&RoutingTable{})
	if err != nil {
		return nil, err
	}
	returnSlice := make([]RoutingTable, len(res))
	for i, v := range res {
		returnSlice[i] = *(v.(*RoutingTable))
	}
	return returnSlice, nil
}

func (c Mikrotik) DeleteRoutingTable2(name string) error {
	rt := &RoutingTable{Name: name}
	return c.Delete(rt)
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
	err = Unmarshal(*reply, rule)
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
	err = Unmarshal(*reply, vrf)
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
	_, err = c.RunArgs(cmd)
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
