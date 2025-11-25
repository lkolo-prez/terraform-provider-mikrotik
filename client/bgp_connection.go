package client

import (
	"github.com/go-routeros/routeros/v3"
)

// BgpConnection represents a BGP connection resource in RouterOS v7+
// This replaces the legacy BGP instance and peer configuration
type BgpConnection struct {
	Id                string `mikrotik:".id" codegen:"id,mikrotikID"`
	Name              string `mikrotik:"name" codegen:"name,required,terraformID"`
	AS                int    `mikrotik:"as" codegen:"as,required"`
	Disabled          bool   `mikrotik:"disabled" codegen:"disabled"`
	LocalRole         string `mikrotik:"local.role" codegen:"local_role,optional,computed"`
	LocalAddress      string `mikrotik:"local.address" codegen:"local_address"`
	RemoteAddress     string `mikrotik:"remote.address" codegen:"remote_address,required"`
	RemoteAS          int    `mikrotik:"remote.as" codegen:"remote_as"`
	RemotePort        int    `mikrotik:"remote.port" codegen:"remote_port"`
	RouterID          string `mikrotik:"router-id" codegen:"router_id"`
	Nexthop           string `mikrotik:"nexthop-choice" codegen:"nexthop_choice,optional,computed"`
	HoldTime          string `mikrotik:"hold-time" codegen:"hold_time,optional,computed"`
	KeepaliveTime     string `mikrotik:"keepalive-time" codegen:"keepalive_time,optional,computed"`
	ConnectRetryTime  string `mikrotik:"connect-retry-time" codegen:"connect_retry_time"`
	TTL               string `mikrotik:"ttl" codegen:"ttl,optional,computed"`
	Multihop          bool   `mikrotik:"multihop" codegen:"multihop"`
	UseBFD            bool   `mikrotik:"use-bfd" codegen:"use_bfd"`
	AddressFamily     string `mikrotik:"address-families" codegen:"address_families,optional,computed"`
	Comment           string `mikrotik:"comment" codegen:"comment"`
	Templates         string `mikrotik:"templates" codegen:"templates"`
	InputFilter       string `mikrotik:"input.filter" codegen:"input_filter"`
	OutputFilter      string `mikrotik:"output.filter" codegen:"output_filter"`
	OutputDefaultOriginate string `mikrotik:"output.default-originate" codegen:"output_default_originate,optional,computed"`
	OutputNetwork     string `mikrotik:"output.network" codegen:"output_network"`
	TCPMd5Key         string `mikrotik:"tcp-md5-key" codegen:"tcp_md5_key"`
	UseMPLS           bool   `mikrotik:"use-mpls" codegen:"use_mpls"`
	VPNV4             bool   `mikrotik:"vpnv4" codegen:"vpnv4"`
	VPNV6             bool   `mikrotik:"vpnv6" codegen:"vpnv6"`
}

var _ Resource = (*BgpConnection)(nil)

func (b *BgpConnection) ActionToCommand(a Action) string {
	return map[Action]string{
		Add:    "/routing/bgp/connection/add",
		Find:   "/routing/bgp/connection/print",
		Update: "/routing/bgp/connection/set",
		Delete: "/routing/bgp/connection/remove",
	}[a]
}

func (b *BgpConnection) IDField() string {
	return ".id"
}

func (b *BgpConnection) ID() string {
	return b.Id
}

func (b *BgpConnection) SetID(id string) {
	b.Id = id
}

func (b *BgpConnection) AfterAddHook(r *routeros.Reply) {
	b.Id = r.Done.Map["ret"]
}

func (b *BgpConnection) FindField() string {
	return "name"
}

func (b *BgpConnection) FindFieldValue() string {
	return b.Name
}

func (b *BgpConnection) DeleteField() string {
	return "numbers"
}

func (b *BgpConnection) DeleteFieldValue() string {
	return b.Name
}

// Typed wrappers
func (c Mikrotik) AddBgpConnection(r *BgpConnection) (*BgpConnection, error) {
	res, err := c.Add(r)
	if err != nil {
		return nil, err
	}

	return res.(*BgpConnection), nil
}

func (c Mikrotik) UpdateBgpConnection(r *BgpConnection) (*BgpConnection, error) {
	res, err := c.Update(r)
	if err != nil {
		return nil, err
	}

	return res.(*BgpConnection), nil
}

func (c Mikrotik) FindBgpConnection(name string) (*BgpConnection, error) {
	res, err := c.Find(&BgpConnection{Name: name})
	if err != nil {
		return nil, err
	}

	return res.(*BgpConnection), nil
}

func (c Mikrotik) DeleteBgpConnection(name string) error {
	return c.Delete(&BgpConnection{Name: name})
}
