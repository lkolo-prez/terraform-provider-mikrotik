# BGP VPN/MPLS (VRF) Example
# Layer 3 VPN with MPLS and route distinguishers
# RouterOS 7.20+

terraform {
  required_providers {
    mikrotik = {
      source = "terraform-provider-mikrotik/mikrotik"
    }
  }
}

provider "mikrotik" {
  alias    = "pe1"
  host     = "192.168.88.1"
  username = "admin"
  password = "admin"
}

provider "mikrotik" {
  alias    = "pe2"
  host     = "192.168.88.2"
  username = "admin"
  password = "admin"
}

# ========== PE1 Configuration ==========

# Main BGP instance for MPLS backbone
resource "mikrotik_bgp_instance_v7" "pe1_main" {
  provider = mikrotik.pe1
  
  name      = "main"
  as        = 65300
  router_id = "10.0.3.1"
  
  # Main routing table
  routing_table = "main"
  
  redistribute_connected = true
  
  comment = "PE1 main BGP instance"
}

# VRF for customer A
resource "mikrotik_bgp_instance_v7" "pe1_customer_a" {
  provider = mikrotik.pe1
  
  name      = "customer-a"
  as        = 65300
  router_id = "10.0.3.1"
  
  # VRF settings
  vrf           = "customer-a-vrf"
  routing_table = "customer-a"
  
  redistribute_connected = true
  redistribute_static    = true
  
  comment = "PE1 BGP instance for Customer A VRF"
}

# VRF for customer B
resource "mikrotik_bgp_instance_v7" "pe1_customer_b" {
  provider = mikrotik.pe1
  
  name      = "customer-b"
  as        = 65300
  router_id = "10.0.3.1"
  
  vrf           = "customer-b-vrf"
  routing_table = "customer-b"
  
  redistribute_connected = true
  redistribute_static    = true
  
  comment = "PE1 BGP instance for Customer B VRF"
}

# MPLS backbone connection PE1 -> PE2
resource "mikrotik_bgp_connection" "pe1_to_pe2_mpls" {
  provider = mikrotik.pe1
  
  name           = "pe2-mpls"
  instance       = mikrotik_bgp_instance_v7.pe1_main.name
  remote_address = "10.0.3.2"
  remote_as      = 65300
  
  # Enable MPLS and VPN address families
  use_mpls         = true
  address_families = "ip,vpnv4,vpnv6"
  
  # Accept VPN routes
  input_accept_nlri = "bgp,vpn,bgp-mpls-vpn"
  
  multihop = true
  use_bfd  = true
  
  comment = "MPLS backbone to PE2"
}

# Customer A connection (CE1 -> PE1)
resource "mikrotik_bgp_connection" "pe1_customer_a_ce" {
  provider = mikrotik.pe1
  
  name           = "customer-a-ce1"
  instance       = mikrotik_bgp_instance_v7.pe1_customer_a.name
  remote_address = "172.16.1.1"
  remote_as      = 65001
  
  # VRF and route distinguisher
  vrf                = "customer-a-vrf"
  route_distinguisher = "65300:1001"
  
  # Standard eBGP settings
  address_families = "ip"
  
  comment = "Customer A CE1 connection"
}

# Customer B connection (CE2 -> PE1)
resource "mikrotik_bgp_connection" "pe1_customer_b_ce" {
  provider = mikrotik.pe1
  
  name           = "customer-b-ce2"
  instance       = mikrotik_bgp_instance_v7.pe1_customer_b.name
  remote_address = "172.16.2.1"
  remote_as      = 65002
  
  vrf                = "customer-b-vrf"
  route_distinguisher = "65300:2001"
  
  address_families = "ip"
  
  comment = "Customer B CE2 connection"
}

# ========== PE2 Configuration ==========

resource "mikrotik_bgp_instance_v7" "pe2_main" {
  provider = mikrotik.pe2
  
  name      = "main"
  as        = 65300
  router_id = "10.0.3.2"
  
  routing_table = "main"
  
  redistribute_connected = true
  
  comment = "PE2 main BGP instance"
}

