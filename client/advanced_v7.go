package client

import (
	"fmt"
	"github.com/go-routeros/routeros/v3"
)

// InterfaceVeth represents RouterOS v7 Virtual Ethernet interface
// Reference: https://help.mikrotik.com/docs/display/ROS/Virtual+Ethernet
type InterfaceVeth struct {
	Id       string `mikrotik:".id" codegen:"id,mikrotikID"`
	Name     string `mikrotik:"name" codegen:"name,required"`
	Address  string `mikrotik:"address" codegen:"address,required"` // IP address with netmask
	Gateway  string `mikrotik:"gateway" codegen:"gateway"`
	Mtu      string `mikrotik:"mtu" codegen:"mtu"`
	Disabled string `mikrotik:"disabled" codegen:"disabled"`
	Comment  string `mikrotik:"comment" codegen:"comment"`
}

// WiFiRadio represents RouterOS v7 WiFi radio configuration
// Reference: https://help.mikrotik.com/docs/display/ROS/WiFi
type WiFiRadio struct {
	Id             string `mikrotik:".id" codegen:"id,mikrotikID"`
	Name           string `mikrotik:"name" codegen:"name,required"`
	Band           string `mikrotik:"band" codegen:"band"`                       // 2ghz-ax, 5ghz-ax, 6ghz-ax, etc.
	ChannelWidth   string `mikrotik:"channel-width" codegen:"channel_width"`     // 20mhz, 40mhz, 80mhz, 160mhz
	Frequency      string `mikrotik:"frequency" codegen:"frequency"`             // Center frequency
	SkipDfsChannels string `mikrotik:"skip-dfs-channels" codegen:"skip_dfs_channels"`
	Disabled       string `mikrotik:"disabled" codegen:"disabled"`
	Comment        string `mikrotik:"comment" codegen:"comment"`
}

// InterfaceWiFi represents RouterOS v7 /interface/wifi resource
type InterfaceWiFi struct {
	Id              string `mikrotik:".id" codegen:"id,mikrotikID"`
	Name            string `mikrotik:"name" codegen:"name,required"`
	Configuration   string `mikrotik:"configuration" codegen:"configuration"`
	Datapath        string `mikrotik:"datapath" codegen:"datapath"`
	Channel         string `mikrotik:"channel" codegen:"channel"`
	MasterInterface string `mikrotik:"master-interface" codegen:"master_interface"`
	MacAddress      string `mikrotik:"mac-address" codegen:"mac_address"`
	Mtu             string `mikrotik:"mtu" codegen:"mtu"`
	Arp             string `mikrotik:"arp" codegen:"arp"`
	ArpTimeout      string `mikrotik:"arp-timeout" codegen:"arp_timeout"`
	Disabled        string `mikrotik:"disabled" codegen:"disabled"`
	Comment         string `mikrotik:"comment" codegen:"comment"`
	Running         string `mikrotik:"running" codegen:"running,computed"`
	Radio           string `mikrotik:"radio" codegen:"radio,computed"`
}

