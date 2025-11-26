package client

import "github.com/go-routeros/routeros/v3"

// Snmp represents RouterOS /snmp configuration
// Path: /snmp
//
// SNMP (Simple Network Management Protocol) configuration for monitoring.
// Supports SNMPv2c communities and SNMPv3 users with authentication and encryption.
type Snmp struct {
	Enabled       bool   `mikrotik:"enabled" codegen:"enabled"`
	Contact       string `mikrotik:"contact" codegen:"contact"`
	Location      string `mikrotik:"location" codegen:"location"`
	EngineId      string `mikrotik:"engine-id" codegen:"engine_id,read"`
	TrapVersion   string `mikrotik:"trap-version" codegen:"trap_version"`
	TrapCommunity string `mikrotik:"trap-community" codegen:"trap_community"`
	TrapTarget    string `mikrotik:"trap-target" codegen:"trap_target"`
	TrapGenerators string `mikrotik:"trap-generators" codegen:"trap_generators"`
}

// Resource interface implementation for Snmp
// Note: SNMP uses /snmp/set (not add/remove) as it's a singleton configuration

func (s *Snmp) IDField() string {
	return ".id"
}

func (s *Snmp) ID() string {
	return "*0" // SNMP is singleton, always *0
}

func (s *Snmp) SetID(id string) {
	// SNMP is singleton, ID doesn't change
}

func (s *Snmp) AfterAddHook(reply *routeros.Reply) {
	// SNMP doesn't have add operation
}

func (s *Snmp) FindField() string {
	return ".id"
}

func (s *Snmp) FindFieldValue() string {
	return "*0"
}

func (s *Snmp) DeleteField() string {
	return ".id"
}

func (s *Snmp) DeleteFieldValue() string {
	return "*0"
}

func (s *Snmp) ActionToCommand(action Action) string {
	return map[Action]string{
		Add:    "/snmp/set", // SNMP uses set, not add
		Find:   "/snmp/print",
		Update: "/snmp/set",
		Delete: "/snmp/set", // Can't delete SNMP, only disable
	}[action]
}

// CRUD wrappers for Snmp

func (c *Mikrotik) GetSnmp() (*Snmp, error) {
	res, err := c.Find(&Snmp{})
	if err != nil {
		return nil, err
	}

	return res.(*Snmp), nil
}

func (c *Mikrotik) UpdateSnmp(d *Snmp) (*Snmp, error) {
	res, err := c.Update(d)
	if err != nil {
		return nil, err
	}

	return res.(*Snmp), nil
}

// SnmpCommunity represents RouterOS /snmp/community configuration
// Path: /snmp/community
//
// SNMPv2c communities for read-only or read-write access.
type SnmpCommunity struct {
	Id        string `mikrotik:".id" codegen:"id,mikrotikID,read"`
	Name      string `mikrotik:"name" codegen:"name,required"`
	Security  string `mikrotik:"security" codegen:"security"` // none, authorized, private
	ReadAccess bool  `mikrotik:"read-access" codegen:"read_access"`
	WriteAccess bool `mikrotik:"write-access" codegen:"write_access"`
	Address   string `mikrotik:"address" codegen:"address"`
	Disabled  bool   `mikrotik:"disabled" codegen:"disabled"`
}

// Resource interface implementation for SnmpCommunity

func (c *SnmpCommunity) IDField() string {
	return ".id"
}

func (c *SnmpCommunity) ID() string {
	return c.Id
}

func (c *SnmpCommunity) SetID(id string) {
	c.Id = id
}

func (c *SnmpCommunity) AfterAddHook(reply *routeros.Reply) {
	c.Id = reply.Done.Map["ret"]
}

func (c *SnmpCommunity) FindField() string {
	return "name"
}

func (c *SnmpCommunity) FindFieldValue() string {
	return c.Name
}

func (c *SnmpCommunity) DeleteField() string {
	return "numbers"
}

func (c *SnmpCommunity) DeleteFieldValue() string {
	return c.Id
}

func (c *SnmpCommunity) ActionToCommand(action Action) string {
	return map[Action]string{
		Add:    "/snmp/community/add",
		Find:   "/snmp/community/print",
		Update: "/snmp/community/set",
		Delete: "/snmp/community/remove",
	}[action]
}

// CRUD wrappers for SnmpCommunity

func (c *Mikrotik) AddSnmpCommunity(d *SnmpCommunity) (*SnmpCommunity, error) {
	res, err := c.Add(d)
	if err != nil {
		return nil, err
	}

	return res.(*SnmpCommunity), nil
}

func (c *Mikrotik) UpdateSnmpCommunity(d *SnmpCommunity) (*SnmpCommunity, error) {
	res, err := c.Update(d)
	if err != nil {
		return nil, err
	}

	return res.(*SnmpCommunity), nil
}

func (c *Mikrotik) FindSnmpCommunity(name string) (*SnmpCommunity, error) {
	res, err := c.Find(&SnmpCommunity{Name: name})
	if err != nil {
		return nil, err
	}

	return res.(*SnmpCommunity), nil
}

func (c *Mikrotik) DeleteSnmpCommunity(id string) error {
	return c.Delete(&SnmpCommunity{Id: id})
}
