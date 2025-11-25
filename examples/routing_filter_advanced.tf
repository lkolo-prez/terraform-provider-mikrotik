# Routing Filter - Advanced BGP with Communities
#
# This example demonstrates advanced BGP filtering with:
# - Community-based route preference
# - Multi-ISP setup with different priorities
# - Traffic engineering with local-preference
# - Customer-specific filtering

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
  insecure = true
}

# ============================================================================
# Primary ISP Filter Chain (High Priority)
# ============================================================================

resource "mikrotik_routing_filter_chain" "isp_primary_in" {
  name     = "isp-primary-in"
  dynamic  = true
  comment  = "Primary ISP input - highest preference"
}

resource "mikrotik_routing_filter_rule" "primary_deny_default" {
  chain   = mikrotik_routing_filter_chain.isp_primary_in.name
  rule    = "if (dst == 0.0.0.0/0 || dst == ::/0) { reject }"
  comment = "Block default routes"
}

resource "mikrotik_routing_filter_rule" "primary_deny_bogons" {
  chain = mikrotik_routing_filter_chain.isp_primary_in.name
  rule  = <<-EOT
    if (dst in 192.168.0.0/16 || dst in 10.0.0.0/8 || dst in 172.16.0.0/12) {
      reject
    }
  EOT
  comment = "Block private networks"
}

resource "mikrotik_routing_filter_rule" "primary_high_pref" {
  chain   = mikrotik_routing_filter_chain.isp_primary_in.name
  rule    = "set bgp-local-pref 200; accept"
  comment = "High preference for primary ISP"
}

# ============================================================================
# Backup ISP Filter Chain (Lower Priority)
# ============================================================================

resource "mikrotik_routing_filter_chain" "isp_backup_in" {
  name     = "isp-backup-in"
  dynamic  = true
  comment  = "Backup ISP input - lower preference"
}

resource "mikrotik_routing_filter_rule" "backup_deny_default" {
  chain   = mikrotik_routing_filter_chain.isp_backup_in.name
  rule    = "if (dst == 0.0.0.0/0 || dst == ::/0) { reject }"
  comment = "Block default routes"
}

resource "mikrotik_routing_filter_rule" "backup_deny_bogons" {
  chain = mikrotik_routing_filter_chain.isp_backup_in.name
  rule  = <<-EOT
    if (dst in 192.168.0.0/16 || dst in 10.0.0.0/8 || dst in 172.16.0.0/12) {
      reject
    }
  EOT
  comment = "Block private networks"
}

resource "mikrotik_routing_filter_rule" "backup_low_pref" {
  chain   = mikrotik_routing_filter_chain.isp_backup_in.name
  rule    = "set bgp-local-pref 100; accept"
  comment = "Lower preference for backup ISP"
}

# ============================================================================
# Customer Routes Filter (Community-Based)
# ============================================================================

resource "mikrotik_routing_filter_chain" "customer_routes" {
  name     = "customer-routes-in"
  dynamic  = true
  comment  = "Customer route filtering with community tags"
}

# Deny default from customers
resource "mikrotik_routing_filter_rule" "customer_deny_default" {
  chain   = mikrotik_routing_filter_chain.customer_routes.name
  rule    = "if (dst == 0.0.0.0/0) { reject }"
  comment = "Customers cannot send default route"
}

# High priority customer routes (community 65001:100)
resource "mikrotik_routing_filter_rule" "customer_high_priority" {
  chain = mikrotik_routing_filter_chain.customer_routes.name
  rule  = <<-EOT
    if (dst in 10.0.0.0/8 && bgp-communities includes 65001:100) {
      set bgp-local-pref 250;
      set bgp-weight 200;
      accept
    }
  EOT
  comment = "High priority customer routes"
}

# Medium priority customer routes (community 65001:150)
resource "mikrotik_routing_filter_rule" "customer_medium_priority" {
  chain = mikrotik_routing_filter_chain.customer_routes.name
  rule  = <<-EOT
    if (dst in 10.0.0.0/8 && bgp-communities includes 65001:150) {
      set bgp-local-pref 150;
      accept
    }
  EOT
  comment = "Medium priority customer routes"
}

# Standard customer routes
resource "mikrotik_routing_filter_rule" "customer_standard" {
  chain = mikrotik_routing_filter_chain.customer_routes.name
  rule  = <<-EOT
    if (dst in 10.0.0.0/8) {
      set bgp-local-pref 120;
      accept
    }
  EOT
  comment = "Standard customer routes"
}

# Reject anything else from customers
resource "mikrotik_routing_filter_rule" "customer_reject" {
  chain   = mikrotik_routing_filter_chain.customer_routes.name
  rule    = "reject"
  comment = "Reject unauthorized prefixes"
}

