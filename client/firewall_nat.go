package client

import (
	"github.com/go-routeros/routeros/v3"
)

// FirewallNat defines /ip/firewall/nat resource
type FirewallNat struct {
	Id       string `mikrotik:".id" codegen:"id,mikrotikID"`
	Chain    string `mikrotik:"chain" codegen:"chain,required"`
	Action   string `mikrotik:"action" codegen:"action,required"`
	Disabled bool   `mikrotik:"disabled" codegen:"disabled"`
	Comment  string `mikrotik:"comment" codegen:"comment"`

	// Matching criteria
	SrcAddress     string `mikrotik:"src-address" codegen:"src_address"`
	DstAddress     string `mikrotik:"dst-address" codegen:"dst_address"`
	SrcAddressList string `mikrotik:"src-address-list" codegen:"src_address_list"`
	DstAddressList string `mikrotik:"dst-address-list" codegen:"dst_address_list"`
	Protocol       string `mikrotik:"protocol" codegen:"protocol"`
	SrcPort        string `mikrotik:"src-port" codegen:"src_port"`
	DstPort        string `mikrotik:"dst-port" codegen:"dst_port"`
	InInterface    string `mikrotik:"in-interface" codegen:"in_interface"`
	OutInterface   string `mikrotik:"out-interface" codegen:"out_interface"`
	InInterfaceList  string `mikrotik:"in-interface-list" codegen:"in_interface_list"`
	OutInterfaceList string `mikrotik:"out-interface-list" codegen:"out_interface_list"`

	// Connection tracking
	ConnectionState    string `mikrotik:"connection-state" codegen:"connection_state"`
	ConnectionNatState string `mikrotik:"connection-nat-state" codegen:"connection_nat_state"`
	ConnectionMark     string `mikrotik:"connection-mark" codegen:"connection_mark"`
	PacketMark         string `mikrotik:"packet-mark" codegen:"packet_mark"`
	RoutingMark        string `mikrotik:"routing-mark" codegen:"routing_mark"`

	// NAT action parameters
	ToAddresses string `mikrotik:"to-addresses" codegen:"to_addresses"`
	ToPorts     string `mikrotik:"to-ports" codegen:"to_ports"`

	// Logging
	Log       bool   `mikrotik:"log" codegen:"log"`
	LogPrefix string `mikrotik:"log-prefix" codegen:"log_prefix"`

	// Advanced matching
	IcmpOptions        string `mikrotik:"icmp-options" codegen:"icmp_options"`
	Limit              string `mikrotik:"limit" codegen:"limit"`
	Time               string `mikrotik:"time" codegen:"time"`
	Random             string `mikrotik:"random" codegen:"random"`
	HotspotAuth        string `mikrotik:"hotspot" codegen:"hotspot"`
	ContentType        string `mikrotik:"content-type" codegen:"content_type"`
	Layer7Protocol     string `mikrotik:"layer7-protocol" codegen:"layer7_protocol"`
	Psd                string `mikrotik:"psd" codegen:"psd"`
	TcpFlags           string `mikrotik:"tcp-flags" codegen:"tcp_flags"`
	TcpMss             string `mikrotik:"tcp-mss" codegen:"tcp_mss"`
	DstLimit           string `mikrotik:"dst-limit" codegen:"dst_limit"`
	AddressList        string `mikrotik:"address-list" codegen:"address_list"`
	AddressListTimeout string `mikrotik:"address-list-timeout" codegen:"address_list_timeout"`

	// Packet size
	PacketSize string `mikrotik:"packet-size" codegen:"packet_size"`

	// IPv6
	SrcAddressType string `mikrotik:"src-address-type" codegen:"src_address_type"`
	DstAddressType string `mikrotik:"dst-address-type" codegen:"dst_address_type"`

	// Read-only computed fields
	Bytes   int64 `mikrotik:"bytes" codegen:"bytes,computed"`
	Packets int64 `mikrotik:"packets" codegen:"packets,computed"`
	Dynamic bool  `mikrotik:"dynamic" codegen:"dynamic,computed"`
	Invalid bool  `mikrotik:"invalid" codegen:"invalid,computed"`
}

var _ Resource = (*FirewallNat)(nil)

func (f *FirewallNat) ActionToCommand(a Action) string {
	return map[Action]string{
		Add:    "/ip/firewall/nat/add",
		Find:   "/ip/firewall/nat/print",
		Update: "/ip/firewall/nat/set",
		Delete: "/ip/firewall/nat/remove",
	}[a]
}

func (f *FirewallNat) IDField() string {
	return ".id"
}

func (f *FirewallNat) ID() string {
	return f.Id
}

func (f *FirewallNat) SetID(id string) {
	f.Id = id
}

func (f *FirewallNat) AfterAddHook(r *routeros.Reply) {
	f.Id = r.Done.Map["ret"]
}

func (f *FirewallNat) FindField() string {
	return ".id"
}

func (f *FirewallNat) FindFieldValue() string {
	return f.Id
}

func (f *FirewallNat) DeleteField() string {
	return "numbers"
}

func (f *FirewallNat) DeleteFieldValue() string {
	return f.Id
}

// Typed wrappers
func (c Mikrotik) AddFirewallNat(r *FirewallNat) (*FirewallNat, error) {
	res, err := c.Add(r)
	if err != nil {
		return nil, err
	}

	return res.(*FirewallNat), nil
}

func (c Mikrotik) UpdateFirewallNat(r *FirewallNat) (*FirewallNat, error) {
	res, err := c.Update(r)
	if err != nil {
		return nil, err
	}

	return res.(*FirewallNat), nil
}

func (c Mikrotik) FindFirewallNat(id string) (*FirewallNat, error) {
	res, err := c.Find(&FirewallNat{Id: id})
	if err != nil {
		return nil, err
	}

	return res.(*FirewallNat), nil
}

func (c Mikrotik) DeleteFirewallNat(id string) error {
	return c.Delete(&FirewallNat{Id: id})
}
