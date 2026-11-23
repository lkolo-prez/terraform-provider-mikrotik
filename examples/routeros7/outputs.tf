output "ospf_v3_instance" {
  description = "Name of the OSPFv3 instance."
  value       = mikrotik_ospf_instance_v7.v3.name
}

output "vrrp_interface" {
  description = "VRRP interface carrying the virtual gateway."
  value       = mikrotik_interface_vrrp.gateway.name
}

output "wifi_interface" {
  description = "Wi-Fi interface and whether it is running."
  value = {
    name    = mikrotik_interface_wifi.wifi1.name
    running = mikrotik_interface_wifi.wifi1.running
  }
}

output "backup" {
  description = "Backup file name, size and creation time reported by RouterOS."
  value = {
    name          = mikrotik_system_backup.after_change.name
    size          = mikrotik_system_backup.after_change.size
    creation_time = mikrotik_system_backup.after_change.creation_time
  }
}
