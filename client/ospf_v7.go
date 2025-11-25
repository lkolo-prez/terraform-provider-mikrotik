package client

import (
	"github.com/ddelnano/terraform-provider-mikrotik/client/types"
	"github.com/go-routeros/routeros/v3"
)

// OspfInstanceV7 represents /routing/ospf/instance
type OspfInstanceV7 struct {
	Id                   string
	Name                 string `mikrotik:"name"`
	Version              string `mikrotik:"version"`             // "2" or "3"
	RouterId             string `mikrotik:"router-id"`
	DomainId             string `mikrotik:"domain-id"`
	Disabled             string `mikrotik:"disabled"`
	Comment              string `mikrotik:"comment"`
	Vrf                  string `mikrotik:"vrf"`
	RoutingTable         string `mikrotik:"routing-table"`
	// Redistribution flags
	RedistributeConnected string `mikrotik:"redistribute-connected"` // yes/no
	RedistributeStatic    string `mikrotik:"redistribute-static"`    // yes/no
	RedistributeBgp       string `mikrotik:"redistribute-bgp"`       // yes/no
	RedistributeRip       string `mikrotik:"redistribute-rip"`       // yes/no
	RedistributeOspf      string `mikrotik:"redistribute-ospf"`      // yes/no
	// Default route origination
	OriginateDefault string `mikrotik:"originate-default"` // never, always, if-installed
	// Filters
	InFilterChain  string `mikrotik:"in-filter-chain"`
	OutFilterChain string `mikrotik:"out-filter-chain"`
	// Computed fields
	RoutingMarks string `mikrotik:"routing-marks"`
	Dynamic      string `mikrotik:"dynamic"`
	Invalid      string `mikrotik:"invalid"`
}

// Implement Resource interface for OspfInstanceV7
func (b *OspfInstanceV7) ActionToCommand(a Action) string {
	return "/routing/ospf/instance/" + string(a)
}

func (b *OspfInstanceV7) IDField() string {
	return ".id"
}

func (b *OspfInstanceV7) ID() string {
	return b.Id
}

func (b *OspfInstanceV7) SetID(id string) {
	b.Id = id
}

func (b *OspfInstanceV7) AfterAddHook(r *routeros.Reply) {
	b.Id = r.Done.Map["ret"]
}

func (b *OspfInstanceV7) FindField() string {
	return "name"
}

func (b *OspfInstanceV7) FindFieldValue() string {
	return b.Name
}

func (b *OspfInstanceV7) DeleteField() string {
	return "numbers"
}

func (b *OspfInstanceV7) DeleteFieldValue() string {
	return b.Id
}

// CRUD wrappers for OspfInstanceV7
func (c Mikrotik) AddOspfInstanceV7(d *OspfInstanceV7) (*OspfInstanceV7, error) {
	res, err := c.Add(d)
	if err != nil {
		return nil, err
	}

	return res.(*OspfInstanceV7), nil
}

func (c Mikrotik) FindOspfInstanceV7(name string) (*OspfInstanceV7, error) {
	res, err := c.Find(&OspfInstanceV7{Name: name})
	if err != nil {
		return nil, err
	}

	return res.(*OspfInstanceV7), nil
}

func (c Mikrotik) UpdateOspfInstanceV7(d *OspfInstanceV7) (*OspfInstanceV7, error) {
	res, err := c.Update(d)
	if err != nil {
		return nil, err
	}

	return res.(*OspfInstanceV7), nil
}

func (c Mikrotik) DeleteOspfInstanceV7(name string) error {
	d, err := c.FindOspfInstanceV7(name)
	if err != nil {
		return err
	}

	err = c.Delete(d)
	return err
}

// OspfAreaV7 represents /routing/ospf/area
type OspfAreaV7 struct {
	Id          string
	Name        string `mikrotik:"name"`
	AreaId      string `mikrotik:"area-id"` // 0.0.0.0 format
	Instance    string `mikrotik:"instance"`
	Type        string `mikrotik:"type"` // default, stub, nssa
	Disabled    string `mikrotik:"disabled"`
	Comment     string `mikrotik:"comment"`
	// Stub area options
	DefaultCost string `mikrotik:"default-cost"`
	NoSummaries string `mikrotik:"no-summaries"` // yes/no, for totally stubby
	// NSSA options
	NssaTranslator  string `mikrotik:"nssa-translator"`  // yes/no/candidate
	NssaPropagation string `mikrotik:"nssa-propagation"` // yes/no
	// Computed fields
	Dynamic string `mikrotik:"dynamic"`
	Invalid string `mikrotik:"invalid"`
}

// Implement Resource interface for OspfAreaV7
func (b *OspfAreaV7) ActionToCommand(a Action) string {
	return "/routing/ospf/area/" + string(a)
}

func (b *OspfAreaV7) IDField() string {
	return ".id"
}

func (b *OspfAreaV7) ID() string {
	return b.Id
}

func (b *OspfAreaV7) SetID(id string) {
	b.Id = id
}

func (b *OspfAreaV7) AfterAddHook(r *routeros.Reply) {
	b.Id = r.Done.Map["ret"]
}

func (b *OspfAreaV7) FindField() string {
	return "name"
}

