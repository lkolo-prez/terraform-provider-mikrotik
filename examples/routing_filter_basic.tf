# Routing Filter - Basic BGP Filtering Example
#
# This example demonstrates basic BGP input filtering:
# - Deny default route
# - Deny bogon/private networks
# - Block small prefixes
# - Accept everything else

terraform {
  required_providers {
    mikrotik = {
      source = "lkolo-prez/mikrotik"
    }
  }
}

provider "mikrotik" {
  host     = var.mikrotik_host
  username = var.mikrotik_username
  password = var.mikrotik_password
  tls      = true
  insecure = true  # Only for testing with self-signed cert
}

# ============================================================================
# BGP Input Filter Chain
# ============================================================================

resource "mikrotik_routing_filter_chain" "bgp_in" {
  name     = "bgp-input-filter"
  dynamic  = true
  disabled = false
  comment  = "BGP input filtering - basic security"
}

# Deny default route
resource "mikrotik_routing_filter_rule" "deny_default_v4" {
  chain    = mikrotik_routing_filter_chain.bgp_in.name
  rule     = "if (dst == 0.0.0.0/0) { reject }"
  disabled = false
  comment  = "Block default route IPv4"
}

resource "mikrotik_routing_filter_rule" "deny_default_v6" {
  chain    = mikrotik_routing_filter_chain.bgp_in.name
  rule     = "if (dst == ::/0) { reject }"
  disabled = false
  comment  = "Block default route IPv6"
}

# Deny RFC1918 private networks
resource "mikrotik_routing_filter_rule" "deny_rfc1918" {
  chain    = mikrotik_routing_filter_chain.bgp_in.name
  rule     = <<-EOT
    if (dst in 192.168.0.0/16 || dst in 10.0.0.0/8 || dst in 172.16.0.0/12) {
      reject
    }
  EOT
  disabled = false
  comment  = "Block RFC1918 private networks"
}

# Deny other bogon networks
resource "mikrotik_routing_filter_rule" "deny_bogons" {
  chain    = mikrotik_routing_filter_chain.bgp_in.name
  rule     = <<-EOT
    if (dst in 127.0.0.0/8 || dst in 169.254.0.0/16 || dst in 224.0.0.0/4) {
      reject
    }
  EOT
  disabled = false
  comment  = "Block loopback, link-local, multicast"
}

# Block small prefixes (DDoS mitigation)
resource "mikrotik_routing_filter_rule" "block_small_prefixes" {
  chain    = mikrotik_routing_filter_chain.bgp_in.name
  rule     = "if (dst-len > 24) { reject }"
  disabled = false
  comment  = "Block prefixes longer than /24"
}

# Accept everything else
resource "mikrotik_routing_filter_rule" "accept_all" {
  chain    = mikrotik_routing_filter_chain.bgp_in.name
  rule     = "accept"
  disabled = false
  comment  = "Accept all other routes"
}

# ============================================================================
# BGP Output Filter Chain (what we advertise)
# ============================================================================

resource "mikrotik_routing_filter_chain" "bgp_out" {
  name     = "bgp-output-filter"
  dynamic  = true
  disabled = false
  comment  = "BGP output filtering - advertise only our networks"
}

# Only advertise routes from our AS (no transit)
resource "mikrotik_routing_filter_rule" "out_our_as_only" {
  chain    = mikrotik_routing_filter_chain.bgp_out.name
  rule     = "if (bgp-as-path ~ \"^$\") { accept }"
  disabled = false
  comment  = "Only advertise routes originating from our AS"
}

# Reject everything else (safety)
resource "mikrotik_routing_filter_rule" "out_reject_all" {
  chain    = mikrotik_routing_filter_chain.bgp_out.name
  rule     = "reject"
  disabled = false
  comment  = "Don't advertise anything else"
}

# ============================================================================
# BGP Connection using filters
# ============================================================================

resource "mikrotik_bgp_connection" "isp" {
  name           = "isp-primary"
  template       = "default"
  remote_address = var.bgp_peer_address
  remote_as      = var.bgp_peer_as
  local_as       = var.bgp_local_as
  
  # Apply filters
  input_filter  = mikrotik_routing_filter_chain.bgp_in.name
  output_filter = mikrotik_routing_filter_chain.bgp_out.name
  
  disabled = false
  comment  = "Primary ISP BGP connection with filtering"
}

# ============================================================================
# Outputs
# ============================================================================

output "input_filter_chain" {
  description = "BGP input filter chain name"
  value       = mikrotik_routing_filter_chain.bgp_in.name
}

output "output_filter_chain" {
  description = "BGP output filter chain name"
  value       = mikrotik_routing_filter_chain.bgp_out.name
}
