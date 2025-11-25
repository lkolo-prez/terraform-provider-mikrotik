# OSPF Basic Configuration Example
# Single-area OSPF backbone with simple configuration

terraform {
  required_providers {
    mikrotik = {
      source = "ddelnano/mikrotik"
    }
  }
}

provider "mikrotik" {
  host     = var.mikrotik_host
  username = var.mikrotik_username
  password = var.mikrotik_password
  tls      = true
  insecure = true
}

variable "mikrotik_host" {
  type = string
}

variable "mikrotik_username" {
  type    = string
  default = "admin"
}

variable "mikrotik_password" {
  type      = string
  sensitive = true
}

# OSPF Instance (OSPFv2 for IPv4)
resource "mikrotik_ospf_instance_v7" "default" {
  name      = "default_v2"
  version   = "2"
  router_id = "1.1.1.1"
  comment   = "Basic OSPF instance"
}

# Backbone Area (required)
resource "mikrotik_ospf_area_v7" "backbone" {
  name     = "backbone"
  instance = mikrotik_ospf_instance_v7.default.name
  area_id  = "0.0.0.0"
  comment  = "OSPF backbone area"
}

# Enable OSPF on LAN networks
resource "mikrotik_ospf_interface_template_v7" "lan" {
  area = mikrotik_ospf_area_v7.backbone.name
  networks = [
    "192.168.1.0/24",
    "192.168.2.0/24"
  ]
  cost    = 10
  comment = "LAN segments"
}

# Enable OSPF on WAN link (point-to-point)
resource "mikrotik_ospf_interface_template_v7" "wan" {
  area    = mikrotik_ospf_area_v7.backbone.name
  networks = ["10.0.0.0/30"]
  type    = "ptp"
  cost    = 100
  comment = "WAN uplink"
}

# Advertise internal network without forming adjacencies
resource "mikrotik_ospf_interface_template_v7" "internal" {
  area     = mikrotik_ospf_area_v7.backbone.name
  networks = ["10.10.0.0/16"]
  passive  = true
  comment  = "Internal network - passive"
}

output "ospf_instance_id" {
  value       = mikrotik_ospf_instance_v7.default.id
  description = "OSPF instance ID"
}

output "ospf_router_id" {
  value       = mikrotik_ospf_instance_v7.default.router_id
  description = "OSPF router ID"
}
