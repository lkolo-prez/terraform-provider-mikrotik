# Example: VRRP High Availability Configuration
# Two routers (master and backup) sharing a virtual IP for gateway redundancy

terraform {
  required_providers {
    mikrotik = {
      source = "lkolo-prez/mikrotik"
    }
  }
}

# Configuration for MASTER router
provider "mikrotik" {
  alias    = "master"
  host     = "192.168.1.10"
  username = "admin"
  password = var.master_password
  tls      = true
}

# Configuration for BACKUP router
provider "mikrotik" {
  alias    = "backup"
  host     = "192.168.1.11"
  username = "admin"
  password = var.backup_password
  tls      = true
}

variable "vrrp_password" {
  type      = string
  sensitive = true
  description = "VRRP authentication password"
}

variable "master_password" {
  type      = string
  sensitive = true
}

variable "backup_password" {
  type      = string
  sensitive = true
}

# MASTER Router Configuration
resource "mikrotik_interface_vrrp" "gateway_master" {
  provider = mikrotik.master
  
  name            = "vrrp-gateway"
  interface       = "ether1"
  vrid            = 10
  priority        = 254          # Higher priority = Master
  
  version         = 3
  v3_protocol     = "ipv4"
  
  authentication  = "simple"
  password        = var.vrrp_password
  
  interval        = "1s"
  preemption_mode = true
  
  comment         = "Master VRRP interface for gateway HA"
}

# Virtual IP on Master
resource "mikrotik_ip_address" "vrrp_vip_master" {
  provider = mikrotik.master
  
  address   = "192.168.1.1/24"
  interface = mikrotik_interface_vrrp.gateway_master.name
  comment   = "Virtual IP for VRRP gateway"
}

# BACKUP Router Configuration
resource "mikrotik_interface_vrrp" "gateway_backup" {
  provider = mikrotik.backup
  
  name            = "vrrp-gateway"
  interface       = "ether1"
  vrid            = 10           # Same VRID as master
  priority        = 100          # Lower priority = Backup
  
  version         = 3
  v3_protocol     = "ipv4"
  
  authentication  = "simple"
  password        = var.vrrp_password
  
  interval        = "1s"
  preemption_mode = true
  
  comment         = "Backup VRRP interface for gateway HA"
}

# Virtual IP on Backup
resource "mikrotik_ip_address" "vrrp_vip_backup" {
  provider = mikrotik.backup
  
  address   = "192.168.1.1/24"
  interface = mikrotik_interface_vrrp.gateway_backup.name
  comment   = "Virtual IP for VRRP gateway"
}

# Optional: Scripts for state change notifications
resource "mikrotik_script" "vrrp_master_notify" {
  provider = mikrotik.master
  
  name   = "vrrp-master-notify"
  source = <<-EOT
    :log info "VRRP: Transitioned to MASTER state"
    /tool/e-mail/send to="noc@company.com" subject="VRRP MASTER" body="Router became MASTER"
  EOT
  
  policy = ["read", "write", "policy", "test"]
}

resource "mikrotik_script" "vrrp_backup_notify" {
  provider = mikrotik.master
  
  name   = "vrrp-backup-notify"
  source = <<-EOT
    :log warning "VRRP: Transitioned to BACKUP state"
    /tool/e-mail/send to="noc@company.com" subject="VRRP BACKUP" body="Router became BACKUP"
  EOT
  
  policy = ["read", "write", "policy", "test"]
}

# Update VRRP with script hooks
resource "mikrotik_interface_vrrp" "gateway_master_with_scripts" {
  provider = mikrotik.master
  
  name            = "vrrp-gateway"
  interface       = "ether1"
  vrid            = 10
  priority        = 254
  
  version         = 3
  v3_protocol     = "ipv4"
  
  authentication  = "simple"
  password        = var.vrrp_password
  
  interval        = "1s"
  preemption_mode = true
  
  on_master       = "/system/script/run ${mikrotik_script.vrrp_master_notify.name}"
  on_backup       = "/system/script/run ${mikrotik_script.vrrp_backup_notify.name}"
  
  comment         = "Master VRRP with notifications"
  
  depends_on = [
    mikrotik_script.vrrp_master_notify,
    mikrotik_script.vrrp_backup_notify
  ]
}

# Example: Multiple VRRP groups for load balancing
# Master for LAN1, Backup for LAN2
resource "mikrotik_interface_vrrp" "vrrp_lan1_master" {
  provider = mikrotik.master
  
  name      = "vrrp-lan1"
  interface = "ether2"
  vrid      = 11
  priority  = 254    # Master for LAN1
  
  version   = 3
  comment   = "VRRP for LAN1 - Master"
}

resource "mikrotik_interface_vrrp" "vrrp_lan2_backup" {
  provider = mikrotik.master
  
  name      = "vrrp-lan2"
  interface = "ether3"
  vrid      = 12
  priority  = 100    # Backup for LAN2
  
  version   = 3
  comment   = "VRRP for LAN2 - Backup"
}

# Corresponding configs on backup router with reversed priorities
resource "mikrotik_interface_vrrp" "vrrp_lan1_backup" {
  provider = mikrotik.backup
  
  name      = "vrrp-lan1"
  interface = "ether2"
  vrid      = 11
  priority  = 100    # Backup for LAN1
  
  version   = 3
  comment   = "VRRP for LAN1 - Backup"
}

resource "mikrotik_interface_vrrp" "vrrp_lan2_master" {
  provider = mikrotik.backup
  
  name      = "vrrp-lan2"
  interface = "ether3"
  vrid      = 12
  priority  = 254    # Master for LAN2
  
  version   = 3
  comment   = "VRRP for LAN2 - Master"
}

# Outputs
output "master_vrrp_status" {
  value = {
    id      = mikrotik_interface_vrrp.gateway_master.id
    name    = mikrotik_interface_vrrp.gateway_master.name
    running = mikrotik_interface_vrrp.gateway_master.running
  }
  description = "Master VRRP interface status"
}

output "backup_vrrp_status" {
  value = {
    id      = mikrotik_interface_vrrp.gateway_backup.id
    name    = mikrotik_interface_vrrp.gateway_backup.name
    running = mikrotik_interface_vrrp.gateway_backup.running
  }
  description = "Backup VRRP interface status"
}