// WiFiConfiguration represents RouterOS v7 WiFi configuration profile
type WiFiConfiguration struct {
	Id                    string `mikrotik:".id" codegen:"id,mikrotikID"`
	Name                  string `mikrotik:"name" codegen:"name,required"`
	Mode                  string `mikrotik:"mode" codegen:"mode"` // ap, station, sniffer
	SSID                  string `mikrotik:"ssid" codegen:"ssid"`
	HideSSID              string `mikrotik:"hide-ssid" codegen:"hide_ssid"`
	Security              string `mikrotik:"security" codegen:"security"`
	Country               string `mikrotik:"country" codegen:"country"`
	Installation          string `mikrotik:"installation" codegen:"installation"`
	TxPower               string `mikrotik:"tx-power" codegen:"tx_power"`
	TxPowerMode           string `mikrotik:"tx-power-mode" codegen:"tx_power_mode"`
	SupportedRates        string `mikrotik:"supported-rates" codegen:"supported_rates"`
	BasicRates            string `mikrotik:"basic-rates" codegen:"basic_rates"`
	HEGuardInterval       string `mikrotik:"he-guard-interval" codegen:"he_guard_interval"`
	HEFrameFormat         string `mikrotik:"he-frame-format" codegen:"he_frame_format"`
	Distance              string `mikrotik:"distance" codegen:"distance"`
	AckTimeout            string `mikrotik:"ack-timeout" codegen:"ack_timeout"`
	MulticastHelper       string `mikrotik:"multicast-helper" codegen:"multicast_helper"`
	LoadBalancingGroup    string `mikrotik:"load-balancing-group" codegen:"load_balancing_group"`
	ConnectList           string `mikrotik:"connect-list" codegen:"connect_list"`
	ScanList              string `mikrotik:"scan-list" codegen:"scan_list"`
	IEEE80211w            string `mikrotik:"ieee80211w" codegen:"ieee80211w"`
	WPS                   string `mikrotik:"wps" codegen:"wps"`
	WPSMode               string `mikrotik:"wps-mode" codegen:"wps_mode"`
	BridgeMode            string `mikrotik:"bridge-mode" codegen:"bridge_mode"`
	Steering              string `mikrotik:"steering" codegen:"steering"`
	DisconnectTimeout     string `mikrotik:"disconnect-timeout" codegen:"disconnect_timeout"`
	KeepaliveFrames       string `mikrotik:"keepalive-frames" codegen:"keepalive_frames"`
	MaxStationCount       string `mikrotik:"max-station-count" codegen:"max_station_count"`
	Roaming               string `mikrotik:"roaming" codegen:"roaming"`
	RoamingThreshold      string `mikrotik:"roaming-threshold" codegen:"roaming_threshold"`
	BeaconInterval        string `mikrotik:"beacon-interval" codegen:"beacon_interval"`
	DTIMPeriod            string `mikrotik:"dtim-period" codegen:"dtim_period"`
	Disabled              string `mikrotik:"disabled" codegen:"disabled"`
	Comment               string `mikrotik:"comment" codegen:"comment"`
}

// WiFiSecurity represents RouterOS v7 WiFi security profile
type WiFiSecurity struct {
	Id                         string `mikrotik:".id" codegen:"id,mikrotikID"`
	Name                       string `mikrotik:"name" codegen:"name,required"`
	AuthenticationTypes        string `mikrotik:"authentication-types" codegen:"authentication_types"`
	Encryption                 string `mikrotik:"encryption" codegen:"encryption"`
	Passphrase                 string `mikrotik:"passphrase" codegen:"passphrase"`
	ManagementProtection       string `mikrotik:"management-protection" codegen:"management_protection"`
	ManagementProtectionKey    string `mikrotik:"management-protection-key" codegen:"management_protection_key"`
	EapMethods                 string `mikrotik:"eap-methods" codegen:"eap_methods"`
	EapRadiusServer            string `mikrotik:"eap-radius-server" codegen:"eap_radius_server"`
	EapRadiusSecret            string `mikrotik:"eap-radius-secret" codegen:"eap_radius_secret"`
	EapRadiusPort              string `mikrotik:"eap-radius-port" codegen:"eap_radius_port"`
	EapRadiusAccounting        string `mikrotik:"eap-radius-accounting" codegen:"eap_radius_accounting"`
	EapRadiusAccountingPort    string `mikrotik:"eap-radius-accounting-port" codegen:"eap_radius_accounting_port"`
	EapTlsCertificate          string `mikrotik:"eap-tls-certificate" codegen:"eap_tls_certificate"`
	SAE_PWE                    string `mikrotik:"sae-pwe" codegen:"sae_pwe"`
	SAE_Groups                 string `mikrotik:"sae-groups" codegen:"sae_groups"`
	FT                         string `mikrotik:"ft" codegen:"ft"`
	FTOverDS                   string `mikrotik:"ft-over-ds" codegen:"ft_over_ds"`
	FTPreserveVlan             string `mikrotik:"ft-preserve-vlan" codegen:"ft_preserve_vlan"`
	GroupRekey                 string `mikrotik:"group-rekey" codegen:"group_rekey"`
	Disabled                   string `mikrotik:"disabled" codegen:"disabled"`
	Comment                    string `mikrotik:"comment" codegen:"comment"`
}

