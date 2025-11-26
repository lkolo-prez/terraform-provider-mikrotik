package client

import "github.com/go-routeros/routeros/v3"

// SystemLogging represents RouterOS /system/logging configuration
// Path: /system/logging
//
// System logging routes log topics to specific actions (destinations).
// Each logging rule specifies which topics to log and where to send them.
type SystemLogging struct {
	Id       string `mikrotik:".id" codegen:"id,mikrotikID,read"`
	Topics   string `mikrotik:"topics" codegen:"topics,required"`   // Comma-separated topics
	Action   string `mikrotik:"action" codegen:"action,required"`   // Reference to logging action
	Prefix   string `mikrotik:"prefix" codegen:"prefix"`            // Log message prefix
	Disabled bool   `mikrotik:"disabled" codegen:"disabled"`        // yes/no
}

// Resource interface implementation for SystemLogging

func (l *SystemLogging) IDField() string {
	return ".id"
}

func (l *SystemLogging) ID() string {
	return l.Id
}

func (l *SystemLogging) SetID(id string) {
	l.Id = id
}

func (l *SystemLogging) AfterAddHook(reply *routeros.Reply) {
	l.Id = reply.Done.Map["ret"]
}

func (l *SystemLogging) FindField() string {
	return ".id"
}

func (l *SystemLogging) FindFieldValue() string {
	return l.Id
}

func (l *SystemLogging) DeleteField() string {
	return "numbers"
}

func (l *SystemLogging) DeleteFieldValue() string {
	return l.Id
}

func (l *SystemLogging) ActionToCommand(action Action) string {
	return map[Action]string{
		Add:    "/system/logging/add",
		Find:   "/system/logging/print",
		Update: "/system/logging/set",
		Delete: "/system/logging/remove",
	}[action]
}

// CRUD wrappers for SystemLogging

func (c *Mikrotik) AddSystemLogging(d *SystemLogging) (*SystemLogging, error) {
	res, err := c.Add(d)
	if err != nil {
		return nil, err
	}

	return res.(*SystemLogging), nil
}

func (c *Mikrotik) UpdateSystemLogging(d *SystemLogging) (*SystemLogging, error) {
	res, err := c.Update(d)
	if err != nil {
		return nil, err
	}

	return res.(*SystemLogging), nil
}

func (c *Mikrotik) FindSystemLogging(id string) (*SystemLogging, error) {
	res, err := c.Find(&SystemLogging{Id: id})
	if err != nil {
		return nil, err
	}

	return res.(*SystemLogging), nil
}

func (c *Mikrotik) DeleteSystemLogging(id string) error {
	return c.Delete(&SystemLogging{Id: id})
}
