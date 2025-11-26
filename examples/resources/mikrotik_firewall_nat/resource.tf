# Firewall NAT Examples

terraform {
  required_providers {
    mikrotik = {
      source = "lkolo-prez/mikrotik"
    }
  }
}

provider "mikrotik" {
  host     = var.router_ip
  username = var.router_username
  password = var.router_password
  tls      = true
}

variable "router_ip" {
  type = string
}

variable "router_username" {
  type    = string
  default = "admin"
}

variable "router_password" {
  type      = string
  sensitive = true
}

## Example 1: Basic Internet Sharing (Masquerade)
resource "mikrotik_firewall_nat" "masquerade_wan" {
  chain         = "srcnat"
  action        = "masquerade"
  out_interface = "ether1-wan"
  comment       = "Internet sharing - Masquerade LAN to WAN"
}

## Example 2: Masquerade Specific Network
resource "mikrotik_firewall_nat" "masquerade_lan" {
  chain         = "srcnat"
  action        = "masquerade"
  src_address   = "192.168.1.0/24"
  out_interface = "ether1-wan"
  comment       = "Masquerade LAN network only"
}

## Example 3: Port Forwarding - HTTP Server
resource "mikrotik_firewall_nat" "http_port_forward" {
  chain        = "dstnat"
  action       = "dst-nat"
  protocol     = "tcp"
  dst_port     = "80"
  in_interface = "ether1-wan"
  to_addresses = "192.168.1.100"
  to_ports     = "80"
  comment      = "Forward HTTP to internal web server"
}

## Example 4: Port Forwarding - SSH with Port Change
resource "mikrotik_firewall_nat" "ssh_custom_port" {
  chain        = "dstnat"
  action       = "dst-nat"
  protocol     = "tcp"
  dst_port     = "2222"
  in_interface = "ether1-wan"
  to_addresses = "192.168.1.50"
  to_ports     = "22"
  comment      = "SSH to internal server on custom port 2222"
}

## Example 5: Port Forwarding - RDP
resource "mikrotik_firewall_nat" "rdp_forward" {
  chain        = "dstnat"
  action       = "dst-nat"
  protocol     = "tcp"
  dst_port     = "3389"
  in_interface = "ether1-wan"
  to_addresses = "192.168.1.110"
  comment      = "RDP to Windows server"
}

## Example 6: Port Range Forwarding
resource "mikrotik_firewall_nat" "passive_ftp_forward" {
  chain        = "dstnat"
  action       = "dst-nat"
  protocol     = "tcp"
  dst_port     = "21000-21100"
  in_interface = "ether1-wan"
  to_addresses = "192.168.1.120"
  to_ports     = "21000-21100"
  comment      = "Passive FTP port range"
}

## Example 7: Multiple Ports Forwarding
resource "mikrotik_firewall_nat" "web_services_forward" {
  chain        = "dstnat"
  action       = "dst-nat"
  protocol     = "tcp"
  dst_port     = "80,443,8080"
  in_interface = "ether1-wan"
  to_addresses = "192.168.1.130"
  comment      = "Multiple web ports to server"
}

## Example 8: Hairpin NAT (Internal Access to Public IP)
resource "mikrotik_firewall_nat" "hairpin_nat" {
  chain        = "srcnat"
  action       = "masquerade"
  src_address  = "192.168.1.0/24"
  dst_address  = "192.168.1.100"
  protocol     = "tcp"
  dst_port     = "80"
  comment      = "Hairpin NAT for internal clients"
}

## Example 9: Source NAT for Specific Subnet
resource "mikrotik_firewall_nat" "guest_srcnat" {
  chain        = "srcnat"
  action       = "src-nat"
  src_address  = "10.10.10.0/24"
  to_addresses = "203.0.113.10"
  comment      = "Guest network with static source IP"
}

## Example 10: Multi-WAN Load Balancing NAT
resource "mikrotik_firewall_nat" "wan1_masquerade" {
  chain           = "srcnat"
  action          = "masquerade"
  connection_mark = "wan1-conn"
  out_interface   = "ether1-wan1"
  comment         = "WAN1 masquerade for load balancing"
}

resource "mikrotik_firewall_nat" "wan2_masquerade" {
  chain           = "srcnat"
  action          = "masquerade"
  connection_mark = "wan2-conn"
  out_interface   = "ether2-wan2"
  comment         = "WAN2 masquerade for load balancing"
}