// WiFiChannel represents RouterOS v7 /interface/wifi/channel resource
type WiFiChannel struct {
	Id                     string `mikrotik:".id" codegen:"id,mikrotikID"`
	Name                   string `mikrotik:"name" codegen:"name,required"`
	Band                   string `mikrotik:"band" codegen:"band"`
	Frequency              string `mikrotik:"frequency" codegen:"frequency"`
	Width                  string `mikrotik:"width" codegen:"width"`
	SecondaryFrequency     string `mikrotik:"secondary-frequency" codegen:"secondary_frequency"`
	SkipDFSChannels        string `mikrotik:"skip-dfs-channels" codegen:"skip_dfs_channels"`
	ReuseDFSChannels       string `mikrotik:"reuse-dfs-channels" codegen:"reuse_dfs_channels"`
	ControlChannelPosition string `mikrotik:"control-channel-position" codegen:"control_channel_position"`
	Disabled               string `mikrotik:"disabled" codegen:"disabled"`
	Comment                string `mikrotik:"comment" codegen:"comment"`
}

// WiFiDatapath represents RouterOS v7 /interface/wifi/datapath resource
type WiFiDatapath struct {
	Id                        string `mikrotik:".id" codegen:"id,mikrotikID"`
	Name                      string `mikrotik:"name" codegen:"name,required"`
	Bridge                    string `mikrotik:"bridge" codegen:"bridge"`
	BridgeHorizon             string `mikrotik:"bridge-horizon" codegen:"bridge_horizon"`
	BridgeCost                string `mikrotik:"bridge-cost" codegen:"bridge_cost"`
	VlanID                    string `mikrotik:"vlan-id" codegen:"vlan_id"`
	VlanMode                  string `mikrotik:"vlan-mode" codegen:"vlan_mode"`
	ClientIsolation           string `mikrotik:"client-isolation" codegen:"client_isolation"`
	ClientToClientForwarding  string `mikrotik:"client-to-client-forwarding" codegen:"client_to_client_forwarding"`
	ARP                       string `mikrotik:"arp" codegen:"arp"`
	ARPTimeout                string `mikrotik:"arp-timeout" codegen:"arp_timeout"`
	InterfaceList             string `mikrotik:"interface-list" codegen:"interface_list"`
	L2MTU                     string `mikrotik:"l2mtu" codegen:"l2mtu"`
	MTU                       string `mikrotik:"mtu" codegen:"mtu"`
	OpenFlowSwitch            string `mikrotik:"openflow-switch" codegen:"openflow_switch"`
	Disabled                  string `mikrotik:"disabled" codegen:"disabled"`
	Comment                   string `mikrotik:"comment" codegen:"comment"`
}

// WiFiAccessList represents RouterOS v7 /interface/wifi/access-list resource
type WiFiAccessList struct {
	Id               string `mikrotik:".id" codegen:"id,mikrotikID"`
	MacAddress       string `mikrotik:"mac-address" codegen:"mac_address,required"`
	Action           string `mikrotik:"action" codegen:"action"`
	Interface        string `mikrotik:"interface" codegen:"interface"`
	SSIDRegexp       string `mikrotik:"ssid-regexp" codegen:"ssid_regexp"`
	VlanID           string `mikrotik:"vlan-id" codegen:"vlan_id"`
	VlanMode         string `mikrotik:"vlan-mode" codegen:"vlan_mode"`
	SignalRange      string `mikrotik:"signal-range" codegen:"signal_range"`
	Time             string `mikrotik:"time" codegen:"time"`
	ClientIsolation  string `mikrotik:"client-isolation" codegen:"client_isolation"`
	RadiusAccounting string `mikrotik:"radius-accounting" codegen:"radius_accounting"`
	Disabled         string `mikrotik:"disabled" codegen:"disabled"`
	Comment          string `mikrotik:"comment" codegen:"comment"`
}

