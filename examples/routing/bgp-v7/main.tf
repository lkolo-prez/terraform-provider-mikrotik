# BGP v7 Configuration Example
# RouterOS 7.14.3+

terraform {
  required_version = ">= 1.0"
  
  required_providers {
    mikrotik = {
      source  = "ddelnano/mikrotik"
      version = "~> 1.0"
    }
  }
}

provider "mikrotik" {
  host     = var.mikrotik_host
  username = var.mikrotik_username
  password = var.mikrotik_password
  tls      = var.mikrotik_tls
}

# ============================================================================
# BGP Template (Reusable Configuration)
# ============================================================================

resource "mikrotik_bgp_template" "default" {
  name    = "default-template"
  as      = var.local_as
  
  # Router ID (usually loopback or main IP)
  router_id = var.router_id
  
  # Address families
  address_families = "ip"  # ip, ipv6, vpnv4, vpnv6
  
  # Timers
  holdtime    = "180s"
  keepalive   = "60s"
  
  # Output filter (optional)
  output_filter = "bgp-out"
  
  # Input filter (optional)
  input_filter = "bgp-in"
  
  comment = "Default BGP template for all connections"
}

# ============================================================================
# BGP Connection to ISP
# ============================================================================

resource "mikrotik_bgp_connection" "isp" {
  name = "isp-primary"
  
  # Use template for common settings
  template = mikrotik_bgp_template.default.name
  
  # Remote peer configuration
  remote_address = var.isp_bgp_peer
  remote_as      = var.isp_as
  
  # Local configuration
  local_address = var.local_bgp_address
  local_as      = var.local_as
  
  # Connection options
  multihop         = false
  nexthop_choice   = "default"
  
  # Enable BFD for fast failure detection (optional)
  use_bfd = var.enable_bfd
  
  # MPLS and VPN support (optional)
  vpnv4 = false
  vpnv6 = false
  
  # Connection password (MD5 auth)
  tcp_md5_key = var.bgp_password
  
  # Listen on specific interface
  listen  = true
  connect = true
  
  disabled = false
  comment  = "Primary ISP BGP connection - AS${var.isp_as}"
}

# ============================================================================
# Secondary ISP Connection (Optional for redundancy)
# ============================================================================

resource "mikrotik_bgp_connection" "isp_backup" {
  count = var.enable_backup_isp ? 1 : 0
  
  name = "isp-backup"
  
  template = mikrotik_bgp_template.default.name
  
  remote_address = var.backup_isp_bgp_peer
  remote_as      = var.backup_isp_as
  
  local_address = var.local_bgp_address
  local_as      = var.local_as
  
  # Lower local preference for backup path
  input_filter = "bgp-in-backup"
  
  tcp_md5_key = var.bgp_backup_password
  
  disabled = false
  comment  = "Backup ISP BGP connection - AS${var.backup_isp_as}"
}

# ============================================================================
# Firewall Rules for BGP
# ============================================================================

resource "mikrotik_firewall_filter" "allow_bgp_in" {
  chain    = "input"
  protocol = "tcp"
  dst_port = "179"
  
  src_address_list = "bgp-neighbors"
  connection_state = "new"
  
  action  = "accept"
  comment = "Allow BGP from known neighbors"
  
  # Place before drop rules
  place_before = var.firewall_drop_rule_id
}

resource "mikrotik_firewall_filter" "allow_bgp_established" {
  chain = "input"
  protocol = "tcp"
  src_port = "179"
  
  connection_state = "established,related"
  
  action  = "accept"
  comment = "Allow established BGP connections"
  
  place_before = var.firewall_drop_rule_id
}

# ============================================================================
# Address List for BGP Neighbors
# ============================================================================

resource "mikrotik_ip_address_list" "bgp_neighbors" {
  list    = "bgp-neighbors"
  address = var.isp_bgp_peer
  comment = "Primary ISP BGP peer"
}

resource "mikrotik_ip_address_list" "bgp_neighbors_backup" {
  count = var.enable_backup_isp ? 1 : 0
  
  list    = "bgp-neighbors"
  address = var.backup_isp_bgp_peer
  comment = "Backup ISP BGP peer"
}

# ============================================================================
# BFD Configuration (Optional)
# ============================================================================

# Note: BFD resource not yet implemented in provider
# Configure manually or via script resource:
#
# /routing/bfd/configuration/add
#   multiplier=5
#   min-rx=200ms
#   min-tx=200ms

resource "mikrotik_script" "configure_bfd" {
  count = var.enable_bfd ? 1 : 0
  
  name   = "configure-bgp-bfd"
  source = <<-EOT
    /routing/bfd/configuration
    :if ([len [find]] = 0) do={
      add multiplier=5 min-rx=200ms min-tx=200ms
    }
  EOT
  
  comment = "Configure BFD for BGP"
}

resource "mikrotik_scheduler" "run_bfd_config" {
  count = var.enable_bfd ? 1 : 0
  
  name        = "init-bfd-config"
  on_event    = mikrotik_script.configure_bfd[0].name
  start_time  = "startup"
  interval    = 0
  
  comment = "Run BFD configuration at startup"
}

# ============================================================================
# Outputs
# ============================================================================

output "bgp_template_name" {
  description = "Name of the BGP template"
  value       = mikrotik_bgp_template.default.name
}

output "primary_connection" {
  description = "Primary ISP BGP connection details"
  value = {
    name       = mikrotik_bgp_connection.isp.name
    remote_as  = mikrotik_bgp_connection.isp.remote_as
    remote_ip  = mikrotik_bgp_connection.isp.remote_address
  }
}

output "backup_connection" {
  description = "Backup ISP BGP connection details (if enabled)"
  value = var.enable_backup_isp ? {
    name      = mikrotik_bgp_connection.isp_backup[0].name
    remote_as = mikrotik_bgp_connection.isp_backup[0].remote_as
    remote_ip = mikrotik_bgp_connection.isp_backup[0].remote_address
  } : null
}