func (b *OspfAreaV7) FindFieldValue() string {
	return b.Name
}

func (b *OspfAreaV7) DeleteField() string {
	return "numbers"
}

func (b *OspfAreaV7) DeleteFieldValue() string {
	return b.Id
}

// CRUD wrappers for OspfAreaV7
func (c Mikrotik) AddOspfAreaV7(d *OspfAreaV7) (*OspfAreaV7, error) {
	res, err := c.Add(d)
	if err != nil {
		return nil, err
	}

	return res.(*OspfAreaV7), nil
}

func (c Mikrotik) FindOspfAreaV7(name string) (*OspfAreaV7, error) {
	res, err := c.Find(&OspfAreaV7{Name: name})
	if err != nil {
		return nil, err
	}

	return res.(*OspfAreaV7), nil
}

func (c Mikrotik) UpdateOspfAreaV7(d *OspfAreaV7) (*OspfAreaV7, error) {
	res, err := c.Update(d)
	if err != nil {
		return nil, err
	}

	return res.(*OspfAreaV7), nil
}

func (c Mikrotik) DeleteOspfAreaV7(name string) error {
	d, err := c.FindOspfAreaV7(name)
	if err != nil {
		return err
	}

	err = c.Delete(d)
	return err
}

// OspfInterfaceTemplateV7 represents /routing/ospf/interface-template
type OspfInterfaceTemplateV7 struct {
	Id         string
	Area       string              `mikrotik:"area"`
	Networks   types.MikrotikList  `mikrotik:"networks"` // list of CIDR blocks
	Interfaces types.MikrotikList  `mikrotik:"interfaces"`
	Type       string              `mikrotik:"type"` // broadcast, ptp, ptmp, nbma, ptmp-broadcast, virtual-link
	Disabled   string              `mikrotik:"disabled"`
	Comment    string              `mikrotik:"comment"`
	// Cost and priority
	Cost     string `mikrotik:"cost"`
	Priority string `mikrotik:"priority"`
	Passive  string `mikrotik:"passive"` // yes/no
	// Authentication
	Auth     string `mikrotik:"auth"`     // none, simple, md5, sha1, sha256, sha384, sha512
	AuthKey  string `mikrotik:"auth-key"` // sensitive
	AuthId   string `mikrotik:"auth-id"`  // for MD5 and SHA
	// Timers
	HelloInterval      string `mikrotik:"hello-interval"`
	DeadInterval       string `mikrotik:"dead-interval"`
	RetransmitInterval string `mikrotik:"retransmit-interval"`
	TransmitDelay      string `mikrotik:"transmit-delay"`
	WaitTime           string `mikrotik:"wait-time"`
	// Virtual link options
	VlinkTransitArea string `mikrotik:"vlink-transit-area"`
	VlinkNeighborId  string `mikrotik:"vlink-neighbor-id"` // router-id format
	// Computed fields
	Dynamic string `mikrotik:"dynamic"`
	Invalid string `mikrotik:"invalid"`
}

// Implement Resource interface for OspfInterfaceTemplateV7
func (b *OspfInterfaceTemplateV7) ActionToCommand(a Action) string {
	return "/routing/ospf/interface-template/" + string(a)
}

func (b *OspfInterfaceTemplateV7) IDField() string {
	return ".id"
}

func (b *OspfInterfaceTemplateV7) ID() string {
	return b.Id
}

func (b *OspfInterfaceTemplateV7) SetID(id string) {
	b.Id = id
}

func (b *OspfInterfaceTemplateV7) AfterAddHook(r *routeros.Reply) {
	b.Id = r.Done.Map["ret"]
}

func (b *OspfInterfaceTemplateV7) FindField() string {
	return ".id"
}

func (b *OspfInterfaceTemplateV7) FindFieldValue() string {
	return b.Id
}

func (b *OspfInterfaceTemplateV7) DeleteField() string {
	return "numbers"
}

func (b *OspfInterfaceTemplateV7) DeleteFieldValue() string {
	return b.Id
}

// CRUD wrappers for OspfInterfaceTemplateV7
func (c Mikrotik) AddOspfInterfaceTemplateV7(d *OspfInterfaceTemplateV7) (*OspfInterfaceTemplateV7, error) {
	res, err := c.Add(d)
	if err != nil {
		return nil, err
	}

	return res.(*OspfInterfaceTemplateV7), nil
}

func (c Mikrotik) FindOspfInterfaceTemplateV7ById(id string) (*OspfInterfaceTemplateV7, error) {
	res, err := c.Find(&OspfInterfaceTemplateV7{Id: id})
	if err != nil {
		return nil, err
	}

	return res.(*OspfInterfaceTemplateV7), nil
}

func (c Mikrotik) UpdateOspfInterfaceTemplateV7(d *OspfInterfaceTemplateV7) (*OspfInterfaceTemplateV7, error) {
	res, err := c.Update(d)
	if err != nil {
		return nil, err
	}

	return res.(*OspfInterfaceTemplateV7), nil
}

func (c Mikrotik) DeleteOspfInterfaceTemplateV7(id string) error {
	d, err := c.FindOspfInterfaceTemplateV7ById(id)
	if err != nil {
		return err
	}

	err = c.Delete(d)
	return err
}