// QueueType represents RouterOS v7 queue types including CAKE and fq_codel
// Reference: https://help.mikrotik.com/docs/display/ROS/Queues
type QueueType struct {
	Id         string `mikrotik:".id" codegen:"id,mikrotikID"`
	Name       string `mikrotik:"name" codegen:"name,required"`
	Kind       string `mikrotik:"kind" codegen:"kind,required"` // pcq, red, sfq, fifo, cake, fq_codel
	
	// PCQ-specific
	PcqRate           string `mikrotik:"pcq-rate" codegen:"pcq_rate"`
	PcqClassifier     string `mikrotik:"pcq-classifier" codegen:"pcq_classifier"`
	PcqLimit          string `mikrotik:"pcq-limit" codegen:"pcq_limit"`
	PcqBurstRate      string `mikrotik:"pcq-burst-rate" codegen:"pcq_burst_rate"`
	PcqBurstThreshold string `mikrotik:"pcq-burst-threshold" codegen:"pcq_burst_threshold"`
	PcqBurstTime      string `mikrotik:"pcq-burst-time" codegen:"pcq_burst_time"`
	
	// CAKE-specific (RouterOS 7+)
	CakeBandwidth     string `mikrotik:"cake-bandwidth" codegen:"cake_bandwidth"`
	CakeRtt           string `mikrotik:"cake-rtt" codegen:"cake_rtt"`
	CakeOverhead      string `mikrotik:"cake-overhead" codegen:"cake_overhead"`
	CakeMpu           string `mikrotik:"cake-mpu" codegen:"cake_mpu"`
	CakeAtm           string `mikrotik:"cake-atm" codegen:"cake_atm"`
	CakeNat           string `mikrotik:"cake-nat" codegen:"cake_nat"`
	CakeAckFilter     string `mikrotik:"cake-ack-filter" codegen:"cake_ack_filter"`
	
	// fq_codel-specific (RouterOS 7+)
	FqCodelLimit      string `mikrotik:"fq-codel-limit" codegen:"fq_codel_limit"`
	FqCodelTarget     string `mikrotik:"fq-codel-target" codegen:"fq_codel_target"`
	FqCodelInterval   string `mikrotik:"fq-codel-interval" codegen:"fq_codel_interval"`
	FqCodelQuantum    string `mikrotik:"fq-codel-quantum" codegen:"fq_codel_quantum"`
	FqCodelEcn        string `mikrotik:"fq-codel-ecn" codegen:"fq_codel_ecn"`
}

// ActionToCommand returns the RouterOS CLI path for InterfaceVeth
func (veth *InterfaceVeth) ActionToCommand(action Action) string {
	return map[Action]string{
		Add:    "/interface/veth/add",
		Find:   "/interface/veth/print",
		Update: "/interface/veth/set",
		Delete: "/interface/veth/remove",
	}[action]
}

// ActionToCommand returns the RouterOS CLI path for WiFiRadio
func (radio *WiFiRadio) ActionToCommand(action Action) string {
	return map[Action]string{
		Add:    "/interface/wifi/radio/add",
		Find:   "/interface/wifi/radio/print",
		Update: "/interface/wifi/radio/set",
		Delete: "/interface/wifi/radio/remove",
	}[action]
}

// ActionToCommand returns the RouterOS CLI path for InterfaceWiFi
func (w *InterfaceWiFi) ActionToCommand(action Action) string {
	return map[Action]string{
		Add:    "/interface/wifi/add",
		Find:   "/interface/wifi/print",
		Update: "/interface/wifi/set",
		Delete: "/interface/wifi/remove",
	}[action]
}

