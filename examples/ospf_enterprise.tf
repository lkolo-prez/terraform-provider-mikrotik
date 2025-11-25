# OSPF Enterprise Configuration Example
# Multi-area OSPF with stub, NSSA, authentication, and advanced features

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

variable "ospf_auth_key" {
  type        = string
  sensitive   = true
  description = "OSPF authentication key for SHA256"
}

# ============================================================================
# CORE ROUTER CONFIGURATION (ABR connecting all areas)
# ============================================================================

# Main OSPF Instance with redistribution
resource "mikrotik_ospf_instance_v7" "enterprise" {
  name                   = "enterprise_ospf"
  version                = "2"
  router_id              = "172.16.0.1"
  redistribute_connected = false  # Don't redistribute all connected
  redistribute_static    = true   # Redistribute static routes
  redistribute_bgp       = false
  originate_default      = "if-installed"  # Only if 0.0.0.0/0 exists
  comment                = "Enterprise multi-area OSPF"
}

# ============================================================================
# AREA CONFIGURATION
# ============================================================================

# Backbone Area (required)
resource "mikrotik_ospf_area_v7" "backbone" {
  name     = "backbone"
  instance = mikrotik_ospf_instance_v7.enterprise.name
  area_id  = "0.0.0.0"
  comment  = "OSPF backbone - connects all areas"
}

# Area 1: Data Center (standard area, full routing)
resource "mikrotik_ospf_area_v7" "datacenter" {
  name     = "datacenter"
  instance = mikrotik_ospf_instance_v7.enterprise.name
  area_id  = "0.0.0.1"
  type     = "default"
  comment  = "Data center area - full LSA database"
}

# Area 2: Branch Offices (stub area - single exit point)
resource "mikrotik_ospf_area_v7" "branches" {
  name         = "branches"
  instance     = mikrotik_ospf_instance_v7.enterprise.name
  area_id      = "0.0.0.2"
  type         = "stub"
  default_cost = 100  # Cost of injected default route
  comment      = "Branch offices - stub area"
}

# Area 3: Small Remote Sites (totally stubby - minimal routing)
resource "mikrotik_ospf_area_v7" "remote_sites" {
  name         = "remote_sites"
  instance     = mikrotik_ospf_instance_v7.enterprise.name
  area_id      = "0.0.0.3"
  type         = "stub"
  no_summaries = true  # Totally stubby - only default route
  default_cost = 200
  comment      = "Small remote sites - totally stubby"
}

# Area 4: NSSA (branch with Internet connection)
resource "mikrotik_ospf_area_v7" "nssa_branch" {
  name             = "nssa_branch"
  instance         = mikrotik_ospf_instance_v7.enterprise.name
  area_id          = "0.0.0.4"
  type             = "nssa"
  nssa_translator  = "candidate"  # Participate in translator election
  nssa_propagation = true
  comment          = "Branch with Internet - NSSA area"
}

# ============================================================================
# INTERFACE TEMPLATES - BACKBONE AREA
# ============================================================================

# Core router interconnects (authenticated point-to-point)
resource "mikrotik_ospf_interface_template_v7" "backbone_core" {
  area               = mikrotik_ospf_area_v7.backbone.name
  networks           = ["172.16.0.0/24"]
  type               = "ptp"
  cost               = 10
  auth               = "sha256"
  auth_key           = var.ospf_auth_key
  auth_id            = 1
  hello_interval     = "5s"
  dead_interval      = "20s"
  comment            = "Backbone core links - authenticated"
}

# Management network (passive)
resource "mikrotik_ospf_interface_template_v7" "backbone_mgmt" {
  area     = mikrotik_ospf_area_v7.backbone.name
  networks = ["172.16.255.0/24"]
  passive  = true
  comment  = "Management network - passive"
}

# ============================================================================
# INTERFACE TEMPLATES - DATA CENTER AREA
# ============================================================================

# Data center interconnects
resource "mikrotik_ospf_interface_template_v7" "dc_core" {
  area     = mikrotik_ospf_area_v7.datacenter.name
  networks = ["10.0.0.0/16"]
  type     = "broadcast"
  cost     = 10
  priority = 200  # High priority for core switches
  auth     = "sha256"
  auth_key = var.ospf_auth_key
  auth_id  = 1
  comment  = "Data center core"
}

# Server VLANs (passive)
resource "mikrotik_ospf_interface_template_v7" "dc_servers" {
  area     = mikrotik_ospf_area_v7.datacenter.name
  networks = ["10.1.0.0/16", "10.2.0.0/16"]
  passive  = true
  comment  = "Server subnets - passive"
}

# ============================================================================
# INTERFACE TEMPLATES - BRANCH OFFICES (STUB)
# ============================================================================

