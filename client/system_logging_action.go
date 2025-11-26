package client

import "github.com/go-routeros/routeros/v3"

// SystemLoggingAction represents RouterOS /system/logging/action configuration
// Path: /system/logging/action
//
// Logging actions define where and how logs are stored or transmitted.
// Actions can target memory, disk, email, or remote syslog servers.
type SystemLoggingAction struct {
	Id         string `mikrotik:".id" codegen:"id,mikrotikID,read"`
	Name       string `mikrotik:"name" codegen:"name,required"`
	Target     string `mikrotik:"target" codegen:"target,required"`     // disk, echo, email, memory, remote
	Remote     string `mikrotik:"remote" codegen:"remote"`              // IP:port for remote syslog
	RemotePort string `mikrotik:"remote-port" codegen:"remote_port"`    // Remote port (deprecated, use remote)
	BsdSyslog  bool   `mikrotik:"bsd-syslog" codegen:"bsd_syslog"`      // yes/no - BSD syslog format
	SyslogFacility string `mikrotik:"syslog-facility" codegen:"syslog_facility"` // kern, user, mail, daemon, auth, syslog, etc.
	SrcAddress string `mikrotik:"src-address" codegen:"src_address"`    // Source IP for remote logging
	Memory     string `mikrotik:"memory" codegen:"memory"`              // Memory lines (for memory target)
	DiskFileName string `mikrotik:"disk-file-name" codegen:"disk_file_name"` // File name for disk target
	DiskFileCount string `mikrotik:"disk-file-count" codegen:"disk_file_count"` // Number of files for rotation
	DiskLinesPerFile string `mikrotik:"disk-lines-per-file" codegen:"disk_lines_per_file"` // Lines per file
	Remember   bool   `mikrotik:"remember" codegen:"remember"`          // yes/no - remember logs
}

// Resource interface implementation for SystemLoggingAction

func (a *SystemLoggingAction) IDField() string {
	return ".id"
}

func (a *SystemLoggingAction) ID() string {
	return a.Id
}

func (a *SystemLoggingAction) SetID(id string) {
	a.Id = id
}

func (a *SystemLoggingAction) AfterAddHook(reply *routeros.Reply) {
	a.Id = reply.Done.Map["ret"]
}

func (a *SystemLoggingAction) FindField() string {
	return "name"
}

func (a *SystemLoggingAction) FindFieldValue() string {
	return a.Name
}

func (a *SystemLoggingAction) DeleteField() string {
	return "numbers"
}

func (a *SystemLoggingAction) DeleteFieldValue() string {
	return a.Id
}

func (a *SystemLoggingAction) ActionToCommand(action Action) string {
	return map[Action]string{
		Add:    "/system/logging/action/add",
		Find:   "/system/logging/action/print",
		Update: "/system/logging/action/set",
		Delete: "/system/logging/action/remove",
	}[action]
}

// CRUD wrappers for SystemLoggingAction

func (c *Mikrotik) AddSystemLoggingAction(d *SystemLoggingAction) (*SystemLoggingAction, error) {
	res, err := c.Add(d)
	if err != nil {
		return nil, err
	}

	return res.(*SystemLoggingAction), nil
}

func (c *Mikrotik) UpdateSystemLoggingAction(d *SystemLoggingAction) (*SystemLoggingAction, error) {
	res, err := c.Update(d)
	if err != nil {
		return nil, err
	}

	return res.(*SystemLoggingAction), nil
}

func (c *Mikrotik) FindSystemLoggingAction(name string) (*SystemLoggingAction, error) {
	res, err := c.Find(&SystemLoggingAction{Name: name})
	if err != nil {
		return nil, err
	}

	return res.(*SystemLoggingAction), nil
}

func (c *Mikrotik) DeleteSystemLoggingAction(id string) error {
	return c.Delete(&SystemLoggingAction{Id: id})
}