// ActionToCommand returns the RouterOS CLI path for WiFiConfiguration
func (cfg *WiFiConfiguration) ActionToCommand(action Action) string {
	return map[Action]string{
		Add:    "/interface/wifi/configuration/add",
		Find:   "/interface/wifi/configuration/print",
		Update: "/interface/wifi/configuration/set",
		Delete: "/interface/wifi/configuration/remove",
	}[action]
}

// ActionToCommand returns the RouterOS CLI path for WiFiSecurity
func (sec *WiFiSecurity) ActionToCommand(action Action) string {
	return map[Action]string{
		Add:    "/interface/wifi/security/add",
		Find:   "/interface/wifi/security/print",
		Update: "/interface/wifi/security/set",
		Delete: "/interface/wifi/security/remove",
	}[action]
}

// ActionToCommand returns the RouterOS CLI path for WiFiChannel
func (ch *WiFiChannel) ActionToCommand(action Action) string {
	return map[Action]string{
		Add:    "/interface/wifi/channel/add",
		Find:   "/interface/wifi/channel/print",
		Update: "/interface/wifi/channel/set",
		Delete: "/interface/wifi/channel/remove",
	}[action]
}

// ActionToCommand returns the RouterOS CLI path for WiFiDatapath
func (d *WiFiDatapath) ActionToCommand(action Action) string {
	return map[Action]string{
		Add:    "/interface/wifi/datapath/add",
		Find:   "/interface/wifi/datapath/print",
		Update: "/interface/wifi/datapath/set",
		Delete: "/interface/wifi/datapath/remove",
	}[action]
}

// ActionToCommand returns the RouterOS CLI path for WiFiAccessList
func (a *WiFiAccessList) ActionToCommand(action Action) string {
	return map[Action]string{
		Add:    "/interface/wifi/access-list/add",
		Find:   "/interface/wifi/access-list/print",
		Update: "/interface/wifi/access-list/set",
		Delete: "/interface/wifi/access-list/remove",
	}[action]
}

// ActionToCommand returns the RouterOS CLI path for QueueType
func (qt *QueueType) ActionToCommand(action Action) string {
	return map[Action]string{
		Add:    "/queue/type/add",
		Find:   "/queue/type/print",
		Update: "/queue/type/set",
		Delete: "/queue/type/remove",
	}[action]
}

// FindInterfaceVeth finds a veth interface by name
func (client Mikrotik) FindInterfaceVeth(name string) (*InterfaceVeth, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := []string{"/interface/veth/print", "?name=" + name}
	reply, err := c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}

	if len(reply.Re) == 0 {
		return nil, NewNotFound(fmt.Sprintf("veth interface '%s' not found", name))
	}

	veth := &InterfaceVeth{}
	err = Unmarshal(*reply, veth)
	if err != nil {
		return nil, err
	}

	return veth, nil
}

// CreateInterfaceVeth creates a new veth interface
func (client Mikrotik) CreateInterfaceVeth(veth *InterfaceVeth) (*InterfaceVeth, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := Marshal("/interface/veth/add", veth)
	_, err = c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}

	return client.FindInterfaceVeth(veth.Name)
}

// UpdateInterfaceVeth updates an existing veth interface
func (client Mikrotik) UpdateInterfaceVeth(veth *InterfaceVeth) (*InterfaceVeth, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := Marshal("/interface/veth/set", veth)
	_, err = c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}

	return client.FindInterfaceVeth(veth.Name)
}

// DeleteInterfaceVeth deletes a veth interface
func (client Mikrotik) DeleteInterfaceVeth(name string) error {
	c, err := client.getMikrotikClient()
	if err != nil {
		return err
	}

	veth, err := client.FindInterfaceVeth(name)
	if err != nil {
		return err
	}

	cmd := []string{"/interface/veth/remove", "=.id=" + veth.Id}
	_, err = c.RunArgs(cmd)
	return err
}

