variable "mikrotik_host" {
  description = "RouterOS API address, e.g. 192.168.88.1:8729 (API-SSL)."
  type        = string
}

variable "mikrotik_username" {
  description = "RouterOS user with write access."
  type        = string
  default     = "admin"
}

variable "mikrotik_password" {
  description = "Password of the RouterOS user."
  type        = string
  sensitive   = true
}

variable "mikrotik_insecure" {
  description = "Skip TLS certificate verification (lab only)."
  type        = bool
  default     = false
}

variable "lan_bridge" {
  description = "Name of the LAN bridge created by this example."
  type        = string
  default     = "bridge-lan"
}

variable "router_id" {
  description = "OSPF router ID (IPv4 dotted notation, also used for OSPFv3)."
  type        = string
  default     = "10.255.0.1"
}

variable "vrrp_virtual_ip" {
  description = "Virtual gateway address shared by the VRRP pair (CIDR)."
  type        = string
  default     = "192.168.10.1/24"
}

variable "vrrp_priority" {
  description = "VRRP priority: higher wins (e.g. 200 on the primary, 100 on the secondary)."
  type        = number
  default     = 200
}

variable "wifi_ssid" {
  description = "SSID broadcast by the access point."
  type        = string
  default     = "aircloud-lab"
}

variable "wifi_passphrase" {
  description = "WPA2/WPA3 passphrase (min. 8 characters)."
  type        = string
  sensitive   = true
}

variable "wifi_country" {
  description = "Regulatory country for the radio."
  type        = string
  default     = "Poland"
}

variable "backup_password" {
  description = "Password used to encrypt the binary backup file."
  type        = string
  sensitive   = true
}