# Branch office WAN links
resource "mikrotik_ospf_interface_template_v7" "branch_wan" {
  area     = mikrotik_ospf_area_v7.branches.name
  networks = ["192.168.0.0/22"]  # 192.168.0.0 - 192.168.3.255
  type     = "ptp"
  cost     = 100
  auth     = "sha256"
  auth_key = var.ospf_auth_key
  auth_id  = 1
  comment  = "Branch WAN links"
}

# Branch LAN networks (passive)
resource "mikrotik_ospf_interface_template_v7" "branch_lan" {
  area     = mikrotik_ospf_area_v7.branches.name
  networks = ["192.168.4.0/22", "192.168.8.0/22"]
  passive  = true
  comment  = "Branch LAN networks"
}

# ============================================================================
# INTERFACE TEMPLATES - REMOTE SITES (TOTALLY STUBBY)
# ============================================================================

# Remote site WAN (slow links)
resource "mikrotik_ospf_interface_template_v7" "remote_wan" {
  area               = mikrotik_ospf_area_v7.remote_sites.name
  networks           = ["192.168.200.0/24"]
  type               = "ptp"
  cost               = 1000  # High cost for slow links
  hello_interval     = "30s"  # Slower timers for stability
  dead_interval      = "120s"
  auth               = "md5"
  auth_key           = var.ospf_auth_key
  auth_id            = 1
  comment            = "Remote site WAN - slow links"
}

# ============================================================================
# INTERFACE TEMPLATES - NSSA BRANCH
# ============================================================================

# NSSA branch WAN
resource "mikrotik_ospf_interface_template_v7" "nssa_wan" {
  area     = mikrotik_ospf_area_v7.nssa_branch.name
  networks = ["192.168.100.0/24"]
  type     = "ptp"
  cost     = 100
  auth     = "sha256"
  auth_key = var.ospf_auth_key
  auth_id  = 1
  comment  = "NSSA branch WAN"
}

# NSSA branch LAN (passive)
resource "mikrotik_ospf_interface_template_v7" "nssa_lan" {
  area     = mikrotik_ospf_area_v7.nssa_branch.name
  networks = ["192.168.101.0/24"]
  passive  = true
  comment  = "NSSA branch LAN"
}

# ============================================================================
# ROUTING FILTERS (optional)
# ============================================================================

# Filter incoming routes
resource "mikrotik_routing_filter_chain" "ospf_in" {
  name    = "ospf_in_filter"
  comment = "Filter incoming OSPF routes"
}

# Reject bogon networks
resource "mikrotik_routing_filter_rule" "reject_bogons" {
  chain = mikrotik_routing_filter_chain.ospf_in.name
  rule = <<-EOT
    if (dst in 0.0.0.0/8 || dst in 127.0.0.0/8 || 
        dst in 169.254.0.0/16 || dst in 224.0.0.0/4) {
      reject
    }
  EOT
}

# Filter outgoing routes
resource "mikrotik_routing_filter_chain" "ospf_out" {
  name    = "ospf_out_filter"
  comment = "Filter outgoing redistributed routes"
}

# Only redistribute specific static routes
resource "mikrotik_routing_filter_rule" "allow_public" {
  chain = mikrotik_routing_filter_chain.ospf_out.name
  rule = <<-EOT
    if (dst in 203.0.113.0/24) {
      set ospf-ext-metric 1000;
      accept
    } else {
      reject
    }
  EOT
}

# Apply filters to instance (uncomment to use)
# resource "mikrotik_ospf_instance_v7" "enterprise" {
#   in_filter_chain  = mikrotik_routing_filter_chain.ospf_in.name
#   out_filter_chain = mikrotik_routing_filter_chain.ospf_out.name
# }

# ============================================================================
# OUTPUTS
# ============================================================================

output "ospf_instance" {
  value = {
    id        = mikrotik_ospf_instance_v7.enterprise.id
    name      = mikrotik_ospf_instance_v7.enterprise.name
    router_id = mikrotik_ospf_instance_v7.enterprise.router_id
  }
  description = "OSPF instance details"
}

output "ospf_areas" {
  value = {
    backbone     = mikrotik_ospf_area_v7.backbone.name
    datacenter   = mikrotik_ospf_area_v7.datacenter.name
    branches     = mikrotik_ospf_area_v7.branches.name
    remote_sites = mikrotik_ospf_area_v7.remote_sites.name
    nssa_branch  = mikrotik_ospf_area_v7.nssa_branch.name
  }
  description = "Configured OSPF areas"
}

output "area_types" {
  value = {
    backbone     = "standard (0.0.0.0)"
    datacenter   = "standard"
    branches     = "stub"
    remote_sites = "totally stubby"
    nssa_branch  = "NSSA"
  }
  description = "Area types"
}
