# Route Reflector Example
# Hub-and-spoke iBGP topology using route reflector
# Reduces connections from N*(N-1)/2 to N-1
# RouterOS 7.20+

terraform {
  required_providers {
    mikrotik = {
      source = "terraform-provider-mikrotik/mikrotik"
    }
  }
}

# ========== Route Reflector (Hub) ==========

provider "mikrotik" {
  alias    = "rr"
  host     = "192.168.88.1"
  username = "admin"
  password = "admin"
}

resource "mikrotik_bgp_instance_v7" "route_reflector" {
  provider = mikrotik.rr
  
  name      = "rr-main"
  as        = 65200
  router_id = "10.0.2.1"
  
  # Enable route reflection
  client_to_client_reflection = true
  
  # Optional: Set cluster ID for multiple RRs
  cluster_id = "10.0.2.1"
  
  redistribute_connected = true
  redistribute_static    = true
  redistribute_ospf      = true
  
  comment = "Route Reflector instance"
}

# Template for route reflector clients
resource "mikrotik_bgp_template" "rr_client_template" {
  provider = mikrotik.rr
  
  name = "rr-clients"
  as   = 65200
  
  # Enable route reflection
  route_reflect = true
  
  # Accept all NLRI from clients
  input_accept_nlri = "bgp,vpn,rtc,flow,srte,rip,ospf,connected,static"
  
  # Full capabilities
  capabilities = "mp,refresh,as4,route-refresh"
  
  address_families = "ip,ipv6"
  
  hold_time      = "3m"
  keepalive_time = "1m"
  
  comment = "Template for route reflector clients"
}

# ========== Client Routers (Spokes) ==========

# Client 1
provider "mikrotik" {
  alias    = "client1"
  host     = "192.168.88.11"
  username = "admin"
  password = "admin"
}

resource "mikrotik_bgp_instance_v7" "client1" {
  provider = mikrotik.client1
  
  name      = "main"
  as        = 65200
  router_id = "10.0.2.11"
  
  # Clients don't do route reflection
  client_to_client_reflection = false
  
  redistribute_connected = true
  
  comment = "Client 1 BGP instance"
}

resource "mikrotik_bgp_connection" "client1_to_rr" {
  provider = mikrotik.client1
  
  name           = "to-route-reflector"
  instance       = mikrotik_bgp_instance_v7.client1.name
  remote_address = "10.0.2.1"
  remote_as      = 65200
  
  # Use multihop for non-directly connected RR
  multihop = true
  use_bfd  = true
  
  address_families = "ip,ipv6"
  
  comment = "Connection to route reflector"
}

# Client 2
provider "mikrotik" {
  alias    = "client2"
  host     = "192.168.88.12"
  username = "admin"
  password = "admin"
}

resource "mikrotik_bgp_instance_v7" "client2" {
  provider = mikrotik.client2
  
  name      = "main"
  as        = 65200
  router_id = "10.0.2.12"
  
  client_to_client_reflection = false
  
  redistribute_connected = true
  
  comment = "Client 2 BGP instance"
}

resource "mikrotik_bgp_connection" "client2_to_rr" {
  provider = mikrotik.client2
  
  name           = "to-route-reflector"
  instance       = mikrotik_bgp_instance_v7.client2.name
  remote_address = "10.0.2.1"
  remote_as      = 65200
  
  multihop = true
  use_bfd  = true
  
  address_families = "ip,ipv6"
  
  comment = "Connection to route reflector"
}

# Client 3
provider "mikrotik" {
  alias    = "client3"
  host     = "192.168.88.13"
  username = "admin"
  password = "admin"
}

resource "mikrotik_bgp_instance_v7" "client3" {
  provider = mikrotik.client3
  
  name      = "main"
  as        = 65200
  router_id = "10.0.2.13"
  
  client_to_client_reflection = false
  
  redistribute_connected = true
  
  comment = "Client 3 BGP instance"
}

resource "mikrotik_bgp_connection" "client3_to_rr" {
  provider = mikrotik.client3
  
  name           = "to-route-reflector"
  instance       = mikrotik_bgp_instance_v7.client3.name
  remote_address = "10.0.2.1"
  remote_as      = 65200
  
  multihop = true
  use_bfd  = true
  
  address_families = "ip,ipv6"
  
  comment = "Connection to route reflector"
}

# ========== Route Reflector Connections ==========

# RR -> Client 1
resource "mikrotik_bgp_connection" "rr_to_client1" {
  provider = mikrotik.rr
  
  name           = "client-1"
  instance       = mikrotik_bgp_instance_v7.route_reflector.name
  remote_address = "10.0.2.11"
  remote_as      = 65200
  
  templates = [mikrotik_bgp_template.rr_client_template.name]
  
  multihop = true
  use_bfd  = true
  
  # This is a route reflector client
  local_role = "route-reflector"
  
  comment = "RR connection to Client 1"
}

# RR -> Client 2
resource "mikrotik_bgp_connection" "rr_to_client2" {
  provider = mikrotik.rr
  
  name           = "client-2"
  instance       = mikrotik_bgp_instance_v7.route_reflector.name
  remote_address = "10.0.2.12"
  remote_as      = 65200
  
  templates = [mikrotik_bgp_template.rr_client_template.name]
  
  multihop   = true
  use_bfd    = true
  local_role = "route-reflector"
  
  comment = "RR connection to Client 2"
}

# RR -> Client 3
resource "mikrotik_bgp_connection" "rr_to_client3" {
  provider = mikrotik.rr
  
  name           = "client-3"
  instance       = mikrotik_bgp_instance_v7.route_reflector.name
  remote_address = "10.0.2.13"
  remote_as      = 65200
  
  templates = [mikrotik_bgp_template.rr_client_template.name]
  
  multihop   = true
  use_bfd    = true
  local_role = "route-reflector"
  
  comment = "RR connection to Client 3"
}

# ========== Monitoring ==========

data "mikrotik_bgp_session" "rr_sessions" {
  provider = mikrotik.rr
  
  depends_on = [
    mikrotik_bgp_connection.rr_to_client1,
    mikrotik_bgp_connection.rr_to_client2,
    mikrotik_bgp_connection.rr_to_client3
  ]
}

# ========== Outputs ==========

output "route_reflector_status" {
  value = {
    cluster_id          = mikrotik_bgp_instance_v7.route_reflector.cluster_id
    total_clients       = 3
    sessions_established = data.mikrotik_bgp_session.rr_sessions.established
    state               = data.mikrotik_bgp_session.rr_sessions.state
  }
  description = "Route reflector status"
}

output "topology_efficiency" {
  value = {
    full_mesh_connections   = "N*(N-1)/2 = 4*(3)/2 = 6"
    route_reflector_connections = "N-1 = 4-1 = 3"
    reduction               = "50%"
  }
  description = "Connection reduction using route reflector"
}
