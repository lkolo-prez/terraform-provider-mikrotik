# BGP Communities and Filtering Example
# Advanced route filtering using BGP communities
# RouterOS 7.20+

terraform {
  required_providers {
    mikrotik = {
      source = "terraform-provider-mikrotik/mikrotik"
    }
  }
}

provider "mikrotik" {
  alias    = "isp"
  host     = "192.168.88.1"
  username = "admin"
  password = "admin"
}

provider "mikrotik" {
  alias    = "customer"
  host     = "192.168.88.2"
  username = "admin"
  password = "admin"
}

# ========== ISP Router Configuration ==========

resource "mikrotik_bgp_instance_v7" "isp" {
  provider = mikrotik.isp
  
  name      = "main"
  as        = 65400
  router_id = "10.0.4.1"
  
  # ISP redistributes connected and static routes
  redistribute_connected = true
  redistribute_static    = true
  
  comment = "ISP BGP instance"
}

# BGP Template for customer connections
resource "mikrotik_bgp_template" "customer_template" {
  provider = mikrotik.isp
  
  name = "customers"
  as   = 65400
  
  # Address families
  address_families = "ip"
  
  # Accept customer routes
  input_accept_nlri = "bgp"
  
  # Accept specific communities
  # Format: AS:Value
  # 65400:100 = customer routes
  # 65400:200 = backup routes
  # 65400:300 = premium routes
  input_accept_communities = "65400:100,65400:200,65400:300"
  
  # Don't accept originated routes back
  input_accept_originated = false
  
  # Output settings
  output_default_originate = "if-installed"
  
  # Rate limiting for customer prefixes
  limit_process_routes_ipv4 = 1000
  limit_process_routes_ipv6 = 100
  
  # Timers
  hold_time      = "3m"
  keepalive_time = "1m"
  
  comment = "Template for customer connections"
}

# Customer connection with community filtering
resource "mikrotik_bgp_connection" "customer" {
  provider = mikrotik.isp
  
  name           = "customer-65500"
  instance       = mikrotik_bgp_instance_v7.isp.name
  remote_address = "192.168.1.2"
  remote_as      = 65500
  
  # Apply template
  templates = [mikrotik_bgp_template.customer_template.name]
  
  # Connection settings
  listen           = true
  address_families = "ip"
  
  # Input filtering - accept only specific communities
  input_accept_communities = "65400:100,65400:200"
  
  # Output filtering - redistribute with community tagging
  output_redistribute = "connected,static,bgp"
  
  comment = "Customer AS 65500 with community filtering"
  
  depends_on = [mikrotik_bgp_template.customer_template]
}

# ========== Customer Router Configuration ==========

resource "mikrotik_bgp_instance_v7" "customer" {
  provider = mikrotik.customer
  
  name      = "main"
  as        = 65500
  router_id = "10.0.5.1"
  
  redistribute_connected = true
  redistribute_static    = true
  
  comment = "Customer BGP instance"
}

# Template for tagging customer routes
resource "mikrotik_bgp_template" "customer_routes" {
  provider = mikrotik.customer
  
  name = "to-isp"
  as   = 65500
  
  address_families = "ip"
  
  # Accept default route from ISP
  input_accept_nlri = "bgp"
  
  # Capabilities
  capabilities = "mp,refresh,as4"
  
  comment = "Template for ISP connection"
}

resource "mikrotik_bgp_connection" "to_isp" {
  provider = mikrotik.customer
  
  name           = "to-isp-65400"
  instance       = mikrotik_bgp_instance_v7.customer.name
  remote_address = "192.168.1.1"
  remote_as      = 65400
  
  templates = [mikrotik_bgp_template.customer_routes.name]
  
  connect          = true
  address_families = "ip"
  
  # Output settings - advertise customer routes
  output_redistribute = "connected,static"
  
  comment = "Connection to ISP with community tagging"
  
  depends_on = [mikrotik_bgp_template.customer_routes]
}

# ========== Monitoring ==========

data "mikrotik_bgp_session" "isp_sessions" {
  provider = mikrotik.isp
  
  depends_on = [mikrotik_bgp_connection.customer]
}

data "mikrotik_bgp_session" "customer_sessions" {
  provider = mikrotik.customer
  
  depends_on = [mikrotik_bgp_connection.to_isp]
}

# ========== Outputs ==========

output "isp_status" {
  value = {
    customer_session = {
      established    = data.mikrotik_bgp_session.isp_sessions.established
      state          = data.mikrotik_bgp_session.isp_sessions.state
      remote_as      = data.mikrotik_bgp_session.isp_sessions.remote_as
      prefix_count   = data.mikrotik_bgp_session.isp_sessions.prefix_count
    }
    communities_accepted = "65400:100, 65400:200"
  }
  description = "ISP BGP status with community filtering"
}

output "customer_status" {
  value = {
    isp_session = {
      established  = data.mikrotik_bgp_session.customer_sessions.established
      state        = data.mikrotik_bgp_session.customer_sessions.state
      remote_as    = data.mikrotik_bgp_session.customer_sessions.remote_as
      prefix_count = data.mikrotik_bgp_session.customer_sessions.prefix_count
    }
  }
  description = "Customer BGP status"
}

output "community_usage" {
  value = {
    description = "BGP Communities for route classification"
    communities = {
      "65400:100" = "Standard customer routes"
      "65400:200" = "Backup routes (lower priority)"
      "65400:300" = "Premium routes (higher priority)"
    }
    use_case = "ISP traffic engineering and customer route management"
  }
  description = "Community values and their meanings"
}