// FindQueueType finds a queue type by name
func (client Mikrotik) FindQueueType(name string) (*QueueType, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := []string{"/queue/type/print", "?name=" + name}
	reply, err := c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}

	if len(reply.Re) == 0 {
		return nil, NewNotFound(fmt.Sprintf("queue type '%s' not found", name))
	}

	qt := &QueueType{}
	err = Unmarshal(*reply, qt)
	if err != nil {
		return nil, err
	}

	return qt, nil
}

// CreateQueueType creates a new queue type
func (client Mikrotik) CreateQueueType(qtype *QueueType) (*QueueType, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := Marshal("/queue/type/add", qtype)
	_, err = c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}

	return client.FindQueueType(qtype.Name)
}

// UpdateQueueType updates an existing queue type
func (client Mikrotik) UpdateQueueType(qtype *QueueType) (*QueueType, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := Marshal("/queue/type/set", qtype)
	_, err = c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}

	return client.FindQueueType(qtype.Name)
}

// DeleteQueueType deletes a queue type
func (client Mikrotik) DeleteQueueType(name string) error {
	c, err := client.getMikrotikClient()
	if err != nil {
		return err
	}

	qtype, err := client.FindQueueType(name)
	if err != nil {
		return err
	}

	cmd := []string{"/queue/type/remove", "=.id=" + qtype.Id}
	_, err = c.RunArgs(cmd)
	return err
}

// Resource interface implementations for WiFi 6 resources

var _ Resource = (*InterfaceWiFi)(nil)
var _ Resource = (*WiFiConfiguration)(nil)
var _ Resource = (*WiFiSecurity)(nil)
var _ Resource = (*WiFiChannel)(nil)
var _ Resource = (*WiFiDatapath)(nil)
var _ Resource = (*WiFiAccessList)(nil)

// InterfaceWiFi Resource implementation
func (w *InterfaceWiFi) IDField() string        { return ".id" }
func (w *InterfaceWiFi) ID() string             { return w.Id }
func (w *InterfaceWiFi) SetID(id string)        { w.Id = id }
func (w *InterfaceWiFi) AfterAddHook(r *routeros.Reply)  { w.Id = r.Done.Map["ret"] }
func (w *InterfaceWiFi) FindField() string      { return "name" }
func (w *InterfaceWiFi) FindFieldValue() string { return w.Name }
func (w *InterfaceWiFi) DeleteField() string    { return "numbers" }
func (w *InterfaceWiFi) DeleteFieldValue() string { return w.Name }

// WiFiConfiguration Resource implementation
func (w *WiFiConfiguration) IDField() string        { return ".id" }
func (w *WiFiConfiguration) ID() string             { return w.Id }
func (w *WiFiConfiguration) SetID(id string)        { w.Id = id }
func (w *WiFiConfiguration) AfterAddHook(r *routeros.Reply)  { w.Id = r.Done.Map["ret"] }
func (w *WiFiConfiguration) FindField() string      { return "name" }
func (w *WiFiConfiguration) FindFieldValue() string { return w.Name }
func (w *WiFiConfiguration) DeleteField() string    { return "numbers" }
func (w *WiFiConfiguration) DeleteFieldValue() string { return w.Name }

// WiFiSecurity Resource implementation
func (w *WiFiSecurity) IDField() string        { return ".id" }
func (w *WiFiSecurity) ID() string             { return w.Id }
func (w *WiFiSecurity) SetID(id string)        { w.Id = id }
func (w *WiFiSecurity) AfterAddHook(r *routeros.Reply)  { w.Id = r.Done.Map["ret"] }
func (w *WiFiSecurity) FindField() string      { return "name" }
func (w *WiFiSecurity) FindFieldValue() string { return w.Name }
func (w *WiFiSecurity) DeleteField() string    { return "numbers" }
func (w *WiFiSecurity) DeleteFieldValue() string { return w.Name }