## Example 11: DMZ with 1:1 NAT (Netmap)
resource "mikrotik_firewall_nat" "dmz_netmap" {
  chain        = "dstnat"
  action       = "netmap"
  dst_address  = "203.0.113.0/24"
  to_addresses = "192.168.100.0/24"
  comment      = "DMZ 1:1 NAT mapping"
}

## Example 12: Transparent Proxy Redirect
resource "mikrotik_firewall_nat" "proxy_redirect" {
  chain      = "dstnat"
  action     = "redirect"
  protocol   = "tcp"
  dst_port   = "80"
  to_ports   = "3128"
  comment    = "Transparent HTTP proxy redirect"
}

## Example 13: NAT Exception (No NAT for VPN)
resource "mikrotik_firewall_nat" "vpn_no_nat" {
  chain            = "srcnat"
  action           = "accept"
  src_address_list = "vpn-clients"
  dst_address_list = "local-networks"
  comment          = "No NAT for VPN to local networks"
}

## Example 14: Conditional NAT with Connection State
resource "mikrotik_firewall_nat" "established_masquerade" {
  chain            = "srcnat"
  action           = "masquerade"
  connection_state = "established,related"
  out_interface    = "ether1-wan"
  comment          = "Masquerade only established connections"
}

## Example 15: NAT with Logging
resource "mikrotik_firewall_nat" "logged_port_forward" {
  chain        = "dstnat"
  action       = "dst-nat"
  protocol     = "tcp"
  dst_port     = "22"
  in_interface = "ether1-wan"
  to_addresses = "192.168.1.50"
  log          = true
  log_prefix   = "SSH-FWD"
  comment      = "Logged SSH port forwarding"
}

## Example 16: Time-Based NAT
resource "mikrotik_firewall_nat" "business_hours_nat" {
  chain        = "srcnat"
  action       = "masquerade"
  src_address  = "192.168.1.0/24"
  out_interface = "ether1-wan"
  time         = "08:00-17:00,mon,tue,wed,thu,fri"
  comment      = "NAT only during business hours"
}

## Example 17: NAT with Rate Limiting
resource "mikrotik_firewall_nat" "limited_port_forward" {
  chain        = "dstnat"
  action       = "dst-nat"
  protocol     = "tcp"
  dst_port     = "80"
  in_interface = "ether1-wan"
  to_addresses = "192.168.1.100"
  limit        = "10/1m,20:packet"
  comment      = "Rate-limited port forwarding"
}

## Example 18: NAT for Specific Interface List
resource "mikrotik_firewall_nat" "wan_list_masquerade" {
  chain              = "srcnat"
  action             = "masquerade"
  out_interface_list = "WAN"
  comment            = "Masquerade for all WAN interfaces"
}

## Example 19: NAT with Packet Mark
resource "mikrotik_firewall_nat" "marked_traffic_nat" {
  chain        = "srcnat"
  action       = "src-nat"
  packet_mark  = "premium-traffic"
  to_addresses = "203.0.113.20"
  comment      = "NAT for marked premium traffic"
}

## Example 20: UDP Port Forwarding (Gaming/VoIP)
resource "mikrotik_firewall_nat" "game_server_forward" {
  chain        = "dstnat"
  action       = "dst-nat"
  protocol     = "udp"
  dst_port     = "27015-27030"
  in_interface = "ether1-wan"
  to_addresses = "192.168.1.200"
  comment      = "Game server UDP ports"
}

## Example 21: NAT with Address List
resource "mikrotik_firewall_nat" "whitelist_nat" {
  chain            = "srcnat"
  action           = "masquerade"
  src_address_list = "allowed-clients"
  out_interface    = "ether1-wan"
  comment          = "NAT only for whitelisted clients"
}

## Example 22: Disabled NAT Rule (for testing)
resource "mikrotik_firewall_nat" "test_forward" {
  chain        = "dstnat"
  action       = "dst-nat"
  protocol     = "tcp"
  dst_port     = "8080"
  in_interface = "ether1-wan"
  to_addresses = "192.168.1.150"
  disabled     = true
  comment      = "Test forward - disabled"
}

## Outputs
output "masquerade_rule_id" {
  value       = mikrotik_firewall_nat.masquerade_wan.id
  description = "ID of the masquerade NAT rule"
}

output "http_forward_stats" {
  value = {
    id      = mikrotik_firewall_nat.http_port_forward.id
    packets = mikrotik_firewall_nat.http_port_forward.packets
    bytes   = mikrotik_firewall_nat.http_port_forward.bytes
  }
  description = "HTTP port forward statistics"
}
