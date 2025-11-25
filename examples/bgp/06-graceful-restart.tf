# BGP Graceful Restart Example
# High availability with graceful restart and BFD
# RouterOS 7.20+

terraform {
  required_providers {
    mikrotik = {
      source = "terraform-provider-mikrotik/mikrotik"
    }
  }
}

provider "mikrotik" {
  alias    = "core1"
  host     = "192.168.88.1"
  username = "admin"
  password = "admin"
}

provider "mikrotik" {
  alias    = "core2"
  host     = "192.168.88.2"
  username = "admin"
  password = "admin"
}

# ========== Core Router 1 ==========

resource "mikrotik_bgp_instance_v7" "core1" {
  provider = mikrotik.core1
  
  name      = "main"
  as        = 65600
  router_id = "10.0.6.1"
  
  client_to_client_reflection = true
  
  redistribute_connected = true
  redistribute_static    = true
  redistribute_ospf      = true
  
  comment = "Core 1 BGP with graceful restart"
}

# Template with graceful restart enabled
resource "mikrotik_bgp_template" "ha_peer" {
  provider = mikrotik.core1
  
  name = "ha-peers"
  as   = 65600
  
  # Address families
  address_families = "ip,ipv6"
  
  # Enable graceful restart capability
  capabilities = "mp,refresh,as4,graceful-restart"
  
  # Graceful restart settings
  # Time to wait for peer to restart (seconds)
  graceful_restart_time = 120
  
  # Preserve forwarding state during restart
  graceful_restart_stale_time = 300
  
  # Accept all NLRI types
  input_accept_nlri = "bgp,vpn,rtc,flow,srte,rip,ospf,connected,static"
  
  # Timers - reduced for faster convergence
  hold_time      = "90s"
  keepalive_time = "30s"
  
  comment = "HA template with graceful restart"
}

resource "mikrotik_bgp_connection" "core1_to_core2" {
  provider = mikrotik.core1
  
  name           = "to-core2"
  instance       = mikrotik_bgp_instance_v7.core1.name
  remote_address = "10.0.6.2"
  remote_as      = 65600
  
  # Apply graceful restart template
  templates = [mikrotik_bgp_template.ha_peer.name]
  
  # Connection settings
  connect          = true
  address_families = "ip,ipv6"
  
  # Enable BFD for fast failure detection (sub-second)
  use_bfd = true
  
  # Multihop for non-directly connected
  multihop = true
  
  comment = "HA connection to Core 2 with GR and BFD"
  
  depends_on = [mikrotik_bgp_template.ha_peer]
}

# ========== Core Router 2 ==========

resource "mikrotik_bgp_instance_v7" "core2" {
  provider = mikrotik.core2
  
  name      = "main"
  as        = 65600
  router_id = "10.0.6.2"
  
  client_to_client_reflection = true
  
  redistribute_connected = true
  redistribute_static    = true
  redistribute_ospf      = true
  
  comment = "Core 2 BGP with graceful restart"
}

resource "mikrotik_bgp_template" "ha_peer_core2" {
  provider = mikrotik.core2
  
  name = "ha-peers"
  as   = 65600
  
  address_families = "ip,ipv6"
  
  # Graceful restart capability
  capabilities = "mp,refresh,as4,graceful-restart"
  
  # Graceful restart timers
  graceful_restart_time       = 120
  graceful_restart_stale_time = 300
  
  input_accept_nlri = "bgp,vpn,rtc,flow,srte,rip,ospf,connected,static"
  
  hold_time      = "90s"
  keepalive_time = "30s"
  
  comment = "HA template with graceful restart"
}

resource "mikrotik_bgp_connection" "core2_to_core1" {
  provider = mikrotik.core2
  
  name           = "to-core1"
  instance       = mikrotik_bgp_instance_v7.core2.name
  remote_address = "10.0.6.1"
  remote_as      = 65600
  
  templates = [mikrotik_bgp_template.ha_peer_core2.name]
  
  connect          = true
  address_families = "ip,ipv6"
  
  # BFD for fast failure detection
  use_bfd = true
  
  multihop = true
  
  comment = "HA connection to Core 1 with GR and BFD"
  
  depends_on = [mikrotik_bgp_template.ha_peer_core2]
}

# ========== Monitoring ==========

data "mikrotik_bgp_session" "core1_sessions" {
  provider = mikrotik.core1
  
  depends_on = [mikrotik_bgp_connection.core1_to_core2]
}

data "mikrotik_bgp_session" "core2_sessions" {
  provider = mikrotik.core2
  
  depends_on = [mikrotik_bgp_connection.core2_to_core1]
}

# ========== Outputs ==========

output "core1_ha_status" {
  value = {
    session = {
      established = data.mikrotik_bgp_session.core1_sessions.established
      state       = data.mikrotik_bgp_session.core1_sessions.state
      uptime      = data.mikrotik_bgp_session.core1_sessions.uptime
    }
    capabilities = {
      graceful_restart = "enabled"
      bfd              = "enabled"
      restart_time     = "120s"
      stale_time       = "300s"
    }
  }
  description = "Core 1 HA status"
}

output "core2_ha_status" {
  value = {
    session = {
      established = data.mikrotik_bgp_session.core2_sessions.established
      state       = data.mikrotik_bgp_session.core2_sessions.state
      uptime      = data.mikrotik_bgp_session.core2_sessions.uptime
    }
    capabilities = {
      graceful_restart = "enabled"
      bfd              = "enabled"
      restart_time     = "120s"
      stale_time       = "300s"
    }
  }
  description = "Core 2 HA status"
}

output "ha_benefits" {
  value = {
    graceful_restart = {
      description = "Preserves forwarding state during BGP restart"
      restart_time = "120 seconds to restart BGP process"
      stale_time   = "300 seconds to keep routes while waiting"
      benefit      = "Zero packet loss during planned maintenance"
    }
    bfd = {
      description  = "Bidirectional Forwarding Detection"
      detection    = "Sub-second failure detection"
      benefit      = "Fast convergence on link failures"
    }
    combined = {
      planned_maintenance = "Graceful restart handles BGP process restart"
      hardware_failure    = "BFD detects and reacts to link failures"
      result              = "High availability with minimal downtime"
    }
  }
  description = "High availability features explanation"
}