// WiFiChannel Resource implementation
func (w *WiFiChannel) IDField() string        { return ".id" }
func (w *WiFiChannel) ID() string             { return w.Id }
func (w *WiFiChannel) SetID(id string)        { w.Id = id }
func (w *WiFiChannel) AfterAddHook(r *routeros.Reply)  { w.Id = r.Done.Map["ret"] }
func (w *WiFiChannel) FindField() string      { return "name" }
func (w *WiFiChannel) FindFieldValue() string { return w.Name }
func (w *WiFiChannel) DeleteField() string    { return "numbers" }
func (w *WiFiChannel) DeleteFieldValue() string { return w.Name }

// WiFiDatapath Resource implementation
func (w *WiFiDatapath) IDField() string        { return ".id" }
func (w *WiFiDatapath) ID() string             { return w.Id }
func (w *WiFiDatapath) SetID(id string)        { w.Id = id }
func (w *WiFiDatapath) AfterAddHook(r *routeros.Reply)  { w.Id = r.Done.Map["ret"] }
func (w *WiFiDatapath) FindField() string      { return "name" }
func (w *WiFiDatapath) FindFieldValue() string { return w.Name }
func (w *WiFiDatapath) DeleteField() string    { return "numbers" }
func (w *WiFiDatapath) DeleteFieldValue() string { return w.Name }

// WiFiAccessList Resource implementation
func (w *WiFiAccessList) IDField() string        { return ".id" }
func (w *WiFiAccessList) ID() string             { return w.Id }
func (w *WiFiAccessList) SetID(id string)        { w.Id = id }
func (w *WiFiAccessList) AfterAddHook(r *routeros.Reply)  { w.Id = r.Done.Map["ret"] }
func (w *WiFiAccessList) FindField() string      { return "mac-address" }
func (w *WiFiAccessList) FindFieldValue() string { return w.MacAddress }
func (w *WiFiAccessList) DeleteField() string    { return "numbers" }
func (w *WiFiAccessList) DeleteFieldValue() string { return w.MacAddress }

// CRUD wrappers for InterfaceWiFi
func (client Mikrotik) AddInterfaceWiFi(w *InterfaceWiFi) (*InterfaceWiFi, error) {
	res, err := client.Add(w)
	if err != nil {
		return nil, err
	}
	return res.(*InterfaceWiFi), nil
}

func (client Mikrotik) FindInterfaceWiFi(name string) (*InterfaceWiFi, error) {
	res, err := client.Find(&InterfaceWiFi{Name: name})
	if err != nil {
		return nil, err
	}
	return res.(*InterfaceWiFi), nil
}

func (client Mikrotik) UpdateInterfaceWiFi(w *InterfaceWiFi) (*InterfaceWiFi, error) {
	res, err := client.Update(w)
	if err != nil {
		return nil, err
	}
	return res.(*InterfaceWiFi), nil
}

func (client Mikrotik) DeleteInterfaceWiFi(name string) error {
	return client.Delete(&InterfaceWiFi{Name: name})
}

// CRUD wrappers for WiFiConfiguration
func (client Mikrotik) AddWiFiConfiguration(w *WiFiConfiguration) (*WiFiConfiguration, error) {
	res, err := client.Add(w)
	if err != nil {
		return nil, err
	}
	return res.(*WiFiConfiguration), nil
}

func (client Mikrotik) FindWiFiConfiguration(name string) (*WiFiConfiguration, error) {
	res, err := client.Find(&WiFiConfiguration{Name: name})
	if err != nil {
		return nil, err
	}
	return res.(*WiFiConfiguration), nil
}

func (client Mikrotik) UpdateWiFiConfiguration(w *WiFiConfiguration) (*WiFiConfiguration, error) {
	res, err := client.Update(w)
	if err != nil {
		return nil, err
	}
	return res.(*WiFiConfiguration), nil
}

