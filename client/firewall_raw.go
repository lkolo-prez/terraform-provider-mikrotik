package client

import (
	"github.com/ddelnano/terraform-provider-mikrotik/client/types"
	"github.com/go-routeros/routeros/v3"
)

// FirewallRaw defines /ip/firewall/raw rule - new in RouterOS v7
// The RAW chain processes packets before connection tracking
type FirewallRaw struct {
	Id               string             `mikrotik:".id" codegen:"id,mikrotikID,terraformID"`
	Action           string             `mikrotik:"action" codegen:"action,required"`
	Chain            string             `mikrotik:"chain" codegen:"chain,required"`
	Comment          string             `mikrotik:"comment" codegen:"comment"`
	Disabled         bool               `mikrotik:"disabled" codegen:"disabled"`
	SrcAddress       string             `mikrotik:"src-address" codegen:"src_address"`
	DstAddress       string             `mikrotik:"dst-address" codegen:"dst_address"`
	SrcPort          string             `mikrotik:"src-port" codegen:"src_port"`
	DstPort          string             `mikrotik:"dst-port" codegen:"dst_port"`
	Protocol         string             `mikrotik:"protocol" codegen:"protocol"`
	InInterface      string             `mikrotik:"in-interface" codegen:"in_interface"`
	OutInterface     string             `mikrotik:"out-interface" codegen:"out_interface"`
	InInterfaceList  string             `mikrotik:"in-interface-list" codegen:"in_interface_list"`
	OutInterfaceList string             `mikrotik:"out-interface-list" codegen:"out_interface_list"`
	AddressList      string             `mikrotik:"address-list" codegen:"address_list"`
	AddressListTimeout string           `mikrotik:"address-list-timeout" codegen:"address_list_timeout"`
	SrcAddressList   string             `mikrotik:"src-address-list" codegen:"src_address_list"`
	DstAddressList   string             `mikrotik:"dst-address-list" codegen:"dst_address_list"`
	ConnectionState  types.MikrotikList `mikrotik:"connection-state" codegen:"connection_state"`
	ConnectionNatState types.MikrotikList `mikrotik:"connection-nat-state" codegen:"connection_nat_state"`
	Bytes            string             `mikrotik:"bytes" codegen:"bytes,readonly"`
	Packets          string             `mikrotik:"packets" codegen:"packets,readonly"`
}

var _ Resource = (*FirewallRaw)(nil)

func (b *FirewallRaw) ActionToCommand(a Action) string {
	return map[Action]string{
		Add:    "/ip/firewall/raw/add",
		Find:   "/ip/firewall/raw/print",
		Update: "/ip/firewall/raw/set",
		Delete: "/ip/firewall/raw/remove",
	}[a]
}

func (b *FirewallRaw) IDField() string {
	return ".id"
}

func (b *FirewallRaw) ID() string {
	return b.Id
}

func (b *FirewallRaw) SetID(id string) {
	b.Id = id
}

func (b *FirewallRaw) AfterAddHook(r *routeros.Reply) {
	b.Id = r.Done.Map["ret"]
}

func (c Mikrotik) AddFirewallRaw(r *FirewallRaw) (*FirewallRaw, error) {
	res, err := c.Add(r)
	if err != nil {
		return nil, err
	}

	return res.(*FirewallRaw), nil
}

func (c Mikrotik) UpdateFirewallRaw(r *FirewallRaw) (*FirewallRaw, error) {
	res, err := c.Update(r)
	if err != nil {
		return nil, err
	}

	return res.(*FirewallRaw), nil
}

func (c Mikrotik) FindFirewallRaw(id string) (*FirewallRaw, error) {
	res, err := c.Find(&FirewallRaw{Id: id})
	if err != nil {
		return nil, err
	}

	return res.(*FirewallRaw), nil
}

func (c Mikrotik) DeleteFirewallRaw(id string) error {
	return c.Delete(&FirewallRaw{Id: id})
}

// FastTrackConnection defines /ip/firewall/connection/tracking settings
// Fast Track is enhanced in RouterOS v7 for better performance
type FastTrackConnection struct {
	Enabled bool `mikrotik:"enabled" codegen:"enabled"`
}
