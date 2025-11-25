# ============================================================================
# MikroTik Connection Variables
# ============================================================================

variable "mikrotik_host" {
  description = "MikroTik router address with API port"
  type        = string
  default     = "192.168.88.1:8728"
}

variable "mikrotik_username" {
  description = "MikroTik admin username"
  type        = string
  default     = "admin"
}

variable "mikrotik_password" {
  description = "MikroTik admin password"
  type        = string
  sensitive   = true
  default     = ""
}

variable "mikrotik_tls" {
  description = "Use TLS for API connection"
  type        = bool
  default     = false
}

# ============================================================================
# BGP Configuration Variables
# ============================================================================

variable "local_as" {
  description = "Your BGP AS number"
  type        = number
  
  validation {
    condition     = var.local_as >= 1 && var.local_as <= 4294967295
    error_message = "AS number must be between 1 and 4294967295"
  }
}

variable "router_id" {
  description = "BGP Router ID (IPv4 address format)"
  type        = string
  
  validation {
    condition     = can(regex("^\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}$", var.router_id))
    error_message = "Router ID must be a valid IPv4 address"
  }
}

variable "local_bgp_address" {
  description = "Local IP address for BGP session"
  type        = string
}

# ============================================================================
# Primary ISP Variables
# ============================================================================

variable "isp_as" {
  description = "ISP BGP AS number"
  type        = number
}

variable "isp_bgp_peer" {
  description = "ISP BGP peer IP address"
  type        = string
}

variable "bgp_password" {
  description = "BGP TCP MD5 authentication password"
  type        = string
  sensitive   = true
  default     = ""
}

# ============================================================================
# Backup ISP Variables (Optional)
# ============================================================================

variable "enable_backup_isp" {
  description = "Enable backup ISP connection"
  type        = bool
  default     = false
}

variable "backup_isp_as" {
  description = "Backup ISP BGP AS number"
  type        = number
  default     = 0
}

variable "backup_isp_bgp_peer" {
  description = "Backup ISP BGP peer IP address"
  type        = string
  default     = ""
}

variable "bgp_backup_password" {
  description = "Backup BGP TCP MD5 authentication password"
  type        = string
  sensitive   = true
  default     = ""
}

# ============================================================================
# Advanced Options
# ============================================================================

variable "enable_bfd" {
  description = "Enable BFD for fast failure detection"
  type        = bool
  default     = false
}

variable "firewall_drop_rule_id" {
  description = "ID of firewall drop rule to place BGP rules before"
  type        = string
  default     = ""
}