func (client Mikrotik) DeleteWiFiConfiguration(name string) error {
	return client.Delete(&WiFiConfiguration{Name: name})
}

// CRUD wrappers for WiFiSecurity
func (client Mikrotik) AddWiFiSecurity(w *WiFiSecurity) (*WiFiSecurity, error) {
	res, err := client.Add(w)
	if err != nil {
		return nil, err
	}
	return res.(*WiFiSecurity), nil
}

func (client Mikrotik) FindWiFiSecurity(name string) (*WiFiSecurity, error) {
	res, err := client.Find(&WiFiSecurity{Name: name})
	if err != nil {
		return nil, err
	}
	return res.(*WiFiSecurity), nil
}

func (client Mikrotik) UpdateWiFiSecurity(w *WiFiSecurity) (*WiFiSecurity, error) {
	res, err := client.Update(w)
	if err != nil {
		return nil, err
	}
	return res.(*WiFiSecurity), nil
}

func (client Mikrotik) DeleteWiFiSecurity(name string) error {
	return client.Delete(&WiFiSecurity{Name: name})
}

// CRUD wrappers for WiFiChannel
func (client Mikrotik) AddWiFiChannel(w *WiFiChannel) (*WiFiChannel, error) {
	res, err := client.Add(w)
	if err != nil {
		return nil, err
	}
	return res.(*WiFiChannel), nil
}

func (client Mikrotik) FindWiFiChannel(name string) (*WiFiChannel, error) {
	res, err := client.Find(&WiFiChannel{Name: name})
	if err != nil {
		return nil, err
	}
	return res.(*WiFiChannel), nil
}

func (client Mikrotik) UpdateWiFiChannel(w *WiFiChannel) (*WiFiChannel, error) {
	res, err := client.Update(w)
	if err != nil {
		return nil, err
	}
	return res.(*WiFiChannel), nil
}

func (client Mikrotik) DeleteWiFiChannel(name string) error {
	return client.Delete(&WiFiChannel{Name: name})
}

// CRUD wrappers for WiFiDatapath
func (client Mikrotik) AddWiFiDatapath(w *WiFiDatapath) (*WiFiDatapath, error) {
	res, err := client.Add(w)
	if err != nil {
		return nil, err
	}
	return res.(*WiFiDatapath), nil
}

func (client Mikrotik) FindWiFiDatapath(name string) (*WiFiDatapath, error) {
	res, err := client.Find(&WiFiDatapath{Name: name})
	if err != nil {
		return nil, err
	}
	return res.(*WiFiDatapath), nil
}

func (client Mikrotik) UpdateWiFiDatapath(w *WiFiDatapath) (*WiFiDatapath, error) {
	res, err := client.Update(w)
	if err != nil {
		return nil, err
	}
	return res.(*WiFiDatapath), nil
}

func (client Mikrotik) DeleteWiFiDatapath(name string) error {
	return client.Delete(&WiFiDatapath{Name: name})
}

// CRUD wrappers for WiFiAccessList
func (client Mikrotik) AddWiFiAccessList(w *WiFiAccessList) (*WiFiAccessList, error) {
	res, err := client.Add(w)
	if err != nil {
		return nil, err
	}
	return res.(*WiFiAccessList), nil
}

func (client Mikrotik) FindWiFiAccessList(macAddress string) (*WiFiAccessList, error) {
	res, err := client.Find(&WiFiAccessList{MacAddress: macAddress})
	if err != nil {
		return nil, err
	}
	return res.(*WiFiAccessList), nil
}

func (client Mikrotik) UpdateWiFiAccessList(w *WiFiAccessList) (*WiFiAccessList, error) {
	res, err := client.Update(w)
	if err != nil {
		return nil, err
	}
	return res.(*WiFiAccessList), nil
}

func (client Mikrotik) DeleteWiFiAccessList(macAddress string) error {
	return client.Delete(&WiFiAccessList{MacAddress: macAddress})
}