resource "mikrotik_bgp_instance_v7" "pe2_customer_a" {
  provider = mikrotik.pe2
  
  name      = "customer-a"
  as        = 65300
  router_id = "10.0.3.2"
  
  vrf           = "customer-a-vrf"
  routing_table = "customer-a"
  
  redistribute_connected = true
  redistribute_static    = true
  
  comment = "PE2 BGP instance for Customer A VRF"
}

resource "mikrotik_bgp_instance_v7" "pe2_customer_b" {
  provider = mikrotik.pe2
  
  name      = "customer-b"
  as        = 65300
  router_id = "10.0.3.2"
  
  vrf           = "customer-b-vrf"
  routing_table = "customer-b"
  
  redistribute_connected = true
  redistribute_static    = true
  
  comment = "PE2 BGP instance for Customer B VRF"
}

# MPLS backbone PE2 -> PE1
resource "mikrotik_bgp_connection" "pe2_to_pe1_mpls" {
  provider = mikrotik.pe2
  
  name           = "pe1-mpls"
  instance       = mikrotik_bgp_instance_v7.pe2_main.name
  remote_address = "10.0.3.1"
  remote_as      = 65300
  
  use_mpls         = true
  address_families = "ip,vpnv4,vpnv6"
  
  input_accept_nlri = "bgp,vpn,bgp-mpls-vpn"
  
  multihop = true
  use_bfd  = true
  
  comment = "MPLS backbone to PE1"
}

# Customer A connection (CE3 -> PE2)
resource "mikrotik_bgp_connection" "pe2_customer_a_ce" {
  provider = mikrotik.pe2
  
  name           = "customer-a-ce3"
  instance       = mikrotik_bgp_instance_v7.pe2_customer_a.name
  remote_address = "172.16.1.2"
  remote_as      = 65001
  
  vrf                = "customer-a-vrf"
  route_distinguisher = "65300:1001"
  
  address_families = "ip"
  
  comment = "Customer A CE3 connection"
}

# Customer B connection (CE4 -> PE2)
resource "mikrotik_bgp_connection" "pe2_customer_b_ce" {
  provider = mikrotik.pe2
  
  name           = "customer-b-ce4"
  instance       = mikrotik_bgp_instance_v7.pe2_customer_b.name
  remote_address = "172.16.2.2"
  remote_as      = 65002
  
  vrf                = "customer-b-vrf"
  route_distinguisher = "65300:2001"
  
  address_families = "ip"
  
  comment = "Customer B CE4 connection"
}

# ========== Monitoring ==========

data "mikrotik_bgp_session" "pe1_sessions" {
  provider = mikrotik.pe1
  
  depends_on = [
    mikrotik_bgp_connection.pe1_to_pe2_mpls,
    mikrotik_bgp_connection.pe1_customer_a_ce,
    mikrotik_bgp_connection.pe1_customer_b_ce
  ]
}

data "mikrotik_bgp_session" "pe2_sessions" {
  provider = mikrotik.pe2
  
  depends_on = [
    mikrotik_bgp_connection.pe2_to_pe1_mpls,
    mikrotik_bgp_connection.pe2_customer_a_ce,
    mikrotik_bgp_connection.pe2_customer_b_ce
  ]
}

# ========== Outputs ==========

output "pe1_vpn_status" {
  value = {
    mpls_enabled     = "yes"
    customer_a_rd    = "65300:1001"
    customer_b_rd    = "65300:2001"
    total_sessions   = 3
    sessions_established = data.mikrotik_bgp_session.pe1_sessions.established
  }
  description = "PE1 VPN/MPLS status"
}

output "pe2_vpn_status" {
  value = {
    mpls_enabled     = "yes"
    customer_a_rd    = "65300:1001"
    customer_b_rd    = "65300:2001"
    total_sessions   = 3
    sessions_established = data.mikrotik_bgp_session.pe2_sessions.established
  }
  description = "PE2 VPN/MPLS status"
}

output "vpn_topology" {
  value = {
    description = "Layer 3 VPN with MPLS"
    customers   = 2
    sites_per_customer = 2
    route_distinguishers = ["65300:1001", "65300:2001"]
    isolation   = "Full VRF isolation between customers"
  }
  description = "VPN topology overview"
}
