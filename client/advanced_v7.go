package client

import (
	"fmt"
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

// WiFiConfiguration represents RouterOS v7 WiFi configuration profile
type WiFiConfiguration struct {
	Id           string `mikrotik:".id" codegen:"id,mikrotikID"`
	Name         string `mikrotik:"name" codegen:"name,required"`
	Mode         string `mikrotik:"mode" codegen:"mode,required"` // ap, station, etc.
	Ssid         string `mikrotik:"ssid" codegen:"ssid,required"`
	Country      string `mikrotik:"country" codegen:"country"`
	Security     string `mikrotik:"security" codegen:"security"`           // Reference to security profile
	Channel      string `mikrotik:"channel" codegen:"channel"`             // Reference to channel config
	Datapath     string `mikrotik:"datapath" codegen:"datapath"`           // Bridge name
	HideSSID     string `mikrotik:"hide-ssid" codegen:"hide_ssid"`
	Disabled     string `mikrotik:"disabled" codegen:"disabled"`
	Comment      string `mikrotik:"comment" codegen:"comment"`
}

// WiFiSecurity represents RouterOS v7 WiFi security profile
type WiFiSecurity struct {
	Id                   string `mikrotik:".id" codegen:"id,mikrotikID"`
	Name                 string `mikrotik:"name" codegen:"name,required"`
	AuthenticationTypes  string `mikrotik:"authentication-types" codegen:"authentication_types"` // wpa2-psk, wpa3-psk, etc.
	Passphrase           string `mikrotik:"passphrase" codegen:"passphrase"`
	GroupCiphers         string `mikrotik:"group-ciphers" codegen:"group_ciphers"`     // aes-ccm, etc.
	UnicastCiphers       string `mikrotik:"unicast-ciphers" codegen:"unicast_ciphers"` // aes-ccm, etc.
	Pmf                  string `mikrotik:"pmf" codegen:"pmf"`                         // optional, required
	Wps                  string `mikrotik:"wps" codegen:"wps"`                         // enable/disable
	Ft                   string `mikrotik:"ft" codegen:"ft"`                           // Fast Transition
	FtPreserveVlanId     string `mikrotik:"ft-preserve-vlan-id" codegen:"ft_preserve_vlan_id"`
	Disabled             string `mikrotik:"disabled" codegen:"disabled"`
	Comment              string `mikrotik:"comment" codegen:"comment"`
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
	return "/interface/veth" + action.ToCommand()
}

// ActionToCommand returns the RouterOS CLI path for WiFiRadio
func (radio *WiFiRadio) ActionToCommand(action Action) string {
	return "/interface/wifi/radio" + action.ToCommand()
}

// ActionToCommand returns the RouterOS CLI path for WiFiConfiguration
func (cfg *WiFiConfiguration) ActionToCommand(action Action) string {
	return "/interface/wifi/configuration" + action.ToCommand()
}

// ActionToCommand returns the RouterOS CLI path for WiFiSecurity
func (sec *WiFiSecurity) ActionToCommand(action Action) string {
	return "/interface/wifi/security" + action.ToCommand()
}

// ActionToCommand returns the RouterOS CLI path for QueueType
func (qt *QueueType) ActionToCommand(action Action) string {
	return "/queue/type" + action.ToCommand()
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
	err = Unmarshal(*reply.Re[0], veth)
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

	qtype := &QueueType{}
	err = Unmarshal(*reply.Re[0], qtype)
	if err != nil {
		return nil, err
	}

	return qtype, nil
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
