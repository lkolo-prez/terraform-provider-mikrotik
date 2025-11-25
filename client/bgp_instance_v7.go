package client

import (
	"github.com/go-routeros/routeros/v3"
)

// BgpInstanceV7 represents a BGP routing instance in RouterOS v7.20+
// Starting from ROSv7.20, BGP routing instances are explicitly defined in instance menu
// instead of auto-detecting based on router-ids.
// 
// BGP routing instance is necessary for best path route selection and other
// instance-dependent features like VPN, EVPN, and so on.
//
// Reference: https://help.mikrotik.com/docs/display/ROS/BGP#BGP-InstanceMenu
type BgpInstanceV7 struct {
	Id                       string `mikrotik:".id" codegen:"id,mikrotikID"`
	Name                     string `mikrotik:"name" codegen:"name,required,terraformID"`
	AS                       int    `mikrotik:"as" codegen:"as,required"`
	RouterID                 string `mikrotik:"router-id" codegen:"router_id"`
	ClientToClientReflection bool   `mikrotik:"client-to-client-reflection" codegen:"client_to_client_reflection"`
	ClusterID                string `mikrotik:"cluster-id" codegen:"cluster_id"`
	Confederation            int    `mikrotik:"confederation" codegen:"confederation"`
	IgnoreAsPathLen          bool   `mikrotik:"ignore-as-path-len" codegen:"ignore_as_path_len"`
	OutFilter                string `mikrotik:"out-filter" codegen:"out_filter"`
	RoutingTable             string `mikrotik:"routing-table" codegen:"routing_table"`
	RedistributeConnected    bool   `mikrotik:"redistribute-connected" codegen:"redistribute_connected"`
	RedistributeOspf         bool   `mikrotik:"redistribute-ospf" codegen:"redistribute_ospf"`
	RedistributeOtherBgp     bool   `mikrotik:"redistribute-other-bgp" codegen:"redistribute_other_bgp"`
	RedistributeRip          bool   `mikrotik:"redistribute-rip" codegen:"redistribute_rip"`
	RedistributeStatic       bool   `mikrotik:"redistribute-static" codegen:"redistribute_static"`
	Disabled                 bool   `mikrotik:"disabled" codegen:"disabled"`
	Comment                  string `mikrotik:"comment" codegen:"comment"`
	
	// VPN/EVPN Support
	VRF                      string `mikrotik:"vrf" codegen:"vrf"`
}

var _ Resource = (*BgpInstanceV7)(nil)

func (b *BgpInstanceV7) ActionToCommand(a Action) string {
	return map[Action]string{
		Add:    "/routing/bgp/instance/add",
		Find:   "/routing/bgp/instance/print",
		Update: "/routing/bgp/instance/set",
		Delete: "/routing/bgp/instance/remove",
	}[a]
}

func (b *BgpInstanceV7) IDField() string {
	return ".id"
}

func (b *BgpInstanceV7) ID() string {
	return b.Id
}

func (b *BgpInstanceV7) SetID(id string) {
	b.Id = id
}

func (b *BgpInstanceV7) AfterAddHook(r *routeros.Reply) {
	b.Id = r.Done.Map["ret"]
}

func (b *BgpInstanceV7) FindField() string {
	return "name"
}

func (b *BgpInstanceV7) FindFieldValue() string {
	return b.Name
}

func (b *BgpInstanceV7) DeleteField() string {
	return "numbers"
}

func (b *BgpInstanceV7) DeleteFieldValue() string {
	return b.Name
}

// Typed wrappers
func (c Mikrotik) AddBgpInstanceV7(r *BgpInstanceV7) (*BgpInstanceV7, error) {
	res, err := c.Add(r)
	if err != nil {
		return nil, err
	}

	return res.(*BgpInstanceV7), nil
}

func (c Mikrotik) UpdateBgpInstanceV7(r *BgpInstanceV7) (*BgpInstanceV7, error) {
	res, err := c.Update(r)
	if err != nil {
		return nil, err
	}

	return res.(*BgpInstanceV7), nil
}

func (c Mikrotik) FindBgpInstanceV7(name string) (*BgpInstanceV7, error) {
	res, err := c.Find(&BgpInstanceV7{Name: name})
	if err != nil {
		return nil, err
	}

	return res.(*BgpInstanceV7), nil
}

func (c Mikrotik) DeleteBgpInstanceV7(name string) error {
	return c.Delete(&BgpInstanceV7{Name: name})
}
