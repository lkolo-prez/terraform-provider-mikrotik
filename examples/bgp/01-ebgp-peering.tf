# eBGP Peering Example
# Two routers establishing external BGP connection
# RouterOS 7.20+

terraform {
  required_providers {
    mikrotik = {
      source = "terraform-provider-mikrotik/mikrotik"
    }
  }
}

# Provider for Router 1 (AS 65001)
provider "mikrotik" {
  alias    = "router1"
  host     = "192.168.88.1"
  username = "admin"
  password = "admin"
  tls      = true
  insecure = true
}

# Provider for Router 2 (AS 65002)
provider "mikrotik" {
  alias    = "router2"
  host     = "192.168.88.2"
  username = "admin"
  password = "admin"
  tls      = true
  insecure = true
}

# ========== Router 1 Configuration (AS 65001) ==========

resource "mikrotik_bgp_instance_v7" "router1" {
  provider = mikrotik.router1
  
  name      = "default"
  as        = 65001
  router_id = "10.0.0.1"
  
  # Redistribute local routes
  redistribute_connected = true
  redistribute_static    = true
  
  comment = "Main BGP instance for AS 65001"
}

resource "mikrotik_bgp_connection" "router1_to_router2" {
  provider = mikrotik.router1
  
  name           = "to-as65002"
  instance       = mikrotik_bgp_instance_v7.router1.name
  remote_address = "192.168.1.2"
  remote_as      = 65002
  
  # Connection settings
  listen             = false
  connect            = true
  address_families   = "ip"
  
  # Timers (default values shown)
  hold_time       = "3m"
  keepalive_time  = "1m"
  
  # Security
  tcp_md5_key = "secret-key-123"
  
  comment = "eBGP peering to Router 2 (AS 65002)"
  
  depends_on = [mikrotik_bgp_instance_v7.router1]
}

# ========== Router 2 Configuration (AS 65002) ==========

resource "mikrotik_bgp_instance_v7" "router2" {
  provider = mikrotik.router2
  
  name      = "default"
  as        = 65002
  router_id = "10.0.0.2"
  
  redistribute_connected = true
  redistribute_static    = true
  
  comment = "Main BGP instance for AS 65002"
}

resource "mikrotik_bgp_connection" "router2_to_router1" {
  provider = mikrotik.router2
  
  name           = "to-as65001"
  instance       = mikrotik_bgp_instance_v7.router2.name
  remote_address = "192.168.1.1"
  remote_as      = 65001
  
  listen           = false
  connect          = true
  address_families = "ip"
  
  hold_time      = "3m"
  keepalive_time = "1m"
  
  tcp_md5_key = "secret-key-123"
  
  comment = "eBGP peering to Router 1 (AS 65001)"
  
  depends_on = [mikrotik_bgp_instance_v7.router2]
}

# ========== Monitor BGP Sessions ==========

data "mikrotik_bgp_session" "router1_sessions" {
  provider = mikrotik.router1
  
  depends_on = [
    mikrotik_bgp_connection.router1_to_router2
  ]
}

data "mikrotik_bgp_session" "router2_sessions" {
  provider = mikrotik.router2
  
  depends_on = [
    mikrotik_bgp_connection.router2_to_router1
  ]
}

# ========== Outputs ==========

output "router1_session_state" {
  value = {
    established = data.mikrotik_bgp_session.router1_sessions.established
    state       = data.mikrotik_bgp_session.router1_sessions.state
    uptime      = data.mikrotik_bgp_session.router1_sessions.uptime
  }
  description = "Router 1 BGP session status"
}

output "router2_session_state" {
  value = {
    established = data.mikrotik_bgp_session.router2_sessions.established
    state       = data.mikrotik_bgp_session.router2_sessions.state
    uptime      = data.mikrotik_bgp_session.router2_sessions.uptime
  }
  description = "Router 2 BGP session status"
}
