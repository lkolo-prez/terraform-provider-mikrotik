package client

// BgpSession represents a read-only BGP session monitoring resource in RouterOS v7+
// This menu shows cached BGP session information including status, capabilities,
// and negotiated parameters.
//
// Even if the BGP session is not active anymore, the cache can still be stored for
// some time. Routes received from a particular session are removed only if the
// cache expires.
//
// Reference: https://help.mikrotik.com/docs/display/ROS/BGP#BGP-SessionMenu
type BgpSession struct {
	Id                    string `mikrotik:".id"`
	Name                  string `mikrotik:"name"`
	Established           bool   `mikrotik:"established"`
	
	// Remote peer information
	RemoteAddress         string `mikrotik:"remote.address"`
	RemoteAS              int    `mikrotik:"remote.as"`
	RemoteID              string `mikrotik:"remote.id"`
	RemoteCapabilities    string `mikrotik:"remote.capabilities"`
	RemoteAFI             string `mikrotik:"remote.afi"`
	RemoteMessages        int    `mikrotik:"remote.messages"`
	RemoteBytes           int    `mikrotik:"remote.bytes"`
	RemoteEOR             string `mikrotik:"remote.eor"`
	RemoteRefusedCapOpt   bool   `mikrotik:"remote.refused-cap-opt"`
	
	// Local information
	LocalAddress          string `mikrotik:"local.address"`
	LocalAS               int    `mikrotik:"local.as"`
	LocalID               string `mikrotik:"local.id"`
	LocalCapabilities     string `mikrotik:"local.capabilities"`
	LocalMessages         int    `mikrotik:"local.messages"`
	LocalBytes            int    `mikrotik:"local.bytes"`
	LocalEOR              string `mikrotik:"local.eor"`
	
	// Session timers and status
	HoldTime              string `mikrotik:"hold-time"`
	KeepaliveTime         string `mikrotik:"keepalive-time"`
	Uptime                string `mikrotik:"uptime"`
	
	// Process information
	OutputProcID          int    `mikrotik:"output.procid"`
	OutputKeepSentAttrs   bool   `mikrotik:"output.keep-sent-attributes"`
	OutputLastNotification string `mikrotik:"output.last-notification"`
	
	InputProcID           int    `mikrotik:"input.procid"`
	InputLimitProcessRoutes int  `mikrotik:"input.limit-process-routes"`
	
	// Session state
	State                 string `mikrotik:"state"`
	
	// Statistics
	PrefixCount           int    `mikrotik:"prefix-count"`
}

var _ Resource = (*BgpSession)(nil)

func (b *BgpSession) ActionToCommand(a Action) string {
	// BGP Session is read-only, only print command is supported
	return map[Action]string{
		Find: "/routing/bgp/session/print",
	}[a]
}

func (b *BgpSession) IDField() string {
	return ".id"
}

func (b *BgpSession) ID() string {
	return b.Id
}

func (b *BgpSession) SetID(id string) {
	b.Id = id
}

func (b *BgpSession) FindField() string {
	return "name"
}

func (b *BgpSession) FindFieldValue() string {
	return b.Name
}

// Read-only methods - no Add/Update/Delete support

func (c Mikrotik) FindBgpSession(name string) (*BgpSession, error) {
	res, err := c.Find(&BgpSession{Name: name})
	if err != nil {
		return nil, err
	}

	return res.(*BgpSession), nil
}

func (client Mikrotik) ListBgpSessions() ([]*BgpSession, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := []string{"/routing/bgp/session/print"}
	reply, err := c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}

	sessions := make([]*BgpSession, 0)
	err = Unmarshal(*reply, &sessions)
	if err != nil {
		return nil, err
	}

	return sessions, nil
}
