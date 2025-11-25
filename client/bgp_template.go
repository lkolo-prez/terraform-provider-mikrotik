package client

import (
	"github.com/go-routeros/routeros/v3"
)

// BgpTemplate represents a BGP template resource in RouterOS v7+
// Templates allow configuration reuse across multiple BGP connections.
// The template contains all BGP protocol-related configuration options.
//
// Reference: https://help.mikrotik.com/docs/display/ROS/BGP#BGP-TemplateMenu
type BgpTemplate struct {
	Id                     string `mikrotik:".id" codegen:"id,mikrotikID"`
	Name                   string `mikrotik:"name" codegen:"name,required,terraformID"`
	AS                     int    `mikrotik:"as" codegen:"as"`
	Disabled               bool   `mikrotik:"disabled" codegen:"disabled"`
	RouterID               string `mikrotik:"router-id" codegen:"router_id"`
	
	// Address families and capabilities
	AddressFamily          string `mikrotik:"address-families" codegen:"address_families,optional,computed"`
	Capabilities           string `mikrotik:"capabilities" codegen:"capabilities"`
	
	// AS manipulation
	AsOverride             bool   `mikrotik:"as-override" codegen:"as_override"`
	Cisco                  bool   `mikrotik:"cisco" codegen:"cisco"`
	RemovePrivateAS        bool   `mikrotik:"remove-private-as" codegen:"remove_private_as"`
	
	Comment                string `mikrotik:"comment" codegen:"comment"`
	
	// Timers
	ConnectRetryTime       string `mikrotik:"connect-retry-time" codegen:"connect_retry_time"`
	HoldTime               string `mikrotik:"hold-time" codegen:"hold_time,optional,computed"`
	KeepaliveTime          string `mikrotik:"keepalive-time" codegen:"keepalive_time,optional,computed"`
	
	// Input filtering and acceptance
	InputAffixFilters      string `mikrotik:"input.affixes" codegen:"input_affixes"`
	InputFilter            string `mikrotik:"input.filter" codegen:"input_filter"`
	InputLimit             int    `mikrotik:"input.limit" codegen:"input_limit"`
	InputAcceptCommunities string `mikrotik:"input.accept-communities" codegen:"input_accept_communities"`
	InputAcceptNLRI        string `mikrotik:"input.accept-nlri" codegen:"input_accept_nlri"`
	InputAcceptOriginated  bool   `mikrotik:"input.accept-originated" codegen:"input_accept_originated"`
	InputIgnoreAsPathLen   bool   `mikrotik:"input.ignore-as-path-len" codegen:"input_ignore_as_path_len"`
	InputLimitProcessRoutesIPv4 int `mikrotik:"input.limit-process-routes-ipv4" codegen:"input_limit_process_routes_ipv4"`
	InputLimitProcessRoutesIPv6 int `mikrotik:"input.limit-process-routes-ipv6" codegen:"input_limit_process_routes_ipv6"`
	
	// Multihop and BFD
	Multihop               bool   `mikrotik:"multihop" codegen:"multihop"`
	UseBFD                 bool   `mikrotik:"use-bfd" codegen:"use_bfd"`
	TTL                    string `mikrotik:"ttl" codegen:"ttl,optional,computed"`
	
	// Nexthop handling
	NexthopChoice          string `mikrotik:"nexthop-choice" codegen:"nexthop_choice,optional,computed"`
	
	// Output filtering and origination
	OutputAffixFilters     string `mikrotik:"output.affixes" codegen:"output_affixes"`
	OutputDefaultOriginate string `mikrotik:"output.default-originate" codegen:"output_default_originate,optional,computed"`
	OutputFilter           string `mikrotik:"output.filter" codegen:"output_filter"`
	OutputFilterChain      string `mikrotik:"output.filter-chain" codegen:"output_filter_chain"`
	OutputKeepaliveTime    string `mikrotik:"output.keepalive-time" codegen:"output_keepalive_time"`
	OutputNetwork          string `mikrotik:"output.network" codegen:"output_network"`
	OutputRedistribute     string `mikrotik:"output.redistribute" codegen:"output_redistribute"`
	
	// Passive mode and Route Reflection
	Passive                bool   `mikrotik:"passive" codegen:"passive"`
	RouteReflect           bool   `mikrotik:"route-reflect" codegen:"route_reflect"`
	
	// Graceful Restart
	GracefulRestart        string `mikrotik:"graceful-restart" codegen:"graceful_restart"`
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
