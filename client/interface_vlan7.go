package client

import (
	"github.com/go-routeros/routeros/v3"
)

// InterfaceVlan7 defines improved VLAN interface for RouterOS v7
// RouterOS 7 has improved VLAN filtering and bridge VLAN support
type InterfaceVlan7 struct {
	Id        string `mikrotik:".id" codegen:"id,mikrotikID"`
	Name      string `mikrotik:"name" codegen:"name,required,terraformID"`
	VlanId    int    `mikrotik:"vlan-id" codegen:"vlan_id,required"`
	Interface string `mikrotik:"interface" codegen:"interface,required"`
	Disabled  bool   `mikrotik:"disabled" codegen:"disabled"`
	Comment   string `mikrotik:"comment" codegen:"comment"`
	MTU       int    `mikrotik:"mtu" codegen:"mtu"`
	UseServiceTag bool `mikrotik:"use-service-tag" codegen:"use_service_tag"`
}

var _ Resource = (*InterfaceVlan7)(nil)

func (b *InterfaceVlan7) ActionToCommand(a Action) string {
	return map[Action]string{
		Add:    "/interface/vlan/add",
		Find:   "/interface/vlan/print",
		Update: "/interface/vlan/set",
		Delete: "/interface/vlan/remove",
	}[a]
}

func (b *InterfaceVlan7) IDField() string {
	return ".id"
}

func (b *InterfaceVlan7) ID() string {
	return b.Id
}

func (b *InterfaceVlan7) SetID(id string) {
	b.Id = id
}

func (b *InterfaceVlan7) AfterAddHook(r *routeros.Reply) {
	b.Id = r.Done.Map["ret"]
}

func (b *InterfaceVlan7) FindField() string {
	return "name"
}

func (b *InterfaceVlan7) FindFieldValue() string {
	return b.Name
}

func (b *InterfaceVlan7) DeleteField() string {
	return "numbers"
}

func (b *InterfaceVlan7) DeleteFieldValue() string {
	return b.Name
}

// BridgeVlanFiltering defines bridge VLAN filtering for RouterOS v7
// This enables hardware-accelerated VLAN filtering on bridges
type BridgeVlanFiltering struct {
	Id               string `mikrotik:".id" codegen:"id,mikrotikID,terraformID"`
	Bridge           string `mikrotik:"bridge" codegen:"bridge,required"`
	VlanFiltering    bool   `mikrotik:"vlan-filtering" codegen:"vlan_filtering"`
	PVIDMode         string `mikrotik:"pvid-mode" codegen:"pvid_mode"`
	FrameTypes       string `mikrotik:"frame-types" codegen:"frame_types"`
	IngressFiltering bool   `mikrotik:"ingress-filtering" codegen:"ingress_filtering"`
	EtherType        string `mikrotik:"ether-type" codegen:"ether_type"`
}

var _ Resource = (*BridgeVlanFiltering)(nil)

func (b *BridgeVlanFiltering) ActionToCommand(a Action) string {
	return map[Action]string{
		Find:   "/interface/bridge/print",
		Update: "/interface/bridge/set",
	}[a]
}

func (b *BridgeVlanFiltering) IDField() string {
	return ".id"
}

func (b *BridgeVlanFiltering) ID() string {
	return b.Id
}

func (b *BridgeVlanFiltering) SetID(id string) {
	b.Id = id
}

func (b *BridgeVlanFiltering) AfterAddHook(r *routeros.Reply) {
	b.Id = r.Done.Map["ret"]
}

// Typed wrappers for VLAN
func (c Mikrotik) AddInterfaceVlan7(r *InterfaceVlan7) (*InterfaceVlan7, error) {
	res, err := c.Add(r)
	if err != nil {
		return nil, err
	}

	return res.(*InterfaceVlan7), nil
}

func (c Mikrotik) UpdateInterfaceVlan7(r *InterfaceVlan7) (*InterfaceVlan7, error) {
	res, err := c.Update(r)
	if err != nil {
		return nil, err
	}

	return res.(*InterfaceVlan7), nil
}

func (c Mikrotik) FindInterfaceVlan7(name string) (*InterfaceVlan7, error) {
	res, err := c.Find(&InterfaceVlan7{Name: name})
	if err != nil {
		return nil, err
	}

	return res.(*InterfaceVlan7), nil
}

func (c Mikrotik) DeleteInterfaceVlan7(name string) error {
	return c.Delete(&InterfaceVlan7{Name: name})
}
