package client

import (
	"github.com/go-routeros/routeros/v3"
)

// InterfaceVrrp represents a VRRP interface configuration
type InterfaceVrrp struct {
	Id              string `mikrotik:".id" codegen:"id,mikrotikID"`
	Name            string `mikrotik:"name" codegen:"name,required,terraformID"`
	Interface       string `mikrotik:"interface" codegen:"interface,required"`
	Vrid            int    `mikrotik:"vrid" codegen:"vrid,required"`
	Priority        int    `mikrotik:"priority" codegen:"priority"`
	Version         int    `mikrotik:"version" codegen:"version"`
	Authentication  string `mikrotik:"authentication" codegen:"authentication"`
	Password        string `mikrotik:"password" codegen:"password"`
	Interval        string `mikrotik:"interval" codegen:"interval"`
	PreemptionMode  bool   `mikrotik:"preemption-mode" codegen:"preemption_mode"`
	V3Protocol      string `mikrotik:"v3-protocol" codegen:"v3_protocol"`
	OnBackup        string `mikrotik:"on-backup" codegen:"on_backup"`
	OnMaster        string `mikrotik:"on-master" codegen:"on_master"`
	Disabled        bool   `mikrotik:"disabled" codegen:"disabled"`
	Comment         string `mikrotik:"comment" codegen:"comment"`
	
	// Read-only status fields
	Running         bool   `mikrotik:"running" codegen:"running,computed"`
}

var _ Resource = (*InterfaceVrrp)(nil)

func (v *InterfaceVrrp) ActionToCommand(a Action) string {
	return map[Action]string{
		Add:    "/interface/vrrp/add",
		Find:   "/interface/vrrp/print",
		Update: "/interface/vrrp/set",
		Delete: "/interface/vrrp/remove",
	}[a]
}

func (v *InterfaceVrrp) IDField() string {
	return ".id"
}

func (v *InterfaceVrrp) ID() string {
	return v.Id
}

func (v *InterfaceVrrp) SetID(id string) {
	v.Id = id
}

func (v *InterfaceVrrp) AfterAddHook(r *routeros.Reply) {
	v.Id = r.Done.Map["ret"]
}

func (v *InterfaceVrrp) FindField() string {
	return "name"
}

func (v *InterfaceVrrp) FindFieldValue() string {
	return v.Name
}

func (v *InterfaceVrrp) DeleteField() string {
	return "numbers"
}

func (v *InterfaceVrrp) DeleteFieldValue() string {
	return v.Name
}

// Typed wrappers
func (c Mikrotik) AddInterfaceVrrp(r *InterfaceVrrp) (*InterfaceVrrp, error) {
	res, err := c.Add(r)
	if err != nil {
		return nil, err
	}

	return res.(*InterfaceVrrp), nil
}

func (c Mikrotik) UpdateInterfaceVrrp(r *InterfaceVrrp) (*InterfaceVrrp, error) {
	res, err := c.Update(r)
	if err != nil {
		return nil, err
	}

	return res.(*InterfaceVrrp), nil
}

func (c Mikrotik) FindInterfaceVrrp(name string) (*InterfaceVrrp, error) {
	res, err := c.Find(&InterfaceVrrp{Name: name})
	if err != nil {
		return nil, err
	}

	return res.(*InterfaceVrrp), nil
}

func (c Mikrotik) DeleteInterfaceVrrp(name string) error {
	return c.Delete(&InterfaceVrrp{Name: name})
}
