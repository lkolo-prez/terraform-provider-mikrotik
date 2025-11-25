# Example: Multi-VRF BGP Setup with Routing Tables

terraform {
  required_providers {
    mikrotik = {
      source = "terraform-provider-mikrotik/mikrotik"
    }
  }
}

provider "mikrotik" {
  host     = var.mikrotik_host
  username = var.mikrotik_username
  password = var.mikrotik_password
  tls      = true
}

# Create routing tables (VRFs) for different customers
resource "mikrotik_routing_table" "customer_a" {
  name    = "customer_a"
  fib     = "main"
  comment = "Customer A VRF - Isolated routing instance"
}

resource "mikrotik_routing_table" "customer_b" {
  name    = "customer_b"
  fib     = "main"
  comment = "Customer B VRF - Isolated routing instance"
}

resource "mikrotik_routing_table" "management" {
  name    = "management"
  fib     = "main"
  comment = "Management VRF - Separate control plane"
}

# BGP instances using VRFs
resource "mikrotik_bgp_instance_v7" "customer_a_bgp" {
  name          = "customer_a_bgp"
  as            = 65001
  router_id     = "10.0.1.1"
  routing_table = mikrotik_routing_table.customer_a.name
  vrf           = mikrotik_routing_table.customer_a.name
  
  redistribute_connected = true
  redistribute_static    = true
  
  comment = "BGP instance for Customer A"
}

resource "mikrotik_bgp_instance_v7" "customer_b_bgp" {
  name          = "customer_b_bgp"
  as            = 65002
  router_id     = "10.0.2.1"
  routing_table = mikrotik_routing_table.customer_b.name
  vrf           = mikrotik_routing_table.customer_b.name
  
  redistribute_connected = true
  
  comment = "BGP instance for Customer B"
}

# BGP connections per VRF
resource "mikrotik_bgp_connection" "customer_a_peer" {
  name           = "customer_a_peer"
  remote_address = "10.0.1.2"
  remote_as      = 65003
  instance       = mikrotik_bgp_instance_v7.customer_a_bgp.name
  
  # Connection operates in customer_a VRF
  routing_table = mikrotik_routing_table.customer_a.name
}

resource "mikrotik_bgp_connection" "customer_b_peer" {
  name           = "customer_b_peer"
  remote_address = "10.0.2.2"
  remote_as      = 65004
  instance       = mikrotik_bgp_instance_v7.customer_b_bgp.name
  
  # Connection operates in customer_b VRF
  routing_table = mikrotik_routing_table.customer_b.name
}

# Output routing table information
output "customer_a_vrf_id" {
  description = "Customer A VRF/Routing Table ID"
  value       = mikrotik_routing_table.customer_a.id
}

output "customer_b_vrf_id" {
  description = "Customer B VRF/Routing Table ID"
  value       = mikrotik_routing_table.customer_b.id
}

output "management_vrf_id" {
  description = "Management VRF/Routing Table ID"
  value       = mikrotik_routing_table.management.id
}