# ============================================================================
# Output Filter (What We Advertise)
# ============================================================================

resource "mikrotik_routing_filter_chain" "bgp_out" {
  name     = "bgp-output-filter"
  dynamic  = true
  comment  = "Advertise only our networks and customers"
}

# Advertise our public network
resource "mikrotik_routing_filter_rule" "out_our_network" {
  chain   = mikrotik_routing_filter_chain.bgp_out.name
  rule    = "if (dst in 203.0.113.0/24) { accept }"
  comment = "Advertise our public network"
}

# Advertise customer networks (with community)
resource "mikrotik_routing_filter_rule" "out_customers" {
  chain   = mikrotik_routing_filter_chain.bgp_out.name
  rule    = "if (dst in 10.0.0.0/8 && bgp-communities includes 65001:100) { accept }"
  comment = "Advertise customer routes"
}

# Reject everything else
resource "mikrotik_routing_filter_rule" "out_reject_all" {
  chain   = mikrotik_routing_filter_chain.bgp_out.name
  rule    = "reject"
  comment = "Don't advertise anything else"
}

# ============================================================================
# Traffic Engineering Filter (Prepend AS-Path)
# ============================================================================

resource "mikrotik_routing_filter_chain" "traffic_engineering_out" {
  name     = "traffic-engineering-out"
  dynamic  = true
  comment  = "AS-path prepending for traffic engineering"
}

# Prepend for backup path (make primary preferred)
resource "mikrotik_routing_filter_rule" "prepend_backup" {
  chain = mikrotik_routing_filter_chain.traffic_engineering_out.name
  rule  = <<-EOT
    if (dst in 203.0.113.0/24) {
      set bgp-as-path-prepend 3;
      accept
    }
  EOT
  comment = "Prepend AS 3 times for backup path"
}

resource "mikrotik_routing_filter_rule" "te_accept_rest" {
  chain   = mikrotik_routing_filter_chain.traffic_engineering_out.name
  rule    = "accept"
  comment = "Accept other routes without prepending"
}

# ============================================================================
# BGP Connections
# ============================================================================

# Primary ISP
resource "mikrotik_bgp_connection" "isp_primary" {
  name           = "isp-primary"
  template       = "default"
  remote_address = var.primary_isp_address
  remote_as      = var.primary_isp_as
  local_as       = var.local_as
  
  input_filter  = mikrotik_routing_filter_chain.isp_primary_in.name
  output_filter = mikrotik_routing_filter_chain.bgp_out.name
  
  disabled = false
  comment  = "Primary ISP - high preference"
}

# Backup ISP
resource "mikrotik_bgp_connection" "isp_backup" {
  name           = "isp-backup"
  template       = "default"
  remote_address = var.backup_isp_address
  remote_as      = var.backup_isp_as
  local_as       = var.local_as
  
  input_filter  = mikrotik_routing_filter_chain.isp_backup_in.name
  output_filter = mikrotik_routing_filter_chain.traffic_engineering_out.name
  
  disabled = false
  comment  = "Backup ISP - lower preference with AS prepend"
}

# Customer connection
resource "mikrotik_bgp_connection" "customer_a" {
  name           = "customer-a"
  template       = "default"
  remote_address = var.customer_a_address
  remote_as      = var.customer_a_as
  local_as       = var.local_as
  
  input_filter  = mikrotik_routing_filter_chain.customer_routes.name
  output_filter = "default-originate"  # Advertise default to customer
  
  disabled = false
  comment  = "Customer A - filtered with communities"
}

# ============================================================================
# Outputs
# ============================================================================

output "filter_chains" {
  description = "All configured filter chains"
  value = {
    primary_in    = mikrotik_routing_filter_chain.isp_primary_in.name
    backup_in     = mikrotik_routing_filter_chain.isp_backup_in.name
    customer_in   = mikrotik_routing_filter_chain.customer_routes.name
    standard_out  = mikrotik_routing_filter_chain.bgp_out.name
    te_out        = mikrotik_routing_filter_chain.traffic_engineering_out.name
  }
}

output "bgp_connections" {
  description = "BGP connections with their filters"
  value = {
    primary = {
      name         = mikrotik_bgp_connection.isp_primary.name
      input_filter = mikrotik_bgp_connection.isp_primary.input_filter
    }
    backup = {
      name         = mikrotik_bgp_connection.isp_backup.name
      input_filter = mikrotik_bgp_connection.isp_backup.input_filter
    }
    customer = {
      name         = mikrotik_bgp_connection.customer_a.name
      input_filter = mikrotik_bgp_connection.customer_a.input_filter
    }
  }
}
