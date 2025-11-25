# iBGP Full Mesh Example
# Three routers in same AS with full mesh connectivity
# RouterOS 7.20+

terraform {
  required_providers {
    mikrotik = {
      source = "terraform-provider-mikrotik/mikrotik"
    }
  }
}

# Provider configurations
provider "mikrotik" {
  alias    = "r1"
  host     = "192.168.88.1"
  username = "admin"
  password = "admin"
}

provider "mikrotik" {
  alias    = "r2"
  host     = "192.168.88.2"
  username = "admin"
  password = "admin"
}

provider "mikrotik" {
  alias    = "r3"
  host     = "192.168.88.3"
  username = "admin"
  password = "admin"
}

# ========== BGP Template for iBGP Peers ==========
# Shared configuration for all iBGP connections

resource "mikrotik_bgp_template" "ibgp_template" {
  provider = mikrotik.r1
  
  name = "ibgp-peers"
  as   = 65100
  
  # Address families
  address_families = "ip,ipv6"
  
  # Capabilities
  capabilities = "mp,refresh,as4"
  
  # iBGP specific settings
  nexthop_choice = "force-self"
  
  # Input filtering
  input_accept_nlri       = "bgp,vpn,rtc,flow,srte,rip,ospf,connected,static,bgp-mpls-vpn"
  input_accept_communities = "all"
  
  # Output filtering
  output_default_originate = "never"
  
  # Timers
  hold_time      = "3m"
  keepalive_time = "1m"
  
  comment = "Template for iBGP full mesh connections"
}

# ========== Router 1 (10.0.1.1) ==========

resource "mikrotik_bgp_instance_v7" "r1_instance" {
  provider = mikrotik.r1
  
  name      = "main"
  as        = 65100
  router_id = "10.0.1.1"
  
  # Enable route reflection for scalability
  client_to_client_reflection = true
  
  redistribute_connected = true
  redistribute_static    = true
  redistribute_ospf      = true
  
  comment = "Router 1 BGP instance"
}

# R1 -> R2 connection
resource "mikrotik_bgp_connection" "r1_to_r2" {
  provider = mikrotik.r1
  
  name           = "ibgp-r2"
  instance       = mikrotik_bgp_instance_v7.r1_instance.name
  remote_address = "10.0.1.2"
  remote_as      = 65100
  
  templates = [mikrotik_bgp_template.ibgp_template.name]
  
  # iBGP requires full connectivity
  multihop     = true
  use_bfd      = true  # Enable BFD for fast failure detection
  
  comment = "iBGP to Router 2"
}

# R1 -> R3 connection
resource "mikrotik_bgp_connection" "r1_to_r3" {
  provider = mikrotik.r1
  
  name           = "ibgp-r3"
  instance       = mikrotik_bgp_instance_v7.r1_instance.name
  remote_address = "10.0.1.3"
  remote_as      = 65100
  
  templates = [mikrotik_bgp_template.ibgp_template.name]
  
  multihop = true
  use_bfd  = true
  
  comment = "iBGP to Router 3"
}

# ========== Router 2 (10.0.1.2) ==========

resource "mikrotik_bgp_instance_v7" "r2_instance" {
  provider = mikrotik.r2
  
  name      = "main"
  as        = 65100
  router_id = "10.0.1.2"
  
  client_to_client_reflection = true
  
  redistribute_connected = true
  redistribute_static    = true
  redistribute_ospf      = true
  
  comment = "Router 2 BGP instance"
}

# R2 -> R1 connection
resource "mikrotik_bgp_connection" "r2_to_r1" {
  provider = mikrotik.r2
  
  name           = "ibgp-r1"
  instance       = mikrotik_bgp_instance_v7.r2_instance.name
  remote_address = "10.0.1.1"
  remote_as      = 65100
  
  templates = [mikrotik_bgp_template.ibgp_template.name]
  
  multihop = true
  use_bfd  = true
  
  comment = "iBGP to Router 1"
}

# R2 -> R3 connection
resource "mikrotik_bgp_connection" "r2_to_r3" {
  provider = mikrotik.r2
  
  name           = "ibgp-r3"
  instance       = mikrotik_bgp_instance_v7.r2_instance.name
  remote_address = "10.0.1.3"
  remote_as      = 65100
  
  templates = [mikrotik_bgp_template.ibgp_template.name]
  
  multihop = true
  use_bfd  = true
  
  comment = "iBGP to Router 3"
}

# ========== Router 3 (10.0.1.3) ==========

resource "mikrotik_bgp_instance_v7" "r3_instance" {
  provider = mikrotik.r3
  
  name      = "main"
  as        = 65100
  router_id = "10.0.1.3"
  
  client_to_client_reflection = true
  
  redistribute_connected = true
  redistribute_static    = true
  redistribute_ospf      = true
  
  comment = "Router 3 BGP instance"
}

# R3 -> R1 connection
resource "mikrotik_bgp_connection" "r3_to_r1" {
  provider = mikrotik.r3
  
  name           = "ibgp-r1"
  instance       = mikrotik_bgp_instance_v7.r3_instance.name
  remote_address = "10.0.1.1"
  remote_as      = 65100
  
  templates = [mikrotik_bgp_template.ibgp_template.name]
  
  multihop = true
  use_bfd  = true
  
  comment = "iBGP to Router 1"
}

# R3 -> R2 connection
resource "mikrotik_bgp_connection" "r3_to_r2" {
  provider = mikrotik.r3
  
  name           = "ibgp-r2"
  instance       = mikrotik_bgp_instance_v7.r3_instance.name
  remote_address = "10.0.1.2"
  remote_as      = 65100
  
  templates = [mikrotik_bgp_template.ibgp_template.name]
  
  multihop = true
  use_bfd  = true
  
  comment = "iBGP to Router 2"
}

# ========== Monitoring ==========

data "mikrotik_bgp_session" "r1_sessions" {
  provider   = mikrotik.r1
  depends_on = [mikrotik_bgp_connection.r1_to_r2, mikrotik_bgp_connection.r1_to_r3]
}

data "mikrotik_bgp_session" "r2_sessions" {
  provider   = mikrotik.r2
  depends_on = [mikrotik_bgp_connection.r2_to_r1, mikrotik_bgp_connection.r2_to_r3]
}

data "mikrotik_bgp_session" "r3_sessions" {
  provider   = mikrotik.r3
  depends_on = [mikrotik_bgp_connection.r3_to_r1, mikrotik_bgp_connection.r3_to_r2]
}

# ========== Outputs ==========

output "full_mesh_status" {
  value = {
    router1 = {
      sessions    = 2
      established = data.mikrotik_bgp_session.r1_sessions.established
    }
    router2 = {
      sessions    = 2
      established = data.mikrotik_bgp_session.r2_sessions.established
    }
    router3 = {
      sessions    = 2
      established = data.mikrotik_bgp_session.r3_sessions.established
    }
  }
  description = "iBGP full mesh status"
}
