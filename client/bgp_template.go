package client

import (
	"github.com/go-routeros/routeros/v3"
)

// BgpTemplate represents a BGP template resource in RouterOS v7+
// Templates allow configuration reuse across multiple BGP connections
type BgpTemplate struct {
	Id                     string `mikrotik:".id" codegen:"id,mikrotikID"`
	Name                   string `mikrotik:"name" codegen:"name,required,terraformID"`
	AS                     int    `mikrotik:"as" codegen:"as"`
	Disabled               bool   `mikrotik:"disabled" codegen:"disabled"`
	RouterID               string `mikrotik:"router-id" codegen:"router_id"`
	AddressFamily          string `mikrotik:"address-families" codegen:"address_families,optional,computed"`
	AsOverride             bool   `mikrotik:"as-override" codegen:"as_override"`
	Cisco                  bool   `mikrotik:"cisco" codegen:"cisco"`
	Comment                string `mikrotik:"comment" codegen:"comment"`
	ConnectRetryTime       string `mikrotik:"connect-retry-time" codegen:"connect_retry_time"`
	HoldTime               string `mikrotik:"hold-time" codegen:"hold_time,optional,computed"`
	InputAffixFilters      string `mikrotik:"input.affixes" codegen:"input_affixes"`
	InputFilter            string `mikrotik:"input.filter" codegen:"input_filter"`
	InputLimit             int    `mikrotik:"input.limit" codegen:"input_limit"`
	KeepaliveTime          string `mikrotik:"keepalive-time" codegen:"keepalive_time,optional,computed"`
	Multihop               bool   `mikrotik:"multihop" codegen:"multihop"`
	NexthopChoice          string `mikrotik:"nexthop-choice" codegen:"nexthop_choice,optional,computed"`
	OutputAffixFilters     string `mikrotik:"output.affixes" codegen:"output_affixes"`
	OutputDefaultOriginate string `mikrotik:"output.default-originate" codegen:"output_default_originate,optional,computed"`
	OutputFilter           string `mikrotik:"output.filter" codegen:"output_filter"`
	OutputFilterChain      string `mikrotik:"output.filter-chain" codegen:"output_filter_chain"`
	OutputKeepaliveTime    string `mikrotik:"output.keepalive-time" codegen:"output_keepalive_time"`
	OutputNetwork          string `mikrotik:"output.network" codegen:"output_network"`
	Passive                bool   `mikrotik:"passive" codegen:"passive"`
	RemovePrivateAS        bool   `mikrotik:"remove-private-as" codegen:"remove_private_as"`
	RouteReflect           bool   `mikrotik:"route-reflect" codegen:"route_reflect"`
	TTL                    string `mikrotik:"ttl" codegen:"ttl,optional,computed"`
	UseBFD                 bool   `mikrotik:"use-bfd" codegen:"use_bfd"`
}

var _ Resource = (*BgpTemplate)(nil)

func (b *BgpTemplate) ActionToCommand(a Action) string {
	return map[Action]string{
		Add:    "/routing/bgp/template/add",
		Find:   "/routing/bgp/template/print",
		Update: "/routing/bgp/template/set",
		Delete: "/routing/bgp/template/remove",
	}[a]
}

func (b *BgpTemplate) IDField() string {
	return ".id"
}

func (b *BgpTemplate) ID() string {
	return b.Id
}

func (b *BgpTemplate) SetID(id string) {
	b.Id = id
}

func (b *BgpTemplate) AfterAddHook(r *routeros.Reply) {
	b.Id = r.Done.Map["ret"]
}

func (b *BgpTemplate) FindField() string {
	return "name"
}

func (b *BgpTemplate) FindFieldValue() string {
	return b.Name
}

func (b *BgpTemplate) DeleteField() string {
	return "numbers"
}

func (b *BgpTemplate) DeleteFieldValue() string {
	return b.Name
}

// Typed wrappers
func (c Mikrotik) AddBgpTemplate(r *BgpTemplate) (*BgpTemplate, error) {
	res, err := c.Add(r)
	if err != nil {
		return nil, err
	}

	return res.(*BgpTemplate), nil
}

func (c Mikrotik) UpdateBgpTemplate(r *BgpTemplate) (*BgpTemplate, error) {
	res, err := c.Update(r)
	if err != nil {
		return nil, err
	}

	return res.(*BgpTemplate), nil
}

func (c Mikrotik) FindBgpTemplate(name string) (*BgpTemplate, error) {
	res, err := c.Find(&BgpTemplate{Name: name})
	if err != nil {
		return nil, err
	}

	return res.(*BgpTemplate), nil
}

func (c Mikrotik) DeleteBgpTemplate(name string) error {
	return c.Delete(&BgpTemplate{Name: name})
}
